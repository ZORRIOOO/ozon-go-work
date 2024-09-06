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
	userId, err := strconv.ParseInt(rawUserId, 10, 64)

	request := GetCartRequest{
		UserId: userId,
	}
	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		errors.NewCustomError("GET /user/{user_id}/cart: Invalid path params", http.StatusBadRequest, w)
		return
	}

	response, err, status := s.cartService.GetCartByUser(userId)
	if err != nil {
		message := fmt.Sprintf("GET /user/{user_id}/cart: %s", err.Error())
		errors.NewCustomError(message, status, w)
		return
	} else {
		rawResponse, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(rawResponse)
	}
}
