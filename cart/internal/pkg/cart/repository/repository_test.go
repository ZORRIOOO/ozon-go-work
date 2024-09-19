package repository

import (
	"github.com/stretchr/testify/require"
	"homework/cart/internal/pkg/cart/model"
	"testing"
)

func TestHandler_AddItem(t *testing.T) {
	skuId := int64(123)
	userId := int64(123)
	count := uint16(1)
	name := "Куртка 'UNIQLO'"
	price := uint32(5000)

	capacity := 100
	cartRepository := NewCartRepository(capacity)

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
	require.NoError(t, err)
	require.Equal(t, expectedResponse, actualResponse)
}

func TestHandler_DeleteItem(t *testing.T) {
	skuId := int64(123)
	userId := int64(123)
	count := uint16(1)
	name := "Книга 'Три мушкетера'"
	price := uint32(500)

	capacity := 100
	cartRepository := NewCartRepository(capacity)

	item := model.CartItem{
		SKU:    skuId,
		Name:   name,
		Count:  count,
		Price:  price,
		UserId: userId,
	}
	cartItem, err := cartRepository.AddItem(item)

	cartParams := model.DeleteCartParameters{
		SKU:    cartItem.SKU,
		UserId: cartItem.UserId,
	}
	actualResponse, err := cartRepository.DeleteItem(cartParams)

	expectedResponse := &model.CartItem{
		SKU:    skuId,
		Name:   name,
		Count:  count,
		Price:  price,
		UserId: userId,
	}

	require.NoError(t, err)
	require.Equal(t, expectedResponse, actualResponse)
}

func TestHandler_DeleteItemsByUser(t *testing.T) {
	skuId := int64(123)
	userId := int64(123)
	count := uint16(5)
	name := "Набор тарелок 'Fargklar'"
	price := uint32(1500)

	capacity := 100
	cartRepository := NewCartRepository(capacity)

	item := model.CartItem{
		SKU:    skuId,
		Name:   name,
		Count:  count,
		Price:  price,
		UserId: userId,
	}
	cartItem, err := cartRepository.AddItem(item)

	actualResponse, err := cartRepository.DeleteItemsByUser(cartItem.UserId)

	require.NoError(t, err)
	require.Equal(t, &userId, actualResponse)
}

func TestHandler_GetItemsByUser(t *testing.T) {
	userId := int64(123)

	capacity := 100
	cartRepository := NewCartRepository(capacity)

	firstItem := model.CartItem{
		SKU:    1234,
		Name:   "Чайник 'Xiaomi'",
		Count:  1,
		Price:  3000,
		UserId: userId,
	}
	cartRepository.AddItem(firstItem)

	secondItem := model.CartItem{
		SKU:    12345,
		Name:   "Фен 'Bosh'",
		Count:  1,
		Price:  6000,
		UserId: userId,
	}
	cartRepository.AddItem(secondItem)

	actualResponse, err := cartRepository.GetItemsByUser(userId)

	expectedResponse := make([]model.CartItem, 0)
	expectedResponse = append(expectedResponse, firstItem)
	expectedResponse = append(expectedResponse, secondItem)

	require.NoError(t, err)
	require.Equal(t, expectedResponse, actualResponse)
}
