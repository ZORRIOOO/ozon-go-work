package delete_cart

import (
	"context"
	"errors"
	"fmt"
	"homework/cart/internal/client/api/product/types"
	"homework/cart/internal/pkg/cart/model"
)

type (
	CartRepository interface {
		AddItem(params model.CartItem) (*model.CartItem, error)
		DeleteItem(params model.DeleteCartParameters) (*model.CartItem, error)
		DeleteItemsByUser(userId model.UserId) (*model.UserId, error)
		GetItemsByUser(userId model.UserId) ([]model.CartItem, error)
	}

	ProductService interface {
		GetProduct(ctx context.Context, request types.ProductRequest) (*types.ProductResponse, error)
		GetSkuList(ctx context.Context, request types.SkusRequest) (*types.SkusResponse, error)
	}

	CartServiceHandler struct {
		repository   CartRepository
		productApi   ProductService
		productToken string
	}
)

func NewHandler(repository CartRepository) *CartServiceHandler {
	return &CartServiceHandler{
		repository: repository,
	}
}

func (cartService CartServiceHandler) DeleteItemsByUser(userId model.UserId) (*model.UserId, error) {
	if userId <= 0 {
		message := fmt.Sprintf("Invalid parameters: userId=%d", userId)
		return nil, errors.New(message)
	}

	return cartService.repository.DeleteItemsByUser(userId)
}
