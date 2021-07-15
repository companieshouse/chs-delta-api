// Package config defines the environment variable and command-line flags
package config

import (
	"sync"

	"github.com/companieshouse/gofigure"
)

var cfg *Config
var mtx sync.Mutex

// Config defines the configuration options for this service.
type Config struct {
	BindAddr          string   `env:"BIND_ADDR" flag:"bind-addr" flagDesc:"Bind address"`
	BrokerAddr        []string `env:"KAFKA_BROKER_ADDR" flag:"broker-addr" flagDesc:"Kafka broker address"`
	SchemaRegistryURL string   `env:"SCHEMA_REGISTRY_URL" flag:"schema-registry-url" flagDesc:"URL for Schema Registry"`
	OfficerDeltaTopic string `env:"OFFICER_DELTA_TOPIC" flag:"officer-delta-topic" flagDesc:"Topic for the officer delta"`
}

// Get returns a pointer to a Config instance populated with values from environment or command-line flags
func Get() (*Config, error) {
	mtx.Lock()
	defer mtx.Unlock()

	if cfg != nil {
		return cfg, nil
	}

	cfg = &Config{}

	err := gofigure.Gofigure(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
