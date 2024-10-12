package producer

import (
	"fmt"
	"github.com/IBM/sarama"
	"homework/loms/internal/infra/kafka/sarama/config"
	"homework/loms/internal/infra/kafka/sarama/options"
	"homework/loms/internal/infra/kafka/types"
)

func NewSyncProducer(conf types.Config, opts ...options.Option) (sarama.SyncProducer, error) {
	syncProducer, err := sarama.NewSyncProducer(conf.Brokers, config.PrepareConfig(opts...))
	if err != nil {
		return nil, fmt.Errorf("NewSyncProducer failed: %w", err)
	}

	return syncProducer, nil
}
