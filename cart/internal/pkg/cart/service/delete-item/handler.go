package delete_item

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	customErrors "homework/cart/core/errors"
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
		GetProduct(context context.Context, request types.ProductRequest) (*types.ProductResponse, error)
		GetSkuList(context context.Context, request types.SkusRequest) (*types.SkusResponse, error)
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

func (cartService CartServiceHandler) DeleteItem(cartParams model.DeleteCartParameters) (*model.CartItem, error) {
	validate := validator.New()
	err := validate.Struct(cartParams)
	if err != nil {
		return nil, errors.New(customErrors.GetValidationErrMsg(err))
	}

	return cartService.repository.DeleteItem(cartParams)
}
