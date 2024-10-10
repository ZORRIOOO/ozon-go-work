package add_item

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	customErrors "homework/cart/core/errors"
	lomsTypes "homework/cart/internal/client/api/loms/types"
	productTypes "homework/cart/internal/client/api/product/types"
	"homework/cart/internal/pkg/cart/model"
	"net/http"
)

type (
	CartRepository interface {
		AddItem(params model.CartItem) (*model.CartItem, error)
		DeleteItem(params model.DeleteCartParameters) (*model.CartItem, error)
		DeleteItemsByUser(userId model.UserId) (*model.UserId, error)
		GetItemsByUser(userId model.UserId) ([]model.CartItem, error)
	}

	ProductService interface {
		GetProduct(context context.Context, request productTypes.ProductRequest) (*productTypes.ProductResponse, error)
		GetSkuList(context context.Context, request productTypes.SkusRequest) (*productTypes.SkusResponse, error)
	}

	LomsService interface {
		CreateOrder(request lomsTypes.OrderCreateRequest) (*lomsTypes.OrderCreateResponse, error)
		StocksInfo(request lomsTypes.StocksInfoRequest) (*lomsTypes.StocksInfoResponse, error)
	}

	CartServiceHandler struct {
		repository   CartRepository
		productApi   ProductService
		lomsApi      LomsService
		productToken string
	}
)

func NewHandler(repository CartRepository, productApi ProductService, lomsApi LomsService, productToken string) *CartServiceHandler {
	return &CartServiceHandler{
		repository:   repository,
		productApi:   productApi,
		lomsApi:      lomsApi,
		productToken: productToken,
	}
}

func (cartService CartServiceHandler) AddItem(cartParams model.CartParameters) (*model.CartItem, error) {
	ctx := context.Background()
	validate := validator.New()
	err := validate.Struct(cartParams)
	if err != nil {
		return nil, errors.New(customErrors.GetValidationErrMsg(err))
	}

	product, err := cartService.productApi.GetProduct(
		ctx,
		productTypes.ProductRequest{
			Sku:   cartParams.SKU,
			Token: cartService.productToken,
		},
	)
	if err != nil {
		return nil, err
	}

	stocks, err := cartService.lomsApi.StocksInfo(
		lomsTypes.StocksInfoRequest{
			Sku: cartParams.SKU,
		},
	)
	if err != nil {
		return nil, err
	}

	totalAvailable := stocks.Count
	if totalAvailable < cartParams.Count {
		message := fmt.Sprintf(
			"%d Failed Precondition Insufficient stock for SKU: %d, Available: %d, Requested: %d",
			http.StatusPreconditionFailed,
			cartParams.SKU,
			totalAvailable,
			cartParams.Count,
		)
		return nil, errors.New(message)
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
