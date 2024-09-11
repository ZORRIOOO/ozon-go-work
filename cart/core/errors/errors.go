package errors

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"strings"
)

func NewCustomError(message string, code int, writer http.ResponseWriter) {
	writer.WriteHeader(code)
	_, errOut := fmt.Fprintf(writer, "Error: %s", message)
	if errOut != nil {
		log.Printf(errOut.Error())
	}
}

func GetStatusCode(errMsg string) int {
	if strings.Contains(errMsg, "404") {
		return http.StatusNotFound
	} else if strings.Contains(errMsg, "429") {
		return http.StatusTooManyRequests
	} else {
		return http.StatusInternalServerError

	}
}

func GetValidationErrMsg(err error) string {
	return fmt.Sprintf("Validation: Field=%s", err.(validator.ValidationErrors)[0].Field())
}
