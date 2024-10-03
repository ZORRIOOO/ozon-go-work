package get_cart

import (
	"errors"
	"fmt"
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

	CartChannel interface {
		FetchProductsInParallel(items []model.CartItem, userId model.UserId) ([]model.CartItem, uint32, error)
	}

	CartServiceHandler struct {
		repository CartRepository
		channel    CartChannel
	}
)

func NewHandler(repository CartRepository, cartChannel CartChannel) *CartServiceHandler {
	return &CartServiceHandler{
		repository: repository,
		channel:    cartChannel,
	}
}

func (handler CartServiceHandler) GetCartByUser(userId model.UserId) (*model.Cart, error) {
	if userId <= 0 {
		message := fmt.Sprintf("Invalid parameters: userId=%d", userId)
		return nil, errors.New(message)
	}

	items, err := handler.repository.GetItemsByUser(userId)
	if err != nil {
		return nil, err
	}

	if items == nil || len(items) == 0 {
		message := fmt.Sprintf("Cart is empty")
		return nil, errors.New(message)
	}

	responseItems, totalPrice, err := handler.channel.FetchProductsInParallel(items, userId)
	if err != nil {
		return nil, err
	}

	sort.Slice(responseItems, func(i, j int) bool {
		return responseItems[i].SKU < responseItems[j].SKU
	})

	return &model.Cart{
		Items:      responseItems,
		TotalPrice: totalPrice,
	}, nil
}
