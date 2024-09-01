package model

type UserId = int64

type CartItem struct {
	SKU    int64
	Name   string
	Count  uint16
	Price  uint32
	UserId UserId
}
