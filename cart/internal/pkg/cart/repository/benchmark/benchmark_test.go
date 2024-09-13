package benchmark

import (
	"fmt"
	"homework/cart/internal/pkg/cart/model"
	"homework/cart/internal/pkg/cart/repository"
	"testing"
)

func BenchmarkCartRepository_AddItem(b *testing.B) {
	capacity := 100
	repo := repository.NewCartRepository(capacity)

	skuId := int64(12345)
	name := "Кофеварка 'Kitfort'"
	count := uint16(1)
	price := uint32(7500)
	userId := int64(123)

	cartItem := model.CartItem{
		SKU:    skuId,
		Name:   name,
		Count:  count,
		Price:  price,
		UserId: userId,
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := repo.AddItem(cartItem)
		if err != nil {
			message := fmt.Sprintf("Unexpected error: %v", err)
			b.Fatalf(message)
		}
	}
}
