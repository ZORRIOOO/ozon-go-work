package server

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"homework/cart/core/errors"
	"net/http"
	"strconv"
)

type DeleteCartRequest struct {
	UserId int64 `json:"user_id" validate:"required,min=1"`
}

func (s *Server) DeleteCart(w http.ResponseWriter, r *http.Request) {
	rawUserId := r.PathValue("user_id")
	userId, _ := strconv.ParseInt(rawUserId, 10, 64)

	deleteRequest := DeleteCartRequest{
		UserId: userId,
	}
	validate := validator.New()
	err := validate.Struct(deleteRequest)
	if err != nil {
		message := fmt.Sprintf("DELETE /user/{user_id}/cart: %s", errors.GetValidationErrMsg(err))
		errors.NewCustomError(message, http.StatusBadRequest, w)
		return
	}

	response, err := s.cartService.DeleteItemsByUser(userId)
	if err != nil {
		message := fmt.Sprintf("DELETE /user/{user_id}/cart: %s", err.Error())
		errors.NewCustomError(message, http.StatusInternalServerError, w)
		return
	}

	if response != nil {
		fmt.Fprint(w, http.StatusNoContent)
	}
}
