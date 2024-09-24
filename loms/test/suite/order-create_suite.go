package suite

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"homework/loms/internal/repository/order"
	"homework/loms/internal/repository/stock"
	lomsService "homework/loms/internal/service/loms"
	"homework/loms/pkg/api/loms/v1"
)

type OrderCreateSuite struct {
	suite.Suite
	service *lomsService.Service
}

func (s *OrderCreateSuite) SetupSuite() {
	const (
		capacity = 1000
		filePath = "../../assets/stock-data.json"
	)

	orderRepository := order.NewRepository(capacity)
	stockRepository := stock.NewRepository(capacity, filePath)
	controller := lomsService.NewService(orderRepository, stockRepository)

	s.service = controller
}

func (s *OrderCreateSuite) TestOrderCreate() {
	ctx := context.Background()
	sku := int64(773297411)
	count := int32(1)
	user := int64(1)
	items := make([]*loms.Item, 0)
	status := "awaiting payment"
	items = append(items, &loms.Item{
		Sku:   sku,
		Count: count,
	})

	orderCreateRequest := &loms.OrderCreateRequest{User: user, Items: items}
	orderCreateResponse, err := s.service.OrderCreate(ctx, orderCreateRequest)

	require.NoError(s.T(), err)

	orderInfoRequest := &loms.OrderInfoRequest{OrderId: orderCreateResponse.OrderId}
	orderInfoResponse, err := s.service.OrderInfo(ctx, orderInfoRequest)

	expectedResponse := &loms.OrderInfoResponse{Status: status, User: user, Items: items}
	require.NoError(s.T(), err)
	require.Equal(s.T(), expectedResponse, orderInfoResponse)
}
