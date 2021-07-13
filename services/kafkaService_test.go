package services

import (
	"errors"
	"github.com/companieshouse/chs-delta-api/config"
	"github.com/companieshouse/chs.go/kafka/producer"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

const (
	topic  = "officer-delta"
	data   = `{"test" : "value"}`
	badSchema = `"bad_schema_value"`
	s = `{"type":"record","namespace":"delta","name":"delta","doc":"SchemaforthedeltathatwillbeusedtotransferdatafromCHIPStoCHS.",
"fields":[{"name":"data","type":"string","doc":"PayloadthatwillbetransferredfromCHIPStoCHSviaKafka"},
{"name":"attempt","type":"int","default":0,"doc":"NumberofattemptstoretrypublishingthemessagetoKafkaTopic"},
{"name":"context_id","type":"string","doc":"Loggingcontextidusedtotracktherequestacrossservices"}]}`
)

func TestNewKafkaService(t *testing.T) {
	Convey("Given I want a new KafkaService", t, func() {
		Ksvc := NewKafkaService()

		Convey("Then I am given a service back that is not nil", func() {
			So(Ksvc, ShouldNotBeNil)
		})
	})
}

func TestKafkaServiceInitSuccessful(t *testing.T) {

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

func TestKafkaServiceInitGetSchemaFails(t *testing.T) {

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

func TestKafkaServiceInitNewProducerFails(t *testing.T) {

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

func TestSendMessageSuccessfully(t *testing.T) {
	Convey("Given I have a Kafka service", t, func() {
		k := NewKafkaService()
		k.Schema = s

		Convey("When I call to send a message via the producer", func() {
			callSend = func(k *KafkaServiceImpl, msg *producer.Message) (int32, int64, error) {
				return int32(0), int64(0), nil
			}

			err := k.SendMessage(topic, data)

			Convey("Then there are no errors", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestSendMessageFailsSchemaMarshalling(t *testing.T) {
	Convey("Given I have a Kafka service", t, func() {
		k := NewKafkaService()
		k.Schema = badSchema

		Convey("When I call to send a message via the producer", func() {

			err := k.SendMessage(topic, data)

			Convey("Then there are errors returned", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestSendMessageFailsWithError(t *testing.T) {
	Convey("Given I have a Kafka service", t, func() {
		k := NewKafkaService()
		k.Schema = s

		Convey("When I call to send a message via the producer", func() {
			callSend = func(k *KafkaServiceImpl, msg *producer.Message) (int32, int64, error) {
				return int32(0), int64(0), errors.New("error sending to kafka producer")
			}

			err := k.SendMessage(topic, data)

			Convey("Then there are errors returned", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}
