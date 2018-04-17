package acme

/*
 * TODO:
 * 1. Listen to bind host/port
 * 2. Aggregate all acme_hosts from all servers
 * 3. Reconfigure if server list or configuration changes
 */

import (
	"../config"
	"../core"
	"../server/tcp"
	"context"
	"crypto/tls"
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

func NewAcmeService(cfg config.Config) *AcmeService {

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
		go http.ListenAndServe(cfg.Acme.Bind, a.certMan.HTTPHandler(nil))
	}

	return a

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

	tcpServer.GetCertificate = a.GetCertificate

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

func (a *AcmeService) GetCertificate(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	return a.certMan.GetCertificate(hello)
}
