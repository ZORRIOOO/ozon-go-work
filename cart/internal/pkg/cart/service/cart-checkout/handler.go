package cart_checkout

import (
	"homework/cart/internal/client/api/loms/types"
	"homework/cart/internal/pkg/cart/model"
)

type (
	LomsService interface {
		CreateOrder(request types.OrderCreateRequest) (*types.OrderCreateResponse, error)
	}

	CartRepository interface {
		AddItem(params model.CartItem) (*model.CartItem, error)
		DeleteItem(params model.DeleteCartParameters) (*model.CartItem, error)
		DeleteItemsByUser(userId model.UserId) (*model.UserId, error)
		GetItemsByUser(userId model.UserId) ([]model.CartItem, error)
	}

	CartServiceHandler struct {
		repository CartRepository
		lomsApi    LomsService
	}
)

func (c CartServiceHandler) CartCheckout(userId model.UserId) (*model.Checkout, error) {
	itemsByUser, err := c.repository.GetItemsByUser(userId)
	if err != nil {
		return nil, err
	}

	items := make([]types.Item, 0)
	for _, item := range itemsByUser {
		items = append(items, types.Item{
			Sku:   item.SKU,
			Count: item.Count,
		})
	}
	request := types.OrderCreateRequest{
		User:  userId,
		Items: items,
	}

	orderResponse, err := c.lomsApi.CreateOrder(request)
	if err != nil {
		return nil, err
	}

	_, deleteErr := c.repository.DeleteItemsByUser(userId)
	if deleteErr != nil {
		return nil, deleteErr
	}

	return &model.Checkout{OrderId: orderResponse.OrderId}, nil
}

func NewHandler(repository CartRepository, lomsApi LomsService) *CartServiceHandler {
	return &CartServiceHandler{
		repository: repository,
		lomsApi:    lomsApi,
	}
}
