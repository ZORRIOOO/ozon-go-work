package order

type User = int64

type Id = int64

type Status = string

type Order struct {
	OrderId Id
	Status  Status
	User    User
	Items   []Item
}

type Item struct {
	Sku   int64
	Count int32
}
