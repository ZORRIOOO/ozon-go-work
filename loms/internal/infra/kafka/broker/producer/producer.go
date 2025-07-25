package producer

import (
	"github.com/IBM/sarama"
	"homework/loms/internal/infra/kafka/broker/config"
	"homework/loms/internal/infra/kafka/broker/flags"
	"homework/loms/internal/infra/kafka/sarama/options"
	"homework/loms/internal/infra/kafka/sarama/producer"
	"log"
	"time"
)

const (
	defaultAcks  = sarama.WaitForAll
	openRequests = 1
	retries      = 5
	duration     = 10 * time.Millisecond
)

func NewKafkaProducer(brokerAddr string) sarama.SyncProducer {
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

	return syncProducer
}
