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
	cartCheckout "homework/cart/internal/pkg/cart/service/cart-checkout"
	"time"
)

type CartCheckoutSuite struct {
	suite.Suite
	addCartItemHandler  *addItem.CartServiceHandler
	cartCheckoutHandler *cartCheckout.CartServiceHandler
}

func (s *CartCheckoutSuite) SetupSuite() {
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
	s.cartCheckoutHandler = cartCheckout.NewHandler(cartRepository, lomsServiceApi)
}

func (s *CartCheckoutSuite) TestCartCheckout() {
	skuId := int64(773297411)
	userId := int64(1)
	count := uint16(1)

	cartParameters := model.CartParameters{
		SKU:    skuId,
		UserId: userId,
		Count:  count,
	}
	cartItem, err := s.addCartItemHandler.AddItem(cartParameters)

	require.NoError(s.T(), err)

	actualResponse, err := s.cartCheckoutHandler.CartCheckout(cartItem.UserId)
	expectedResponse := &model.Checkout{OrderId: actualResponse.OrderId}

	require.NoError(s.T(), err)
	require.Equal(s.T(), expectedResponse, actualResponse)
}
