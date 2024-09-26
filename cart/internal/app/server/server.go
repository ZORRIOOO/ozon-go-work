package server

import (
	"homework/cart/internal/pkg/cart/model"
)

type AddItemHandler interface {
	AddItem(params model.CartParameters) (*model.CartItem, error)
}

type DeleteItemHandler interface {
	DeleteItem(item model.DeleteCartParameters) (*model.CartItem, error)
}

type DeleteCartHandler interface {
	DeleteItemsByUser(userId model.UserId) (*model.UserId, error)
}

type GetCartHandler interface {
	GetCartByUser(userId model.UserId) (*model.Cart, error)
}

type CartCheckoutHandler interface {
	CartCheckout(userId model.UserId) (*model.Checkout, error)
}

type Server struct {
	addItemHandler      AddItemHandler
	deleteItemHandler   DeleteItemHandler
	deleteCartHandler   DeleteCartHandler
	getCartHandler      GetCartHandler
	cartCheckoutHandler CartCheckoutHandler
}

func NewServer(
	addItemHandler AddItemHandler,
	deleteItemHandler DeleteItemHandler,
	deleteCartHandler DeleteCartHandler,
	getCartHandler GetCartHandler,
	cartCheckoutHandler CartCheckoutHandler,
) *Server {
	return &Server{
		addItemHandler:      addItemHandler,
		deleteItemHandler:   deleteItemHandler,
		deleteCartHandler:   deleteCartHandler,
		getCartHandler:      getCartHandler,
		cartCheckoutHandler: cartCheckoutHandler,
	}
}
