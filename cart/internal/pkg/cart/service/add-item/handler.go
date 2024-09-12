package add_item

import (
	"errors"
	"github.com/go-playground/validator/v10"
	customErrors "homework/cart/core/errors"
	productServiceApi "homework/cart/internal/client/api/product/service"
	"homework/cart/internal/client/api/product/types"
	"homework/cart/internal/pkg/cart/model"
)

type CartRepository interface {
	AddItem(params model.CartItem) (*model.CartItem, error)
	DeleteItem(params model.DeleteCartParameters) (*model.CartItem, error)
	DeleteItemsByUser(userId model.UserId) (*model.UserId, error)
	GetItemsByUser(userId model.UserId) ([]model.CartItem, error)
}

type ProductService interface {
	GetProduct(request types.ProductRequest) (*types.ProductResponse, error)
	GetSkuList(request types.SkusRequest) (*types.SkusResponse, error)
}

type CartServiceHandler struct {
	repository   CartRepository
	productApi   ProductService
	productToken string
}

func NewHandler(repository CartRepository, productApi productServiceApi.ProductService, productToken string) *CartServiceHandler {
	return &CartServiceHandler{
		repository:   repository,
		productApi:   productApi,
		productToken: productToken,
	}
}

func (cartService CartServiceHandler) AddItem(cartParams model.CartParameters) (*model.CartItem, error) {
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
