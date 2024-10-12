package config

import (
	"homework/loms/internal/infra/kafka/broker/flags"
	"homework/loms/internal/infra/kafka/types"
	"time"
)

type (
	AppConfig struct {
		RepeatCnt int           `json:"repeat_cnt"`
		StartID   int           `json:"start_id"`
		Count     int           `json:"count"`
		Interval  time.Duration `json:"interval"`
	}

	ProducerConfig struct {
		Topic string `json:"topic"`
	}

	Config struct {
		App      AppConfig      `json:"app"`
		Kafka    types.Config   `json:"kafka"`
		Producer ProducerConfig `json:"producer"`
	}
)

func NewConfig(brokerAddr string, f flags.Flags) Config {
	return Config{
		App: AppConfig{
			RepeatCnt: f.RepeatCnt,
			StartID:   f.StartID,
			Count:     f.Count,
			Interval:  f.Interval,
		},
		Kafka: types.Config{
			Brokers: []string{
				brokerAddr,
			},
		},
		Producer: ProducerConfig{
			Topic: f.Topic,
		},
	}
}
