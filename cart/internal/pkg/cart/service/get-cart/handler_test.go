package get_cart

import (
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"homework/cart/internal/client/api/product/types"
	"homework/cart/internal/pkg/cart/model"
	"homework/cart/internal/pkg/cart/service/get-cart/mock"
	"testing"
)

func TestHandler_GetCartByUser(t *testing.T) {
	ctrl := minimock.NewController(t)
	cartRepositoryMock := mock.NewCartRepositoryMock(ctrl)
	productServiceMock := mock.NewProductServiceMock(ctrl)
	productToken := "testtoken"

	getCartByUserHandler := NewHandler(cartRepositoryMock, productServiceMock, productToken)

	userId := int64(123)

	items := make([]model.CartItem, 0)
	items = append(items, model.CartItem{
		SKU:    int64(123),
		Name:   "Стейк говяжий 'Мираторг'",
		Count:  uint16(1),
		Price:  uint32(500),
		UserId: userId,
	})
	cartRepositoryMock.GetItemsByUserMock.Expect(userId).Return(items, nil)

	productServiceMock.GetProductMock.Expect(types.ProductRequest{
		Sku:   int64(123),
		Token: productToken,
	}).Return(&types.ProductResponse{
		Name:  "Стейк говяжий 'Мираторг'",
		Price: uint32(500),
	}, nil)

	actualResponse, err := getCartByUserHandler.GetCartByUser(userId)

	expectedItems := make([]model.CartItem, 0)
	expectedItems = append(expectedItems, model.CartItem{
		SKU:    int64(123),
		Name:   "Стейк говяжий 'Мираторг'",
		Count:  uint16(1),
		Price:  uint32(500),
		UserId: userId,
	})
	expectedResponse := &model.Cart{
		Items:      expectedItems,
		TotalPrice: uint32(500),
	}
	require.Equal(t, expectedResponse, actualResponse)
	require.NoError(t, err)
}
