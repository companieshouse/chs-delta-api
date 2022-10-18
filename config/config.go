// Package config defines the environment variable and command-line flags
package config

import (
	"errors"
	"sync"

	"github.com/companieshouse/chs.go/log"
	"github.com/companieshouse/gofigure"
)

var (
	cfg                *Config
	mtx                sync.Mutex
	CallValidateConfig = validateConfigs
)

// Config defines the configuration options for this service.
type Config struct {
	BindAddr                    string   `env:"BIND_ADDR" flag:"bind-addr" flagDesc:"Bind address"`
	BrokerAddr                  []string `env:"KAFKA_BROKER_ADDR" flag:"broker-addr" flagDesc:"Kafka broker address (Comma separated list if there is more than one address)"`
	SchemaRegistryURL           string   `env:"SCHEMA_REGISTRY_URL" flag:"schema-registry-url" flagDesc:"URL for Kafka Schema Registry"`
	OfficerDeltaTopic           string   `env:"OFFICER_DELTA_TOPIC" flag:"officer-delta-topic" flagDesc:"Topic for officer deltas"`
	OpenApiSpec                 string   `env:"OPEN_API_SPEC" flag:"open-api-spec" flagDesc:"OpenAPI schema location"`
	InsolvencyDeltaTopic        string   `env:"INSOLVENCY_DELTA_TOPIC" flag:"insolvency-delta-topic" flagDesc:"Topic for insolvency deltas"`
	ChargesDeltaTopic           string   `env:"CHARGES_DELTA_TOPIC" flag:"charges-delta-topic" flagDesc:"Topic for charges deltas"`
	DisqualifiedDeltaTopic      string   `env:"DISQUALIFIED_OFFICERS_DELTA_TOPIC" flag:"disqualified-officers-delta-topic" flagDesc:"Topic for disqualification deltas"`
	CompanyDeltaTopic           string   `env:"COMPANY_DELTA_TOPIC" flag:"company-delta-topic" flagDesc:"Topic for company deltas"`
	CompanyExemptionsDeltaTopic string   `env:"COMPANY_EXEMPTIONS_DELTA_TOPIC" flag:"company-exemptions-delta-topic" flagDesc:"Topic for company exemptions"`
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

	err = CallValidateConfig(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func validateConfigs(cfg *Config) error {

	mandatoryElementMissing := false

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

	if cfg.InsolvencyDeltaTopic == "" {
		log.Info("INSOLVENCY_DELTA_TOPIC not set in environment")
		mandatoryElementMissing = true
	}

	if cfg.ChargesDeltaTopic == "" {
		log.Info("CHARGES_DELTA_TOPIC not set in environment")
		mandatoryElementMissing = true
	}

	if cfg.DisqualifiedDeltaTopic == "" {
		log.Info("DISQUALIFIED_OFFICERS_DELTA_TOPIC not set in environment")
		mandatoryElementMissing = true
	}

	if cfg.OpenApiSpec == "" {
		log.Info("OPEN_API_SPEC not set in environment")
		mandatoryElementMissing = true
	}

	if mandatoryElementMissing {
		return errors.New("mandatory configs missing from environment")
	}

	return nil
}
