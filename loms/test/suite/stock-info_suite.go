package suite

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"homework/loms/core/reader"
	"homework/loms/core/utils"
	"homework/loms/internal/repository/order"
	"homework/loms/internal/repository/stock"
	lomsService "homework/loms/internal/service/loms"
	"homework/loms/pkg/api/loms/v1"
)

type StockInfoSuite struct {
	suite.Suite
	service *lomsService.Service
}

func (s *StockInfoSuite) SetupSuite() {
	const (
		capacity = 1000
		filePath = "../../assets/stock-data.json"
	)

	stocks, err := reader.ReadStocks(utils.GetEnv("DOCKER_PATH_ASSETS", filePath))
	if err != nil {
		fmt.Sprintf("Read stocks failed: %v", err.Error())
	}
	orderRepository := order.NewRepository(capacity)
	stockRepository := stock.NewRepository(capacity, stocks)
	controller := lomsService.NewService(orderRepository, stockRepository)

	s.service = controller
}

func (s *StockInfoSuite) TestStockInfo() {
	ctx := context.Background()
	sku := int64(773297411)
	count := int32(150)

	stockInfoRequest := &loms.StocksInfoRequest{Sku: sku}
	stockInfoResponse, err := s.service.StocksInfo(ctx, stockInfoRequest)

	expectedResponse := &loms.StocksInfoResponse{Count: count}
	require.NoError(s.T(), err)
	require.Equal(s.T(), expectedResponse, stockInfoResponse)
}
