package tests

import (
	"context"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	orderModel "homework/loms/internal/model/order"
	"homework/loms/internal/service/loms"
	"homework/loms/internal/service/loms/mocks"
	lomsApi "homework/loms/pkg/api/loms/v1"
	"testing"
)

func Test_OrderCreate(t *testing.T) {
	ctx := context.Background()
	ctrl := minimock.NewController(t)
	orderRepositoryMock := mocks.NewOrderRepositoryMock(ctrl)
	stocksRepositoryMock := mocks.NewStocksRepositoryMock(ctrl)
	lomsService := loms.NewService(orderRepositoryMock, stocksRepositoryMock)

	user := int64(1)
	sku := int64(773297411)
	count := int32(1)
	orderId := int64(1)
	status := "new"
	statusAwaiting := "awaiting payment"
	items := make([]*lomsApi.Item, 0)
	items = append(items, &lomsApi.Item{
		Sku:   sku,
		Count: count,
	})
	order := orderModel.Order{
		Status: status,
		User:   user,
		Items:  items,
	}
	newOrder := orderModel.Order{
		OrderId: orderId,
		Status:  status,
		User:    user,
		Items:   items,
	}

	orderRepositoryMock.CreateMock.Expect(ctx, order).Return(orderId, nil)
	stocksRepositoryMock.ReserveMock.Expect(ctx, newOrder).Return(nil)
	orderRepositoryMock.SetStatusMock.Expect(ctx, orderId, statusAwaiting).Return(nil)

	request := &lomsApi.OrderCreateRequest{
		User:  user,
		Items: items,
	}
	actualResponse, err := lomsService.OrderCreate(ctx, request)
	expectedResponse := &lomsApi.OrderCreateResponse{OrderId: orderId}

	require.NoError(t, err)
	require.Equal(t, expectedResponse, actualResponse)
}
