package service

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	customErrors "homework/cart/core/errors"
	productServiceApi "homework/cart/internal/client/api/product/service"
	"homework/cart/internal/client/api/product/types"
	"homework/cart/internal/pkg/cart/model"
	"sort"
)

type CartRepository interface {
	AddItem(params model.CartItem) (*model.CartItem, error)
	DeleteItem(params model.DeleteCartParameters) (*model.CartItem, error)
	DeleteItemsByUser(userId model.UserId) (*model.UserId, error)
	GetItemsByUser(userId model.UserId) ([]model.CartItem, error)
}

type CartService struct {
	repository   CartRepository
	productApi   productServiceApi.ProductService
	productToken string
}

func NewCartService(repository CartRepository, productApi productServiceApi.ProductService, productToken string) *CartService {
	return &CartService{
		repository:   repository,
		productApi:   productApi,
		productToken: productToken,
	}
}

func (cartService CartService) AddItem(cartParams model.CartParameters) (*model.CartItem, error) {
	validate := validator.New()
	err := validate.Struct(cartParams)
	if err != nil {
		return nil, errors.New(customErrors.GetValidationErrMsg(err))
	}

	request := types.ProductRequest{
		Sku:   cartParams.SKU,
		Token: cartService.productToken,
	}
	product, err := cartService.productApi.GetProduct(request)
	if err != nil {
		return nil, err
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
	validate := validator.New()
	err := validate.Struct(cartParams)
	if err != nil {
		return nil, errors.New(customErrors.GetValidationErrMsg(err))
	}

	return cartService.repository.DeleteItem(cartParams)
}

func (cartService CartService) DeleteItemsByUser(userId model.UserId) (*model.UserId, error) {
	if userId <= 0 {
		message := fmt.Sprintf("Invalid parameters: userId=%d", userId)
		return nil, errors.New(message)
	}

	return cartService.repository.DeleteItemsByUser(userId)
}

func (cartService CartService) GetCartByUser(userId model.UserId) (*model.Cart, error) {
	if userId <= 0 {
		message := fmt.Sprintf("Invalid parameters: userId=%d", userId)
		return nil, errors.New(message)
	}

	items, err := cartService.repository.GetItemsByUser(userId)
	if err != nil {
		return nil, err
	}

	if items == nil || len(items) == 0 {
		message := fmt.Sprintf("Cart is empty")
		return nil, errors.New(message)
	}

	responseItems := make([]model.CartItem, 0, len(items))
	totalPrice := uint32(0)
	for _, item := range items {
		request := types.ProductRequest{
			Sku:   item.SKU,
			Token: cartService.productToken,
		}
		product, err := cartService.productApi.GetProduct(request)
		if err != nil {
			return nil, err
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

	sort.Slice(responseItems, func(i, j int) bool {
		return responseItems[i].SKU < responseItems[j].SKU
	})

	cart := &model.Cart{
		Items:      responseItems,
		TotalPrice: totalPrice,
	}
	return cart, nil
}
