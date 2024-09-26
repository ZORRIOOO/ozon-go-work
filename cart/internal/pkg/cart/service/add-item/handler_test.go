package add_item

import (
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	lomsTypes "homework/cart/internal/client/api/loms/types"
	"homework/cart/internal/client/api/product/types"
	"homework/cart/internal/pkg/cart/model"
	"homework/cart/internal/pkg/cart/service/add-item/mock"
	"testing"
)

func TestHandler_AddItem(t *testing.T) {
	ctrl := minimock.NewController(t)
	cartRepositoryMock := mock.NewCartRepositoryMock(ctrl)
	productServiceMock := mock.NewProductServiceMock(ctrl)
	lomsServiceMock := mock.NewLomsServiceMock(ctrl)
	productToken := "testtoken"

	addItemHandler := NewHandler(cartRepositoryMock, productServiceMock, lomsServiceMock, productToken)

	skuId := int64(123)
	userId := int64(123)
	count := uint16(1)
	cartParams := model.CartParameters{
		SKU:    skuId,
		UserId: userId,
		Count:  count,
	}
	name := "Кроссовки 'Nike'"
	price := uint32(7500)

	productResponse := types.ProductResponse{
		Name:  name,
		Price: price,
	}
	request := types.ProductRequest{
		Sku:   skuId,
		Token: productToken,
	}
	productServiceMock.GetProductMock.Expect(request).Return(&productResponse, nil)

	stockResponse := lomsTypes.StocksInfoResponse{
		Count: 10,
	}
	stockRequest := lomsTypes.StocksInfoRequest{
		Sku: skuId,
	}
	lomsServiceMock.StocksInfoMock.Expect(stockRequest).Return(&stockResponse, nil)

	cartItem := model.CartItem{
		SKU:    skuId,
		Name:   name,
		Count:  count,
		Price:  price,
		UserId: userId,
	}
	cartRepositoryMock.AddItemMock.Expect(cartItem).Return(&cartItem, nil)

	actualResponse, err := addItemHandler.AddItem(cartParams)

	expectedResponse := &model.CartItem{
		SKU:    skuId,
		Name:   name,
		Count:  count,
		Price:  price,
		UserId: userId,
	}
	require.NoError(t, err)
	require.Equal(t, expectedResponse, actualResponse)

	productAfterCounter := productServiceMock.GetProductMock.Expect(request).Return(&productResponse, nil).GetProductAfterCounter()

	require.EqualValues(t, productAfterCounter, 1)
}
