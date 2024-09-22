package order

import (
	"homework/loms/pkg/api/loms/v1"
)

type User = int64

type Id = int64

type Status = string

type Order struct {
	OrderId Id
	Status  Status
	User    User
	Items   []*loms.Item
}
