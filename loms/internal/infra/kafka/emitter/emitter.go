package emitter

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"homework/loms/internal/infra/kafka/broker/producer"
	"homework/loms/internal/infra/kafka/types"
	"log"
	"strconv"
	"time"
)

const Topic = "loms.order-events"

type Emitter struct {
	syncProducer sarama.SyncProducer
}

func NewEmitter(brokerAddr string) *Emitter {
	kafkaProducer := producer.NewKafkaProducer(brokerAddr)
	return &Emitter{syncProducer: kafkaProducer}
}

func (e *Emitter) SendMessage(messagePayload types.MessagePayload) error {
	bytes, err := json.Marshal(messagePayload)
	if err != nil {
		log.Printf("Failed to marshal message payload: %v", err)
		return err
	}

	message := &sarama.ProducerMessage{
		Topic: Topic,
		Key:   sarama.StringEncoder(strconv.FormatInt(messagePayload.OrderId, 10)),
		Value: sarama.ByteEncoder(bytes),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("app-name"),
				Value: []byte("loms-sync-prod"),
			},
		},
		Timestamp: time.Now(),
	}

	partition, offset, err := e.syncProducer.SendMessage(message)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		return err
	}

	log.Printf("Message sent to partition %d with offset %d", partition, offset)
	return nil
}
