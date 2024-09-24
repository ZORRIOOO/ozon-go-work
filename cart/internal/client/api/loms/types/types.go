package types

type OrderId = int64

type OrderCreateRequest struct {
	User  int64  `json:"user"`
	Items []Item `json:"items"`
}

type Item struct {
	Sku   int64  `json:"sku"`
	Count uint16 `json:"count"`
}

type OrderCreateResponse struct {
	OrderId OrderId `json:"order_id"`
}

type StocksInfoRequest struct {
	Sku int64 `json:"sku"`
}

type StocksInfoResponse struct {
	Count uint16 `json:"count"`
}
