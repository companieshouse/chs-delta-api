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
	schemaName = "chs-delta"
)

// KafkaService defines all Methods needed to successfully send a message onto a Kafka topic.
type KafkaService interface {
	Init(cfg *config.Config) error
	SendMessage(topic, data string) error
}

// KafkaServiceImpl is a concrete implementation of the KafkaService interface.
type KafkaServiceImpl struct {
	Schema   string
	Producer *producer.Producer
}

// NewKafkaService constructs and returns a KafkaServiceImpl using a provided Config.
func NewKafkaService() KafkaServiceImpl {
	return KafkaServiceImpl{}
}

func (kSvc *KafkaServiceImpl) Init(cfg *config.Config) error {
	// Retrieve the generic chs-delta avro schema.
	log.Trace("Get schema from Avro", log.Data{"schema_name": schemaName})
	sch, err := schema.Get(cfg.SchemaRegistryURL, schemaName)
	if err != nil {
		log.Error(fmt.Errorf("error receiving %s schema: %s", schemaName, err))
		return err
	}
	log.Trace("Successfully received schema", log.Data{"schema_name": schemaName})

	// Create a new Kafka Producer which will be used to publish our message onto a given Kafka topic.
	log.Trace("Using Streaming Kafka broker Address", log.Data{"Brokers": cfg.BrokerAddr})
	p, err := producer.New(&producer.Config{Acks: &producer.WaitForAll, BrokerAddrs: cfg.BrokerAddr})
	if err != nil {
		log.Error(fmt.Errorf("error initialising producer: %s", err))
		return err
	}

	kSvc.Schema = sch
	kSvc.Producer = p

	return nil
}

// SendMessage publishes a given data string retrieved from a REST request onto a chosen Kafka topic.
func (kSvc *KafkaServiceImpl) SendMessage(topic, data string) error {

	// Retrieve our chs-delta avro schema using the chs go avro package.
	chsDeltaAvro := &avro.Schema{
		Definition: kSvc.Schema,
	}

	// Construct a chs-delta using provided data.
	deltaData := models.ChsDelta{
		ContextId: "uu-aa-dd",
		DeltaAt: 2014, // TODO: Should these be removed from the avro schema?
		CreatedAt: "20140925104551", // TODO: Should these be removed from the avro schema?
		Data: data,
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
	partition, offset, err := kSvc.Producer.Send(producerMessage)
	log.Info("Sending message", log.Data{"topic": producerMessage.Topic, "partition": partition, "offset": offset})
	log.Trace("Sending message", log.Data{"message source": deltaData}) // TODO: Temp, we need to find a way not to output sensitive data.
	return err
}