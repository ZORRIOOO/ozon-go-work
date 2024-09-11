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
	log.Printf("Incoming request: Method=%s, URL=%s, Headers=%v", r.Method, r.URL.String(), r.Header)
	m.h.ServeHTTP(w, r)
	log.Printf("Request processed: Method=%s, URL=%s", r.Method, r.URL.String())
}
