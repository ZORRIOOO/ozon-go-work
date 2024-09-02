package server

import (
	"cart/internal/pkg/cart/model"
	"context"
)

type CartService interface {
	AddItem(ctx context.Context, item model.CartItem) (*model.CartItem, error)
	DeleteItem(ctx context.Context, item model.CartItem) (*model.CartItem, error)
}

type Server struct {
	cartService CartService
}

func NewServer(cartService CartService) *Server {
	return &Server{cartService: cartService}
}
