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

type KafkaService interface {
	SendMessage(topic , data string) error
}

type KafkaServiceImpl struct {
	Schema   string
	Producer *producer.Producer
}

const (
	schemaName = "chs-delta"
)
func NewKafkaService(cfg *config.Config) (KafkaServiceImpl, error) {

	log.Info("Get schema from Avro", log.Data{"schema_name": schemaName})
	sch, err := schema.Get(cfg.SchemaRegistryURL, schemaName)
	if err != nil {
		log.Error(fmt.Errorf("error receiving %s schema: %s", schemaName, err))
		return KafkaServiceImpl{}, err
	}
	log.Info("Successfully received schema", log.Data{"schema_name": schemaName})

	log.Info("Using Streaming Kafka broker Address", log.Data{"Brokers": cfg.BrokerAddr})
	p, err := producer.New(&producer.Config{Acks: &producer.WaitForAll, BrokerAddrs: cfg.BrokerAddr})
	if err != nil {
		log.Error(fmt.Errorf("error initialising producer: %s", err))
		return KafkaServiceImpl{}, err
	}

	return KafkaServiceImpl{
		Schema: sch,
		Producer: p,
	}, nil
}

func (kSvc *KafkaServiceImpl) SendMessage(topic, data string) error {
	chsDeltaAvro := &avro.Schema{
		Definition: kSvc.Schema,
	}

	deltaData := models.ChsDelta{
		ContextId: "0",
		Data: data,
		Attempt: 1,
	}
	messageBytes, err := chsDeltaAvro.Marshal(deltaData)
	if err != nil {
		return err
	}

	producerMessage := &producer.Message{
		Topic:     topic,
		Value:     sarama.ByteEncoder(messageBytes),
		Partition: 0,
	}

	partition, offset, err := kSvc.Producer.Send(producerMessage)
	log.Info("Start send message", log.Data{"topic": producerMessage.Topic, "partition": partition, "offset": offset})
	log.Trace("Start send message", log.Data{"message source": deltaData})

	if err != nil {
		fmt.Println("Error sending message")
	}
	log.Info("End send message", log.Data{"topic": producerMessage.Topic, "partition": partition, "offset": offset})
	return err
}