package tests

import (
	"context"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
	orderModel "homework/loms/internal/model/order"
	"homework/loms/internal/service/loms"
	"homework/loms/internal/service/loms/mocks"
	lomsApi "homework/loms/pkg/api/loms/v1"
	"testing"
)

func Test_OrderCancel(t *testing.T) {
	ctx := context.Background()
	ctrl := minimock.NewController(t)
	orderRepositoryMock := mocks.NewOrderRepositoryMock(ctrl)
	stocksRepositoryMock := mocks.NewStocksRepositoryMock(ctrl)
	lomsService := loms.NewService(orderRepositoryMock, stocksRepositoryMock)

	orderId := int64(1)
	user := int64(1)
	sku := int64(773297411)
	count := int32(10)
	status := "awaiting payment"
	statusCancelled := "cancelled"
	items := make([]*lomsApi.Item, 0)
	items = append(items, &lomsApi.Item{
		Sku:   sku,
		Count: count,
	})
	order := &orderModel.Order{
		OrderId: orderId,
		Status:  status,
		User:    user,
		Items:   items,
	}

	orderRepositoryMock.GetByIdMock.Expect(ctx, orderId).Return(order, nil)
	stocksRepositoryMock.ReserveCancelMock.Expect(ctx, order).Return(nil)
	orderRepositoryMock.SetStatusMock.Expect(ctx, orderId, statusCancelled).Return(nil)

	request := &lomsApi.OrderCancelRequest{OrderId: orderId}
	actualResponse, err := lomsService.OrderCancel(ctx, request)
	expectedResponse := &emptypb.Empty{}

	require.NoError(t, err)
	require.Equal(t, expectedResponse, actualResponse)
}
