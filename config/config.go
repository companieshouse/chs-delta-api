// Package config defines the environment variable and command-line flags
package config

import (
	"errors"
	"github.com/companieshouse/chs.go/log"
	"github.com/companieshouse/gofigure"
	"sync"
)

var cfg *Config
var mtx sync.Mutex
var mandatoryElementMissing bool

// Config defines the configuration options for this service.
type Config struct {
	BindAddr          string    `env:"BIND_ADDR" flag:"bind-addr" flagDesc:"Bind address"`
	BrokerAddr        []string  `env:"KAFKA_BROKER_ADDR" flag:"broker-addr" flagDesc:"Kafka broker address (Comma separated list if there is more than one address"`
	SchemaRegistryURL string    `env:"SCHEMA_REGISTRY_URL" flag:"schema-registry-url" flagDesc:"URL for Kafka Schema Registry"`
	OfficerDeltaTopic string 	`env:"OFFICER_DELTA_TOPIC" flag:"officer-delta-topic" flagDesc:"Topic for officer deltas"`
	OpenApiSpec		  string    `env:"OPEN_API_SPEC" flag:"open-api-spec" flagDesc:"OpenAPI schema location"`
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

	mandatoryElementMissing = validateConfigs(cfg)

	if mandatoryElementMissing {
		return nil, errors.New("mandatory configs missing from environment")
	}

	return cfg, nil
}

func validateConfigs(cfg *Config ) bool {

	mandatoryElementMissing = false

	if cfg.BindAddr == "" {
		log.Info("BIND_ADDR not set in environment")
		mandatoryElementMissing = true
	}

	if len(cfg.BrokerAddr) == 0 {
		log.Info("KAFKA_BROKER_ADDR not set in environment")
		mandatoryElementMissing = true
	}

	if cfg.SchemaRegistryURL == "" {
		log.Info("SCHEMA_REGISTRY_URL not set in environment")
		mandatoryElementMissing = true
	}

	if cfg.OfficerDeltaTopic == "" {
		log.Info("OFFICER_DELTA_TOPIC not set in environment")
		mandatoryElementMissing = true
	}

	if cfg.OpenApiSpec == "" {
		log.Info("OPEN_API_SPEC not set in environment")
		mandatoryElementMissing = true
	}
	return mandatoryElementMissing
}