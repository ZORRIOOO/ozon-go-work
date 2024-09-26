package stock

type SKU = int64

type TotalCount = int32

type Stock struct {
	SKU        SKU        `json:"sku"`
	TotalCount TotalCount `json:"total_count"`
	Reserved   int32      `json:"reserved"`
}
