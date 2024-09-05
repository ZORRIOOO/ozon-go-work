package errors

import (
	"fmt"
	"log"
	"net/http"
)

func NewCustomError(message string, code int, writer http.ResponseWriter) {
	writer.WriteHeader(code)
	_, errOut := fmt.Fprintf(writer, "Error: %s", message)
	if errOut != nil {
		log.Printf(errOut.Error())
	}
}
