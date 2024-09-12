package delete_item

import (
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"homework/cart/internal/pkg/cart/model"
	"homework/cart/internal/pkg/cart/service/delete-item/mock"
	"testing"
)

func TestHandler_DeleteItem(t *testing.T) {
	ctrl := minimock.NewController(t)
	cartRepositoryMock := mock.NewCartRepositoryMock(ctrl)

	deleteItemHandler := NewHandler(cartRepositoryMock)

	skuId := int64(123)
	userId := int64(123)
	count := uint16(3)
	name := "Тетрадь (В клетку)"
	price := uint32(50)

	deleteCartParameters := model.DeleteCartParameters{
		SKU:    skuId,
		UserId: userId,
	}
	cartItem := model.CartItem{
		SKU:    skuId,
		Name:   name,
		Count:  count,
		Price:  price,
		UserId: userId,
	}

	cartRepositoryMock.DeleteItemMock.Expect(deleteCartParameters).Return(&cartItem, nil)

	cartParams := model.DeleteCartParameters{
		SKU:    skuId,
		UserId: userId,
	}
	actualResponse, err := deleteItemHandler.DeleteItem(cartParams)

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
