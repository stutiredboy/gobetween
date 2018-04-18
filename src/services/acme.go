package services

import (
	"../config"
	"../core"
	"../server/tcp"
	"context"
	"fmt"
	"golang.org/x/crypto/acme/autocert"
	"net/http"
	"sync"
)

type AcmeService struct {
	certMan *autocert.Manager
	hosts   map[string]bool
	sync.RWMutex
}

func NewAcmeService(cfg config.Config) core.Service {

	if cfg.Acme == nil {
		return nil
	}

	a := &AcmeService{
		certMan: &autocert.Manager{
			Cache:  autocert.DirCache(cfg.Acme.CacheDir),
			Prompt: autocert.AcceptTOS,
		},
		hosts: make(map[string]bool),
	}

	a.certMan.HostPolicy = func(_ context.Context, host string) error {
		a.RLock()
		defer a.RUnlock()

		if a.hosts[host] {
			return nil
		}

		return fmt.Errorf("Acme: host %s is not configured", host)
	}

	//accept http challenge
	if cfg.Acme.Challenge == "http" {
		go http.ListenAndServe(cfg.Acme.HttpBind, a.certMan.HTTPHandler(nil))
	}

	return a

}

func (a *AcmeService) String() string {
	return ""
}

func (a *AcmeService) Apply(server core.Server) error {

	if a == nil {
		return nil
	}

	serverCfg := server.Cfg()

	if serverCfg.Tls == nil {
		return nil
	}

	tcpServer, ok := server.(*tcp.Server)

	if !ok {
		return nil
	}

	tcpServer.GetCertificate = a.certMan.GetCertificate

	a.Lock()
	defer a.Unlock()

	for _, host := range serverCfg.Tls.AcmeHosts {

		if a.hosts[host] {
			return fmt.Errorf("Acme host %s is already configured", host)
		}

		a.hosts[host] = true
	}

	return nil
}

func (a *AcmeService) Forget(server core.Server) error {

	serverCfg := server.Cfg()

	if serverCfg.Tls == nil {
		return nil
	}

	a.Lock()
	defer a.Unlock()

	for _, host := range serverCfg.Tls.AcmeHosts {
		delete(a.hosts, host)
	}

	return nil
}
