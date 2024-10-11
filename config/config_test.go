package config

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUnitGetHasErrors(t *testing.T) {
	Convey("When I try to get the config via the Get method", t, func() {
		cfg, err := Get()
		Convey("Then I am given an error when config validation function is executed", func() {
			So(cfg, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})
	})
}

func TestUnitGetNoErrors(t *testing.T) {
	_ = os.Setenv("BIND_ADDR", "bind_addr")
	_ = os.Setenv("KAFKA_BROKER_ADDR", "kafka_broker_addr,kafka_broker_addr")
	_ = os.Setenv("SCHEMA_REGISTRY_URL", "schema_registry_url")
	_ = os.Setenv("OFFICER_DELTA_TOPIC", "officer_delta_topic")
	_ = os.Setenv("INSOLVENCY_DELTA_TOPIC", "insolvency_delta_topic")
	_ = os.Setenv("CHARGES_DELTA_TOPIC", "charges_delta_topic")
	_ = os.Setenv("DISQUALIFIED_OFFICERS_DELTA_TOPIC", "disqualified-officers-delta-topic")
	_ = os.Setenv("COMPANY_DELTA_TOPIC", "company-delta-topic")
	_ = os.Setenv("EXEMPTION_DELTA_TOPIC", "exemption-delta-topic")
	_ = os.Setenv("PSC_STATEMENT_DELTA_TOPIC", "psc-statement-delta-topic")
	_ = os.Setenv("PSC_DELTA_TOPIC", "psc-delta-topic")
	_ = os.Setenv("FILING_HISTORY_DELTA_TOPIC", "filing-history-delta-topic")
	_ = os.Setenv("DOCUMENT_STORE_DELTA_TOPIC", "document-store-delta-topic")
	_ = os.Setenv("REGISTERS_DELTA_TOPIC", "registers-delta-topic")
	_ = os.Setenv("ACSP_PROFILE_DELTA_TOPIC", "acsp-profile-delta-topic")
	_ = os.Setenv("OPEN_API_SPEC", "open_api_spec")
	Convey("When I try to get the config via the Get method and all config vars are provided", t, func() {
		cfg, err := Get()
		Convey("Then I am given a correctly configured config", func() {
			So(cfg, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})
	})
	os.Clearenv()
}
