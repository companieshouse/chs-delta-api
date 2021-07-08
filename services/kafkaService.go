package services

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/companieshouse/chs-delta-api/config"
	"github.com/companieshouse/chs.go/avro"
	"github.com/companieshouse/chs.go/avro/schema"
	"github.com/companieshouse/chs.go/kafka/producer"
	"github.com/companieshouse/chs.go/log"
)

type KafkaService interface {
	func SendMessage(topic , data string) error
}

type KafkaServiceImpl struct {

	Schema string
}

const (
	schemaName       = "chs-delta"
)
func NewKafkaService(cfg *config.Config) ( KafkaServiceImpl, error) {

	log.Info("Get schema from Avro", log.Data{"schema_name": schemaName})
	schema, err := schema.Get(cfg.SchemaRegistryURL, schemaName)
	if err != nil {
		log.Error(fmt.Errorf("error receiving %s schema: %s", schemaName, err))
		return nil, err
	}
	log.Info("Successfully received schema", log.Data{"schema_name": schemaName})
	return KafkaServiceImpl{}
}

func (svc *KafkaServiceImpl) SendMessage(topic, data string) error {
	chsDeltaAvro := &avro.Schema{
		Definition: svc.ProducerSchema,
	}

	//Add avro schema model here when created

	messageBytes, err := callAvroProducerMarshall(producerAvro, StreamFilingHistoryData)
	if err != nil {
		return err
	}

	producerMessage := &producer.Message{
		Topic:     topic,
		Value:     sarama.ByteEncoder(messageBytes),
		Partition: evaluatePartition(i.ResourceKind),
	}

	partition, offset, err := callDoPublishMessage(svc, producerMessage)
	log.InfoC(i.ContextID, "Start send message", log.Data{"topic": producerMessage.Topic, "partition": partition, "offset": offset})
	log.TraceC(i.ContextID, "Start send message", log.Data{"message source": i})

	if err != nil {
		fmt.Println("Error sending message")
	}
	log.InfoC(i.ContextID, "End send message", log.Data{"topic": producerMessage.Topic, "partition": partition, "offset": offset})
	return err
}