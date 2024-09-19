package suite

import (
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	productService "homework/cart/internal/client/api/product/service"
	httpclient "homework/cart/internal/client/base/client"
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
	var productAddress = "http://route256.pavl.uk:8080"
	var productToken = "testtoken"
	client := httpclient.NewHttpClient(10*time.Second, 3, []int{420, 429})
	productServiceApi := productService.NewProductServiceApi(client, productAddress)
	cartRepository := repository.NewCartRepository(100)

	s.addCartItemHandler = addItem.NewHandler(cartRepository, productServiceApi, productToken)
	s.getCartHandler = getCart.NewHandler(cartRepository, productServiceApi, productToken)
}

func (s *GetCartSuite) TestGetCart() {
	skuId := int64(1148162)
	userId := int64(123)
	count := uint16(1)

	cartParameters := model.CartParameters{
		SKU:    skuId,
		UserId: userId,
		Count:  count,
	}
	cartItem, err := s.addCartItemHandler.AddItem(cartParameters)

	actualResponse, err := s.getCartHandler.GetCartByUser(cartItem.UserId)

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
