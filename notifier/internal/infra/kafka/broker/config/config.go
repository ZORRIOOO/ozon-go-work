package config

import (
	"homework/notifier/internal/infra/kafka/broker/flags"
	kafka "homework/notifier/internal/infra/kafka/sarama/config"
)

type Config struct {
	KafkaConfig kafka.Config
}

func NewConfig(f flags.Flags) Config {
	return Config{
		KafkaConfig: kafka.Config{
			Brokers: []string{
				f.BootstrapServer,
			},
		},
	}
}
