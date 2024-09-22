package stock

type SKU = int64

type Stock struct {
	SKU        SKU   `json:"sku"`
	TotalCount int32 `json:"total_count"`
	Reserved   int32 `json:"reserved"`
}
