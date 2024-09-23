package server

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"homework/cart/core/errors"
	"net/http"
	"strconv"
)

type CartCheckoutRequest struct {
	UserId int64 `json:"user_id" validate:"required,min=1"`
}

func (s *Server) CartCheckout(w http.ResponseWriter, r *http.Request) {
	rawUserId := r.PathValue("user_id")
	userId, _ := strconv.ParseInt(rawUserId, 10, 64)

	request := CartCheckoutRequest{
		UserId: userId,
	}
	validate := validator.New()
	err := validate.Struct(request)
	if err != nil {
		message := fmt.Sprintf("POST /user/{user_id}/checkout: %s", errors.GetValidationErrMsg(err))
		errors.NewCustomError(message, http.StatusBadRequest, w)
		return
	}

	response, err := s.cartCheckoutHandler.CartCheckout(userId)
	if err != nil {
		message := fmt.Sprintf("POST /user/{user_id}/checkout: %s", err.Error())
		errors.NewCustomError(message, http.StatusInternalServerError, w)
		return
	} else {
		rawResponse, marshalErr := json.Marshal(response)
		if marshalErr != nil {
			message := fmt.Sprintf("POST /user/{user_id}/checkout: %s", marshalErr.Error())
			errors.NewCustomError(message, http.StatusInternalServerError, w)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(rawResponse)
	}

}
