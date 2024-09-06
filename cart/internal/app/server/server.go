package server

import (
	"homework/cart/internal/pkg/cart/model"
)

type CartService interface {
	AddItem(params model.CartParameters) (*model.CartItem, error, int)
	DeleteItem(item model.DeleteCartParameters) (*model.CartItem, error)
	DeleteItemsByUser(userId model.UserId) (*model.UserId, error)
	GetCartByUser(userId model.UserId) (*model.Cart, error, int)
}

type Server struct {
	cartService CartService
}

func NewServer(cartService CartService) *Server {
	return &Server{cartService: cartService}
}
