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

func Test_OrderInfo(t *testing.T) {
	ctx := context.Background()
	ctrl := minimock.NewController(t)
	orderRepositoryMock := mocks.NewOrderRepositoryMock(ctrl)
	stocksRepositoryMock := mocks.NewStocksRepositoryMock(ctrl)
	lomsService := loms.NewService(orderRepositoryMock, stocksRepositoryMock)

	orderId := int64(1)
	user := int64(1)
	sku := int64(773297411)
	status := "awaiting payment"
	count := int32(10)
	items := []*lomsApi.Item{{
		Sku:   sku,
		Count: count,
	}}
	order := &orderModel.Order{
		OrderId: orderId,
		Status:  status,
		User:    user,
		Items:   items,
	}

	orderRepositoryMock.GetByIdMock.Expect(ctx, orderId).Return(order, nil)

	request := &lomsApi.OrderInfoRequest{OrderId: orderId}
	actualResponse, err := lomsService.OrderInfo(ctx, request)
	expectedResponse := &lomsApi.OrderInfoResponse{
		User:   user,
		Status: status,
		Items:  items,
	}

	require.NoError(t, err)
	require.Equal(t, expectedResponse, actualResponse)
}
