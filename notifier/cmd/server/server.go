package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/IBM/sarama"
	"homework/notifier/internal/infra/kafka/broker/config"
	consumerGroup "homework/notifier/internal/infra/kafka/broker/consumer-group"
	"homework/notifier/internal/infra/kafka/broker/flags"
	consumerGroupSarama "homework/notifier/internal/infra/kafka/sarama/consumer-group"
	"homework/notifier/internal/infra/kafka/sarama/handler"
	"homework/notifier/internal/infra/kafka/sarama/options"
	"log"
	"sync"
)

var CLIFlags = flags.Flags{}

func init() {
	flag.StringVar(&CLIFlags.Topic, "topic", "loms.order-events", "topic to produce")
	flag.StringVar(&CLIFlags.BootstrapServer, "bootstrap-server", "broker:29092", "kafka broker host and port")
	flag.StringVar(&CLIFlags.ConsumerGroupName, "cg-name", "hw-consumer-group", "topic to produce")

	flag.Parse()
	fmt.Print("Init broker flags")
}

func main() {
	var (
		wg   = &sync.WaitGroup{}
		conf = config.NewConfig(CLIFlags)
		ctx  = consumerGroup.RunSignalHandler(context.Background(), wg)
	)

	consumerGroupHandler := handler.NewConsumerGroupHandler()
	cg, err := consumerGroupSarama.NewConsumerGroup(
		conf.KafkaConfig.Brokers,
		CLIFlags.ConsumerGroupName,
		[]string{CLIFlags.Topic},
		consumerGroupHandler,
		options.WithOffsetsInitial(sarama.OffsetOldest),
		options.WithReturnSuccessesEnabled(true),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer cg.Close()

	consumerGroup.RunCGErrorHandler(ctx, cg, wg)

	cg.Run(ctx, wg)
	wg.Wait()
}
