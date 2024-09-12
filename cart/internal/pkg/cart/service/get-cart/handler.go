package get_cart

import (
	"errors"
	"fmt"
	productServiceApi "homework/cart/internal/client/api/product/service"
	"homework/cart/internal/client/api/product/types"
	"homework/cart/internal/pkg/cart/model"
	"sort"
)

type (
	CartRepository interface {
		AddItem(params model.CartItem) (*model.CartItem, error)
		DeleteItem(params model.DeleteCartParameters) (*model.CartItem, error)
		DeleteItemsByUser(userId model.UserId) (*model.UserId, error)
		GetItemsByUser(userId model.UserId) ([]model.CartItem, error)
	}

	ProductService interface {
		GetProduct(request types.ProductRequest) (*types.ProductResponse, error)
		GetSkuList(request types.SkusRequest) (*types.SkusResponse, error)
	}

	CartServiceHandler struct {
		repository   CartRepository
		productApi   ProductService
		productToken string
	}
)

func NewHandler(repository CartRepository, productApi productServiceApi.ProductService, productToken string) *CartServiceHandler {
	return &CartServiceHandler{
		repository:   repository,
		productApi:   productApi,
		productToken: productToken,
	}
}

func (cartService CartServiceHandler) GetCartByUser(userId model.UserId) (*model.Cart, error) {
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
