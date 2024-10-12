package main

import (
	"errors"
	"fmt"
	"homework/cart/internal/app/server"
	lomsService "homework/cart/internal/client/api/loms/service"
	productService "homework/cart/internal/client/api/product/service"
	httpclient "homework/cart/internal/client/base/client"
	"homework/cart/internal/http/middleware"
	"homework/cart/internal/pkg/cart/channel"
	"homework/cart/internal/pkg/cart/repository"
	addItem "homework/cart/internal/pkg/cart/service/add-item"
	cartCheckout "homework/cart/internal/pkg/cart/service/cart-checkout"
	deleteCart "homework/cart/internal/pkg/cart/service/delete-cart"
	deleteItem "homework/cart/internal/pkg/cart/service/delete-item"
	getCart "homework/cart/internal/pkg/cart/service/get-cart"
	"log"
	"net/http"
	"time"
)

const (
	httpPort       = ":8082"
	productAddress = "http://route256.pavl.uk:8080"
	lomsAddress    = "http://loms:8081"
	productToken   = "testtoken"
	capacity       = 1000
	clientTimeout  = 10 * time.Second
	clientRetries  = 3
	rpc            = 10
	maxRate        = 10
)

func HealthCheckHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

func main() {
	log.Println("Go cart emitter starting")

	client := httpclient.NewHttpClient(clientTimeout, clientRetries, []int{420, 429})
	productServiceApi := productService.NewProductServiceApi(client, productAddress)
	lomsServiceApi := lomsService.NewLomsServiceApi(client, lomsAddress)
	cartRepository := repository.NewCartRepository(capacity)
	cartChannel := channel.NewCartChannel(productServiceApi, productToken, rpc, maxRate)
	addItemHandler := addItem.NewHandler(cartRepository, productServiceApi, lomsServiceApi, productToken)
	deleteItemHandler := deleteItem.NewHandler(cartRepository)
	deleteCartHandler := deleteCart.NewHandler(cartRepository)
	getCartHandler := getCart.NewHandler(cartRepository, cartChannel)
	cartCheckoutHandler := cartCheckout.NewHandler(cartRepository, lomsServiceApi)
	appServer := server.NewServer(
		addItemHandler,
		deleteItemHandler,
		deleteCartHandler,
		getCartHandler,
		cartCheckoutHandler,
	)

	mux := http.NewServeMux()
	mux.Handle("GET /healthcheck", http.HandlerFunc(HealthCheckHandler))
	mux.HandleFunc("POST /user/{user_id}/cart/{sku_id}", appServer.AddCartItem)
	mux.HandleFunc("DELETE /user/{user_id}/cart/{sku_id}", appServer.DeleteCartItem)
	mux.HandleFunc("DELETE /user/{user_id}/cart", appServer.DeleteCart)
	mux.HandleFunc("GET /user/{user_id}/cart", appServer.GetCartByUser)
	mux.HandleFunc("PUT /user/{user_id}/checkout", appServer.CartCheckout)

	loggingMux := middleware.NewLoggingMux(mux)
	log.Println("Go cart emitter ready")
	if err := http.ListenAndServe(httpPort, loggingMux); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Error closed server: %s", err.Error())
		}
		if err != nil {
			log.Fatalf("Error starting server: %s", err.Error())
		}
	}
}
