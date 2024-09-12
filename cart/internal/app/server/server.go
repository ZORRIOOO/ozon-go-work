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

type Server struct {
	addItemHandler    AddItemHandler
	deleteItemHandler DeleteItemHandler
	deleteCartHandler DeleteCartHandler
	getCartHandler    GetCartHandler
}

func NewServer(
	addItemHandler AddItemHandler,
	deleteItemHandler DeleteItemHandler,
	deleteCartHandler DeleteCartHandler,
	getCartHandler GetCartHandler,
) *Server {
	return &Server{
		addItemHandler:    addItemHandler,
		deleteItemHandler: deleteItemHandler,
		deleteCartHandler: deleteCartHandler,
		getCartHandler:    getCartHandler,
	}
}
