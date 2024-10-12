package producer

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"homework/loms/internal/infra/kafka/broker/config"
	"homework/loms/internal/infra/kafka/broker/flags"
	"homework/loms/internal/infra/kafka/sarama/options"
	"homework/loms/internal/infra/kafka/sarama/producer"
	"homework/loms/internal/infra/kafka/types"
	"homework/loms/internal/model/order"
	"log"
	"strconv"
	"time"
)

const (
	defaultAcks  = sarama.WaitForAll
	openRequests = 1
	retries      = 5
	duration     = 10 * time.Millisecond
)

type KafkaProducer struct {
	syncProducer  sarama.SyncProducer
	configuration config.Config
}

func NewKafkaProducer(brokerAddr string) *KafkaProducer {
	configuration := config.NewConfig(brokerAddr, flags.FlagCLI)
	syncProducer, err := producer.NewSyncProducer(configuration.Kafka,
		options.WithIdempotent(),
		options.WithRequiredAcks(defaultAcks),
		options.WithMaxOpenRequests(openRequests),
		options.WithMaxRetries(retries),
		options.WithRetryBackoff(duration),
		options.WithProducerPartitioner(sarama.NewManualPartitioner),
		options.WithProducerPartitioner(sarama.NewRoundRobinPartitioner),
		options.WithProducerPartitioner(sarama.NewRandomPartitioner),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer syncProducer.Close()

	return &KafkaProducer{
		syncProducer:  syncProducer,
		configuration: configuration,
	}
}

func (k *KafkaProducer) SendMessage(messagePayload types.MessagePayload) error {
	bytes, err := json.Marshal(messagePayload)
	if err != nil {
		log.Printf("Failed to marshal message payload: %v", err)
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: k.configuration.Producer.Topic,
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

	partition, offset, err := k.syncProducer.SendMessage(msg)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		return err
	}

	log.Printf("Message sent to partition %d with offset %d", partition, offset)
	return nil
}

func RepackPayload(id order.Id, eventStatus order.Status) types.MessagePayload {
	return types.MessagePayload{
		OrderId:  id,
		Event:    eventStatus,
		Datetime: time.Now().UTC().String(),
	}
}
