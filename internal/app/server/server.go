package server

import (
	"cart/internal/pkg/cart/model"
)

type CartService interface {
	AddItem(params model.CartParameters) (*model.CartItem, error)
	DeleteItem(item model.CartItem) (*model.CartItem, error)
}

type Server struct {
	cartService CartService
}

func NewServer(cartService CartService) *Server {
	return &Server{cartService: cartService}
}
