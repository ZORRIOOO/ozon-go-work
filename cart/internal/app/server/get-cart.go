package server

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"homework/cart/core/errors"
	"net/http"
	"strconv"
)

type GetCartRequest struct {
	UserId int64 `json:"user_id" validate:"required,min=1"`
}

func (s *Server) GetCartByUser(w http.ResponseWriter, r *http.Request) {
	rawUserId := r.PathValue("user_id")
	userId, _ := strconv.ParseInt(rawUserId, 10, 64)

	request := GetCartRequest{
		UserId: userId,
	}
	validate := validator.New()
	err := validate.Struct(request)
	if err != nil {
		message := fmt.Sprintf("GET /user/{user_id}/cart: %s", errors.GetValidationErrMsg(err))
		errors.NewCustomError(message, http.StatusBadRequest, w)
		return
	}

	response, err := s.cartService.GetCartByUser(userId)
	if err != nil {
		message := fmt.Sprintf("GET /user/{user_id}/cart: %s", err.Error())
		errors.NewCustomError(message, http.StatusInternalServerError, w)
		return
	} else {
		rawResponse, marshalErr := json.Marshal(response)
		if marshalErr != nil {
			message := fmt.Sprintf("GET /user/{user_id}/cart: %s", marshalErr.Error())
			errors.NewCustomError(message, http.StatusInternalServerError, w)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(rawResponse)
	}
}
