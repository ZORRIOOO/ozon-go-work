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

	baseSKU := int64(1)
	name := "Кофеварка 'Kitfort'"
	count := uint16(1)
	price := uint32(7500)
	userId := int64(123)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cartItem := model.CartItem{
			SKU:    baseSKU + int64(i),
			Name:   fmt.Sprintf("%s №%d", name, i),
			Count:  count,
			Price:  price,
			UserId: userId,
		}

		_, err := repo.AddItem(cartItem)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}
