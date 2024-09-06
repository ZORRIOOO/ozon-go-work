package model

type UserId = int64

type SKU = int64

type DeleteCartParameters struct {
	SKU    int64 `json:"sku_id" validate:"required,min=1"`
	UserId int64 `json:"user_id" validate:"required,min=1"`
}

type CartParameters struct {
	SKU    int64  `json:"sku_id" validate:"required,min=1"`
	UserId int64  `json:"user_id" validate:"required,min=1"`
	Count  uint16 `json:"count" validate:"required,min=1"`
}

type CartItem struct {
	SKU    SKU
	Name   string
	Count  uint16
	Price  uint32
	UserId UserId
}

type Cart struct {
	Items      []CartItem `json:"items"`
	TotalPrice uint32     `json:"total_price"`
}
