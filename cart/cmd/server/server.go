package main

import (
	"errors"
	"fmt"
	"homework/cart/internal/app/server"
	productService "homework/cart/internal/client/api/product/service"
	httpclient "homework/cart/internal/client/base/client"
	"homework/cart/internal/http/middleware"
	"homework/cart/internal/pkg/cart/repository"
	cartServiceInternal "homework/cart/internal/pkg/cart/service"
	"log"
	"net/http"
	"time"
)

var addr = ":8082"

var productAddress = "http://route256.pavl.uk:8080"

func HealthCheckHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

func main() {
	log.Println("Go app starting")

	cartRepository := repository.NewCartRepository(100)
	client := httpclient.NewHttpClient(10*time.Second, 3, []int{420, 429})
	productServiceApi := productService.NewProductServiceApi(client, productAddress)
	cartService := cartServiceInternal.NewCartService(cartRepository, productServiceApi)
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
		if errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Error closed server: %s", err.Error())
		}
		if err != nil {
			log.Fatalf("Error starting server: %s", err.Error())
		}
	}
}
