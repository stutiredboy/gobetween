package services

import (
	"../config"
	"../core"
	"../logging"
)

/**
 * Registry of factory methods for Services
 */
var registry = make(map[string]func(config.Config) core.Service)

func init() {
	registry["acme"] = NewAcmeService
}

func Services(cfg config.Config) []core.Service {
	log := logging.For("services")

	result := make([]core.Service, 0)

	for name, constructor := range registry {
		service := constructor(cfg)
		log.Info("Creating ", name, " ", service)
		result = append(result, service)
	}

	return result
}
