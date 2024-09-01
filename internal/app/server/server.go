package server

type CartService interface{}

type Server struct {
	cartService CartService
}

func NewServer(cartService CartService) *Server {
	return &Server{cartService: cartService}
}
