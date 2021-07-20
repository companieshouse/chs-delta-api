// Package config defines the environment variable and command-line flags
package config

import (
	"errors"
	"github.com/companieshouse/gofigure"
	"log"
	"sync"
)

var cfg *Config
var mtx sync.Mutex

// Config defines the configuration options for this service.
type Config struct {
	BindAddr          string    `env:"BIND_ADDR" flag:"bind-addr" flagDesc:"Bind address"`
	BrokerAddr        []string  `env:"KAFKA_BROKER_ADDR" flag:"broker-addr" flagDesc:"Kafka broker address"`
	SchemaRegistryURL string    `env:"SCHEMA_REGISTRY_URL" flag:"schema-registry-url" flagDesc:"URL for Schema Registry"`
	OfficerDeltaTopic string 	`env:"OFFICER_DELTA_TOPIC" flag:"officer-delta-topic" flagDesc:"Topic for the officer delta"`
	LogLevel		  string  	`env:"LOGLEVEL" flag:"log-level" flagDesc:"Logging level to set i.e error, info, debug"`
	OpenApiSpec		  string    `env:"OPEN_API_SPEC" flag:"open-api-spec" flagDesc:"OpenAPI schema location"`
}

func Get() (*Config, error) {
	mtx.Lock()
	defer mtx.Unlock()
	mandatoryElementMissing := false

	if cfg != nil {
		return cfg, nil
	}

	cfg = &Config{}

	err := gofigure.Gofigure(cfg)
	if err != nil {
		return nil, err
	}

	if cfg.BindAddr == "" {
		log.Printf("BIND_ADDR not set in environment")
		mandatoryElementMissing = true
	}

	if len(cfg.BrokerAddr) == 0 {
		log.Printf("KAFKA_BROKER_ADDR not set in environment")
		mandatoryElementMissing = true
	}

	if cfg.SchemaRegistryURL == "" {
		log.Printf("SCHEMA_REGISTRY_URL not set in environment")
		mandatoryElementMissing = true
	}

	if cfg.OfficerDeltaTopic == "" {
		log.Printf("OFFICER_DELTA_TOPIC not set in environment")
		mandatoryElementMissing = true
	}

	if cfg.OpenApiSpec == "" {
		log.Printf("OPEN_API_SPEC not set in environment")
		mandatoryElementMissing = true
	}

	if mandatoryElementMissing {
		return nil, errors.New("mandatory configs missing from environment")
	}

	return cfg, nil
}