package suite

import (
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	lomsService "homework/cart/internal/client/api/loms/service"
	productService "homework/cart/internal/client/api/product/service"
	httpclient "homework/cart/internal/client/base/client"
	"homework/cart/internal/pkg/cart/model"
	"homework/cart/internal/pkg/cart/repository"
	addItem "homework/cart/internal/pkg/cart/service/add-item"
	deleteItem "homework/cart/internal/pkg/cart/service/delete-item"
	"time"
)

type DeleteCartItemSuite struct {
	suite.Suite
	deleteCartItemHandler *deleteItem.CartServiceHandler
	addCartItemHandler    *addItem.CartServiceHandler
}

func (s *DeleteCartItemSuite) SetupSuite() {
	const (
		productAddress = "http://route256.pavl.uk:8080"
		productToken   = "testtoken"
		lomsAddress    = "http://localhost:8081"
	)
	client := httpclient.NewHttpClient(10*time.Second, 3, []int{420, 429})
	cartRepository := repository.NewCartRepository(1000)
	productServiceApi := productService.NewProductServiceApi(client, productAddress)
	lomsServiceApi := lomsService.NewLomsServiceApi(client, lomsAddress)

	s.addCartItemHandler = addItem.NewHandler(cartRepository, productServiceApi, lomsServiceApi, productToken)
	s.deleteCartItemHandler = deleteItem.NewHandler(cartRepository)
}

func (s *DeleteCartItemSuite) TestDeleteCartItem() {
	skuId := int64(773297411)
	userId := int64(123)
	count := uint16(1)

	cartParameters := model.CartParameters{
		SKU:    skuId,
		UserId: userId,
		Count:  count,
	}
	cartItem, err := s.addCartItemHandler.AddItem(cartParameters)

	cartParams := model.DeleteCartParameters{
		SKU:    skuId,
		UserId: userId,
	}
	actualResponse, err := s.deleteCartItemHandler.DeleteItem(cartParams)

	expectedResponse := &model.CartItem{
		SKU:    cartItem.SKU,
		Name:   cartItem.Name,
		Count:  cartItem.Count,
		Price:  cartItem.Price,
		UserId: cartItem.UserId,
	}
	require.Equal(s.T(), expectedResponse, actualResponse)
	require.NoError(s.T(), err)
}
