package types

type ProductRequest struct {
	Sku   int64  `json:"sku"`
	Token string `json:"token"`
}

type ProductResponse struct {
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}

type SkusRequest struct {
	Token         string `json:"token"`
	StartAfterSku uint32 `json:"startAfterSku"`
	Count         uint32 `json:"count"`
}

type SkusResponse struct {
	Skus []uint32 `json:"skus"`
}
