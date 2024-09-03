package model

type UserId = int64

type SKU = int64

type CartParameters struct {
	SKU    SKU
	UserId UserId
	Count  uint16
}

type CartItem struct {
	SKU    SKU
	Name   string
	Count  uint16
	Price  uint32
	UserId UserId
}
