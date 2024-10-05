package suite

import (
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	lomsService "homework/cart/internal/client/api/loms/service"
	productService "homework/cart/internal/client/api/product/service"
	httpclient "homework/cart/internal/client/base/client"
	"homework/cart/internal/pkg/cart/channel"
	"homework/cart/internal/pkg/cart/model"
	"homework/cart/internal/pkg/cart/repository"
	addItem "homework/cart/internal/pkg/cart/service/add-item"
	getCart "homework/cart/internal/pkg/cart/service/get-cart"
	"time"
)

type GetCartSuite struct {
	suite.Suite
	addCartItemHandler *addItem.CartServiceHandler
	getCartHandler     *getCart.CartServiceHandler
}

func (s *GetCartSuite) SetupSuite() {
	const (
		productAddress = "http://route256.pavl.uk:8080"
		productToken   = "testtoken"
		lomsAddress    = "http://localhost:8081"
		rpc            = 10
		maxRate        = 10
	)
	client := httpclient.NewHttpClient(10*time.Second, 3, []int{420, 429})
	cartRepository := repository.NewCartRepository(1000)
	productServiceApi := productService.NewProductServiceApi(client, productAddress)
	lomsServiceApi := lomsService.NewLomsServiceApi(client, lomsAddress)
	cartChannel := channel.NewCartChannel(productServiceApi, productToken, rpc, maxRate)

	s.addCartItemHandler = addItem.NewHandler(cartRepository, productServiceApi, lomsServiceApi, productToken)
	s.getCartHandler = getCart.NewHandler(cartRepository, cartChannel)
}

func (s *GetCartSuite) TestGetCart() {
	skuId := int64(773297411)
	userId := int64(123)
	count := uint16(1)

	cartParameters := model.CartParameters{
		SKU:    skuId,
		UserId: userId,
		Count:  count,
	}
	cartItem, err := s.addCartItemHandler.AddItem(cartParameters)

	require.NoError(s.T(), err)

	actualResponse, err := s.getCartHandler.GetCartByUser(cartItem.UserId)

	require.NoError(s.T(), err)

	expectedItems := make([]model.CartItem, 0)
	expectedItem := model.CartItem{
		SKU:    cartItem.SKU,
		Name:   cartItem.Name,
		Count:  cartItem.Count,
		Price:  cartItem.Price,
		UserId: cartItem.UserId,
	}

	expectedItems = append(expectedItems, expectedItem)
	expectedResponse := &model.Cart{
		Items:      expectedItems,
		TotalPrice: expectedItem.Price,
	}

	require.Equal(s.T(), expectedResponse, actualResponse)
	require.NoError(s.T(), err)
}
