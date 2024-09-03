package service

import (
	httpclient "cart/internal/clients/base"
	"cart/internal/clients/product/service"
	"cart/internal/clients/product/types"
	"cart/internal/pkg/cart/model"
	"time"
)

var productAddress = "http://route256.pavl.uk:8080"

var productToken = "testtoken"

type CartRepository interface {
	AddItem(params model.CartItem) (*model.CartItem, error)
	DeleteItem(item model.CartItem) (*model.CartItem, error)
}

type CartService struct {
	repository CartRepository
}

func NewCartService(repository CartRepository) *CartService {
	return &CartService{repository: repository}
}

func (cartService CartService) AddItem(cartParams model.CartParameters) (*model.CartItem, error) {
	client := httpclient.NewHttpClient(10 * time.Second)
	productService := service.NewProductService(client, productAddress)

	request := types.ProductRequest{
		Sku:   cartParams.SKU,
		Token: productToken,
	}
	response, err := productService.GetProduct(request)
	if err != nil {
		return nil, err
	}

	cartItem := model.CartItem{
		SKU:    cartParams.SKU,
		Count:  cartParams.Count,
		UserId: cartParams.UserId,
		Name:   response.Name,
		Price:  response.Price,
	}
	return cartService.repository.AddItem(cartItem)
}

func (cartService CartService) DeleteItem(cartItem model.CartItem) (*model.CartItem, error) {
	return nil, nil
}
