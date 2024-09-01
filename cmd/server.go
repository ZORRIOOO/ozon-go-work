package main

import (
	"cart/internal/http/middleware"
	"log"
	"net/http"
)

var addr = "127.0.0.1:8082"

func HealthCheckHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func main() {
	log.Println("server starting")

	mux := http.NewServeMux()

	mux.Handle("GET /healthcheck", http.HandlerFunc(HealthCheckHandler))

	loggingMux := middleware.NewLoggingMux(mux)

	if err := http.ListenAndServe(addr, loggingMux); err != nil {
		panic(err)
	}
}
