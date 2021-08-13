package services

import (
	"errors"
	"github.com/companieshouse/chs-delta-api/config"
	"github.com/companieshouse/chs.go/kafka/producer"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

const (
	Topic      = "chs-delta"
	Data       = `{"test" : "value"}`
	BadSchema  = `"bad_schema_value"`
	ContextId  = "contextId"
	GoodSchema = `{"type":"record","namespace":"delta","name":"delta","doc":"SchemaforthedeltathatwillbeusedtotransferdatafromCHIPStoCHS.",
"fields":[{"name":"data","type":"string","doc":"PayloadthatwillbetransferredfromCHIPStoCHSviaKafka"},
{"name":"attempt","type":"int","default":0,"doc":"NumberofattemptstoretrypublishingthemessagetoKafkaTopic"},
{"name":"context_id","type":"string","doc":"Loggingcontextidusedtotracktherequestacrossservices"}]}`
)

// TestUnitNewKafkaService asserts that the KafkaService constructor returns a non-nil reference to a KafkaServiceImpl.
func TestUnitNewKafkaService(t *testing.T) {
	Convey("Given I want a new KafkaService", t, func() {
		k := NewKafkaService()

		Convey("Then I am given a service back that is not nil", func() {
			So(k, ShouldNotBeNil)
		})
	})
}

// TestUnitKafkaServiceInitSuccessful asserts that the Init method successfully creates and initialises a KafkaServiceImpl.
func TestUnitKafkaServiceInitSuccessful(t *testing.T) {

	config.CallValidateConfig = func(cfg *config.Config) error {
		return nil
	}
	cfg, _ := config.Get()

	Convey("Given a call to init a Kafka service", t, func() {
		k := NewKafkaService()

		callSchemaGet = func(url, name string) (string, error) {
			return "mock_url", nil
		}

		callProducerNew = func(config *producer.Config) (*producer.Producer, error) {
			return &producer.Producer{}, nil
		}

		err := k.Init(cfg)

		Convey("Then the error is nil", func() {
			So(err, ShouldBeNil)
		})
	})
}

// TestUnitKafkaServiceInitGetSchemaFails asserts that errors are captured and returned when retrieving a schema fails.
func TestUnitKafkaServiceInitGetSchemaFails(t *testing.T) {

	cfg, _ := config.Get()

	Convey("Given a call to init a Kafka service", t, func() {
		k := NewKafkaService()

		callSchemaGet = func(url, name string) (string, error) {
			return "", errors.New("error retrieving schema")
		}

		callProducerNew = func(config *producer.Config) (*producer.Producer, error) {
			return &producer.Producer{}, nil
		}

		err := k.Init(cfg)

		Convey("Then the error is not nil", func() {
			So(err, ShouldNotBeNil)
		})
	})
}

// TestUnitKafkaServiceInitNewProducerFails asserts that errors are captured and returned when creating a producer fails.
func TestUnitKafkaServiceInitNewProducerFails(t *testing.T) {

	cfg, _ := config.Get()

	Convey("Given a call to init a Kafka service", t, func() {
		k := NewKafkaService()

		callSchemaGet = func(url, name string) (string, error) {
			return "mock_url", nil
		}

		callProducerNew = func(config *producer.Config) (*producer.Producer, error) {
			return nil, errors.New("error creating producer")
		}

		err := k.Init(cfg)

		Convey("Then the error is not nil", func() {
			So(err, ShouldNotBeNil)

		})
	})
}

// TestUnitSendMessageSuccessfully asserts that when no errors occur, a message can be published onto a kafka topic.
func TestUnitSendMessageSuccessfully(t *testing.T) {
	Convey("Given I have a Kafka service", t, func() {
		k := NewKafkaService()
		k.schema = GoodSchema

		Convey("When I call to send a message via the producer", func() {
			callSend = func(k *KafkaServiceImpl, msg *producer.Message) (int32, int64, error) {
				return int32(0), int64(0), nil
			}

			err := k.SendMessage(Topic, Data, ContextId)

			Convey("Then there are no errors", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}

// TestUnitSendMessageFailsSchemaMarshalling asserts that errors are handled and returned when marshalling a schema fails.
func TestUnitSendMessageFailsSchemaMarshalling(t *testing.T) {
	Convey("Given I have a Kafka service", t, func() {
		k := NewKafkaService()
		k.schema = BadSchema

		Convey("When I call to send a message via the producer", func() {

			err := k.SendMessage(Topic, Data, ContextId)

			Convey("Then there are errors returned", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

// TestUnitSendMessageFailsWithError asserts that errors are captured and returned when calling to send a message fails.
func TestUnitSendMessageFailsWithError(t *testing.T) {
	Convey("Given I have a Kafka service", t, func() {
		k := NewKafkaService()
		k.schema = GoodSchema

		Convey("When I call to send a message via the producer", func() {
			callSend = func(k *KafkaServiceImpl, msg *producer.Message) (int32, int64, error) {
				return int32(0), int64(0), errors.New("error sending to kafka producer")
			}

			err := k.SendMessage(Topic, Data, ContextId)

			Convey("Then there are errors returned", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}
