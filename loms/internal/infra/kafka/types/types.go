package types

type MessagePayload struct {
	OrderId  int64  `json:"order_id"`
	Event    string `json:"event"`
	Datetime string `json:"datetime"`
	Extra    string `json:"extra"`
}

type Config struct {
	Brokers []string
}
