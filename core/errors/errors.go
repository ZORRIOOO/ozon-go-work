package errors

import (
	"fmt"
	"log"
	"net/http"
)

type ErrorContext struct {
	error   error
	writer  http.ResponseWriter
	message string
}

func NewCustomError(writer http.ResponseWriter, message string, code int, extra any) {
	writer.WriteHeader(code)
	_, errOut := fmt.Fprintf(writer, "{\"message\":\"%s\"}", extra)
	if errOut != nil {
		log.Printf("Error: %s - %s", message, errOut.Error())
	}
}
