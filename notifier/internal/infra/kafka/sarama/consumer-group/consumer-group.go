package consumer_group

import (
	"context"
	"homework/notifier/internal/infra/kafka/sarama/options"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
)

type ConsumerGroup struct {
	sarama.ConsumerGroup
	Handler sarama.ConsumerGroupHandler
	Topics  []string
}

func (c *ConsumerGroup) Run(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("[consumer-group] run")

		for {
			if err := c.ConsumerGroup.Consume(ctx, c.Topics, c.Handler); err != nil {
				log.Printf("Error from consume: %v\n", err)
			}
			if ctx.Err() != nil {
				log.Printf("[consumer-group]: ctx closed: %s\n", ctx.Err().Error())
				return
			}
		}
	}()

}

func NewConsumerGroup(brokers []string, groupID string, topics []string, consumerGroupHandler sarama.ConsumerGroupHandler, opts ...options.Option) (*ConsumerGroup, error) {
	config := sarama.NewConfig()
	config.Version = sarama.MaxVersion
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Group.ResetInvalidOffsets = true
	config.Consumer.Group.Heartbeat.Interval = 3 * time.Second
	config.Consumer.Group.Session.Timeout = 60 * time.Second
	config.Consumer.Group.Rebalance.Timeout = 60 * time.Second
	config.Consumer.Return.Errors = true

	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second

	for _, opt := range opts {
		opt.Apply(config)
	}

	cg, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}

	return &ConsumerGroup{
		ConsumerGroup: cg,
		Handler:       consumerGroupHandler,
		Topics:        topics,
	}, nil
}
