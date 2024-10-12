package flags

import (
	"flag"
	"time"
)

type Flags struct {
	RepeatCnt int
	StartID   int
	Count     int
	Topic     string
	Interval  time.Duration
}

var (
	FlagCLI         = Flags{}
	DefaultInterval = 100 * time.Millisecond
)

func Init() {
	flag.IntVar(&FlagCLI.RepeatCnt, "repeat-count", 3, "count times all messages sent")
	flag.IntVar(&FlagCLI.StartID, "start-id", 1, "start order-id of all messages")
	flag.IntVar(&FlagCLI.Count, "count", 10, "count of orders to emit events")
	flag.StringVar(&FlagCLI.Topic, "topic", "route256-example", "topic to produce")
	flag.DurationVar(&FlagCLI.Interval, "interval", DefaultInterval, "duration between messages")

	flag.Parse()
}
