package service

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"homework/cart/internal/client/api/product/service"
	"homework/cart/internal/client/api/product/types"
	httpclient "homework/cart/internal/client/base/client"
	"homework/cart/internal/pkg/cart/model"
	"net/http"
	"time"
)

var productAddress = "http://route256.pavl.uk:8080"

var productToken = "testtoken"

type CartRepository interface {
	AddItem(params model.CartItem) (*model.CartItem, error, int)
	DeleteItem(params model.DeleteCartParameters) (*model.CartItem, error)
	DeleteItemsByUser(userId model.UserId) (*model.UserId, error)
	GetItemsByUser(userId model.UserId) ([]model.CartItem, error, int)
}

type CartService struct {
	repository CartRepository
}

func NewCartService(repository CartRepository) *CartService {
	return &CartService{repository: repository}
}

func (cartService CartService) AddItem(cartParams model.CartParameters) (*model.CartItem, error, int) {
	validate := validator.New()
	err := validate.Struct(cartParams)
	if err != nil {
		message := fmt.Sprintf("Invalid cart parameters")
		return nil, errors.New(message), http.StatusBadRequest
	}

	client := httpclient.NewHttpClient(10*time.Second, 3, []int{420, 429})
	productService := service.NewProductService(client, productAddress)

	request := types.ProductRequest{
		Sku:   cartParams.SKU,
		Token: productToken,
	}
	product, err, status := productService.GetProduct(request)
	if err != nil {
		return nil, err, status
	}

	cartItem := model.CartItem{
		SKU:    cartParams.SKU,
		Count:  cartParams.Count,
		UserId: cartParams.UserId,
		Name:   product.Name,
		Price:  product.Price,
	}
	return cartService.repository.AddItem(cartItem)
}

func (cartService CartService) DeleteItem(cartParams model.DeleteCartParameters) (*model.CartItem, error) {
	fmt.Println(cartParams)
	validate := validator.New()
	err := validate.Struct(cartParams)
	if err != nil {
		message := fmt.Sprintf("Invalid cart parameters")
		return nil, errors.New(message)
	}

	return cartService.repository.DeleteItem(cartParams)
}

func (cartService CartService) DeleteItemsByUser(userId model.UserId) (*model.UserId, error) {
	if userId <= 0 {
		message := fmt.Sprintf("Invalid parameters")
		return nil, errors.New(message)
	}

	return cartService.repository.DeleteItemsByUser(userId)
}

func (cartService CartService) GetCartByUser(userId model.UserId) (*model.Cart, error, int) {
	if userId <= 0 {
		message := fmt.Sprintf("Invalid parameters")
		return nil, errors.New(message), http.StatusInternalServerError
	}

	items, err, status := cartService.repository.GetItemsByUser(userId)
	if err != nil {
		return nil, err, status
	}

	if len(items) == 0 {
		message := fmt.Sprintf("Cart is empty")
		return nil, errors.New(message), http.StatusNotFound
	}

	client := httpclient.NewHttpClient(10*time.Second, 3, []int{420, 429})
	productService := service.NewProductService(client, productAddress)

	responseItems := make([]model.CartItem, 0, len(items))
	totalPrice := uint32(0)
	for _, item := range items {
		request := types.ProductRequest{
			Sku:   item.SKU,
			Token: productToken,
		}
		product, err, status := productService.GetProduct(request)
		if err != nil {
			return nil, err, status
		}

		cartItem := model.CartItem{
			SKU:   item.SKU,
			Count: item.Count,
			Name:  product.Name,
			Price: product.Price,
		}
		totalPrice += product.Price * uint32(item.Count)
		responseItems = append(responseItems, cartItem)
	}

	cart := &model.Cart{
		Items:      responseItems,
		TotalPrice: totalPrice,
	}
	return cart, nil, http.StatusOK
}
