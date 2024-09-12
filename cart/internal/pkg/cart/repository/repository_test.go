package repository

import (
	"github.com/stretchr/testify/require"
	"homework/cart/internal/pkg/cart/model"
	"testing"
)

func TestHandler_AddItem(t *testing.T) {
	capacity := 100
	cartRepository := NewCartRepository(capacity)

	skuId := int64(123)
	userId := int64(123)
	count := uint16(1)
	name := "Куртка 'UNIQLO'"
	price := uint32(5000)
	item := model.CartItem{
		SKU:    skuId,
		Name:   name,
		Count:  count,
		Price:  price,
		UserId: userId,
	}
	actualResponse, err := cartRepository.AddItem(item)

	expectedResponse := &model.CartItem{
		SKU:    skuId,
		Name:   name,
		Count:  count,
		Price:  price,
		UserId: userId,
	}
	require.Equal(t, expectedResponse, actualResponse)
	require.NoError(t, err)
}

func TestHandler_DeleteItem(t *testing.T) {

}

func TestHandler_DeleteItemsByUser(t *testing.T) {

}

func TestHandler_GetItemsByUser(t *testing.T) {

}
