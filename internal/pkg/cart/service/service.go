package service

import (
	"cart/internal/client/api/product/service"
	"cart/internal/client/api/product/types"
	httpclient "cart/internal/client/base"
	"cart/internal/pkg/cart/model"
	"errors"
	"fmt"
	"time"
)

var productAddress = "http://route256.pavl.uk:8080"

var productToken = "testtoken"

type CartRepository interface {
	AddItem(params model.CartItem) (*model.CartItem, error)
	DeleteItem(skuId model.SKU, userId model.UserId) (*model.CartItem, error)
	DeleteItemsByUser(userId model.UserId) (*model.UserId, error)
}

type CartService struct {
	repository CartRepository
}

func NewCartService(repository CartRepository) *CartService {
	return &CartService{repository: repository}
}

func (cartService CartService) AddItem(cartParams model.CartParameters) (*model.CartItem, error) {
	if cartParams.Count <= 0 || cartParams.SKU <= 0 || cartParams.UserId <= 0 {
		message := fmt.Sprintf("Invalid cart parameters")
		return nil, errors.New(message)
	}

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

func (cartService CartService) DeleteItem(cartParams model.CartParameters) (*model.CartItem, error) {
	if cartParams.SKU <= 0 || cartParams.UserId <= 0 {
		message := fmt.Sprintf("Invalid cart parameters")
		return nil, errors.New(message)
	}

	return cartService.repository.DeleteItem(cartParams.SKU, cartParams.UserId)
}

func (cartService CartService) DeleteItemsByUser(userId model.UserId) (*model.UserId, error) {
	if userId <= 0 {
		message := fmt.Sprintf("Invalid parameters")
		return nil, errors.New(message)
	}

	return cartService.repository.DeleteItemsByUser(userId)
}

func (cartService CartService) GetCartByUser(userId model.UserId) (any, error) {
	return nil, nil
}
