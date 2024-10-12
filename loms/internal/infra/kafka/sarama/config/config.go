package config

import (
	"github.com/IBM/sarama"
	"homework/loms/internal/infra/kafka/sarama/options"
	"time"
)

func PrepareConfig(opts ...options.Option) *sarama.Config {
	c := sarama.NewConfig()

	{
		c.Producer.Partitioner = sarama.NewHashPartitioner
	}

	{
		c.Producer.RequiredAcks = sarama.WaitForAll
	}

	{
		c.Producer.Idempotent = false
	}

	{
		c.Producer.Retry.Max = 100
		c.Producer.Retry.Backoff = 5 * time.Millisecond
	}

	{
		c.Net.MaxOpenRequests = 1
	}

	{
		c.Producer.CompressionLevel = sarama.CompressionLevelDefault
		c.Producer.Compression = sarama.CompressionGZIP
	}

	{
		c.Producer.Return.Successes = true
		c.Producer.Return.Errors = true
	}

	for _, opt := range opts {
		_ = opt.Apply(c)
	}

	return c
}
