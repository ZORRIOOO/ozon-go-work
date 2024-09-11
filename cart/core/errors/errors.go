package errors

import (
	"errors"
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
	var validationErr = err.(validator.ValidationErrors)
	if len(validationErr) > 0 && errors.Is(err, &validationErr) {
		return fmt.Sprintf("Validation: Field=%s", validationErr[0].Field())
	} else {
		return fmt.Sprint("Validation error")
	}
}
