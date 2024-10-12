package flags

import (
	"flag"
	"fmt"
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

func init() {
	flag.IntVar(&FlagCLI.RepeatCnt, "repeat-count", 3, "count times all messages sent")
	flag.IntVar(&FlagCLI.StartID, "start-id", 1, "start order-id of all messages")
	flag.IntVar(&FlagCLI.Count, "count", 1, "count of orders to emit events")
	flag.StringVar(&FlagCLI.Topic, "topic", "loms.order-events", "topic to produce")
	flag.DurationVar(&FlagCLI.Interval, "interval", DefaultInterval, "duration between messages")

	flag.Parse()
	fmt.Print("Init broker flags")
}
