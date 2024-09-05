package main

import (
	"fmt"
	"homework/cart/internal/app/server"
	"homework/cart/internal/http/middleware"
	"homework/cart/internal/pkg/cart/repository"
	"homework/cart/internal/pkg/cart/service"
	"log"
	"net/http"
)

var addr = "localhost:8082"

func HealthCheckHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

func main() {
	log.Println("Go app starting")

	cartRepository := repository.NewCartRepository(100)
	cartService := service.NewCartService(cartRepository)
	appServer := server.NewServer(cartService)
	log.Println("Server starting")

	mux := http.NewServeMux()
	mux.Handle("GET /healthcheck", http.HandlerFunc(HealthCheckHandler))
	mux.HandleFunc("POST /user/{user_id}/cart/{sku_id}", appServer.AddCartItem)
	mux.HandleFunc("DELETE /user/{user_id}/cart/{sku_id}", appServer.DeleteCartItem)
	mux.HandleFunc("DELETE /user/{user_id}/cart", appServer.DeleteCart)
	mux.HandleFunc("GET /user/{user_id}/cart", appServer.GetCartByUser)

	loggingMux := middleware.NewLoggingMux(mux)

	if err := http.ListenAndServe(addr, loggingMux); err != nil {
		panic(err)
	}
}
