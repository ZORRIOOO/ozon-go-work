package middleware

import (
	"log"
	"net/http"
)

type LoggingMux struct {
	h http.Handler
}

func NewLoggingMux(h http.Handler) http.Handler {
	return &LoggingMux{h: h}
}

func (m *LoggingMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("request gotten")
	m.h.ServeHTTP(w, r)
	log.Printf("request processed")
}
