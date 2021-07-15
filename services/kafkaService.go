package services

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/companieshouse/chs-delta-api/config"
	"github.com/companieshouse/chs-delta-api/models"
	"github.com/companieshouse/chs.go/avro"
	"github.com/companieshouse/chs.go/avro/schema"
	"github.com/companieshouse/chs.go/kafka/producer"
	"github.com/companieshouse/chs.go/log"
)

const (
	SchemaName = "chs-delta"
)

// Used for unit testing. By Adding variables which link to certain package level functions / methods, we can
// override them during unit testing and change them to point to mock implementations to assert functionality.
var (
	callSchemaGet   = schema.Get
	callProducerNew = producer.New
	callSend        = sendViaProducer
)

// KafkaService defines all Methods needed to successfully send a message onto a Kafka topic.
type KafkaService interface {
	Init(cfg *config.Config) error
	SendMessage(topic, data string) error
}

// KafkaServiceImpl is a concrete implementation of the KafkaService interface.
type KafkaServiceImpl struct {
	schema string
	P      *producer.Producer
}

// NewKafkaService returns a KafkaServiceImpl that isn't configured.
func NewKafkaService() KafkaServiceImpl {
	return KafkaServiceImpl{}
}

// Init initialises a KafkaService using a provided config.
func (kSvc *KafkaServiceImpl) Init(cfg *config.Config) error {

	// Initialise the avro schema.
	sch, err := initSchema(cfg)
	if err != nil {
		return err
	}

	// Initialise the kafka producer.
	p, err := initProducer(cfg)
	if err != nil {
		return err
	}

	kSvc.schema = sch
	kSvc.P = p

	return nil
}

func initSchema(cfg *config.Config) (string, error) {
	// Retrieve the generic chs-delta avro schema.
	log.Trace("Get schema from Avro", log.Data{"schema_name": SchemaName})
	sch, err := callSchemaGet(cfg.SchemaRegistryURL, SchemaName)
	if err != nil {
		log.Error(fmt.Errorf("error receiving %s schema: %s", SchemaName, err))
		return "", err
	}
	log.Trace("Successfully received schema", log.Data{"schema_name": SchemaName})
	return sch, nil
}

func initProducer(cfg *config.Config) (*producer.Producer, error) {
	// Create a new Kafka Producer which will be used to publish our message onto a given Kafka topic.
	log.Trace("Using Streaming Kafka broker Address", log.Data{"Brokers": cfg.BrokerAddr})
	p, err := callProducerNew(&producer.Config{Acks: &producer.WaitForAll, BrokerAddrs: cfg.BrokerAddr})
	if err != nil {
		log.Error(fmt.Errorf("error initialising producer: %s", err))
		return nil, err
	}

	return p, nil
}

// SendMessage publishes a given data string retrieved from a REST request onto a chosen Kafka topic.
func (kSvc *KafkaServiceImpl) SendMessage(topic, data string) error {

	// Retrieve our chs-delta avro schema using the chs go avro package.
	chsDeltaAvro := &avro.Schema{
		Definition: kSvc.schema,
	}

	// Construct a chs-delta using provided data.
	deltaData := models.ChsDelta{
		ContextId: "uu-aa-dd", // TODO: Create contextId.
		Data:      data,
	}

	// Marshall the chs-delta previously created into the avro schema and convert it to a []byte for sending.
	messageBytes, err := chsDeltaAvro.Marshal(deltaData)
	if err != nil {
		return err
	}

	// Create the producer message which will contain a topic, our message and a default partition.
	producerMessage := &producer.Message{
		Topic:     topic,
		Value:     sarama.ByteEncoder(messageBytes),
		Partition: 0,
	}

	// Finally try to send the message.
	partition, offset, err := callSend(kSvc, producerMessage)
	log.Info("Sending message", log.Data{"topic": producerMessage.Topic, "partition": partition, "offset": offset})
	log.Trace("Sending message", log.Data{"message source": deltaData}) // TODO: Don't log entire data blob
	return err
}

// sendViaProducer is used to add an abstraction layer for unit testing when calling to send a message via a producer.
func sendViaProducer(k *KafkaServiceImpl, msg *producer.Message) (int32, int64, error) {
	return k.P.Send(msg)
}
