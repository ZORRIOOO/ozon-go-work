package cart_checkout

import (
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"homework/cart/internal/client/api/loms/types"
	"homework/cart/internal/pkg/cart/model"
	"homework/cart/internal/pkg/cart/service/cart-checkout/mock"
	"testing"
)

func TestHandler_CartCheckout(t *testing.T) {
	ctrl := minimock.NewController(t)
	cartRepositoryMock := mock.NewCartRepositoryMock(ctrl)
	lomsServiceMock := mock.NewLomsServiceMock(ctrl)

	cartCheckoutHandler := NewHandler(cartRepositoryMock, lomsServiceMock)

	userId := int64(1)
	orderId := int64(1)
	sku := int64(773297411)
	count := uint16(1)
	items := []model.CartItem{{
		SKU:    sku,
		Name:   "Кроссовки",
		Count:  count,
		Price:  uint32(4500),
		UserId: userId,
	}}
	cartRepositoryMock.GetItemsByUserMock.Expect(userId).Return(items, nil)

	requestItems := []types.Item{{
		Sku:   sku,
		Count: count,
	}}
	request := types.OrderCreateRequest{
		User:  userId,
		Items: requestItems,
	}
	response := &types.OrderCreateResponse{OrderId: orderId}
	lomsServiceMock.CreateOrderMock.Expect(request).Return(response, nil)

	cartRepositoryMock.DeleteItemsByUserMock.Expect(userId).Return(&userId, nil)

	actualResponse, err := cartCheckoutHandler.CartCheckout(userId)
	expectedResponse := &model.Checkout{
		OrderId: orderId,
	}

	require.NoError(t, err)
	require.Equal(t, expectedResponse, actualResponse)
}
