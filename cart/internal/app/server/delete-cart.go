package server

import (
	"fmt"
	"homework/cart/core/errors"
	"net/http"
	"strconv"
)

func (s *Server) DeleteCart(w http.ResponseWriter, r *http.Request) {
	rawUserId := r.PathValue("user_id")
	userId, err := strconv.ParseInt(rawUserId, 10, 64)
	if userId <= 0 || err != nil {
		errors.NewCustomError("DELETE /user/{user_id}/cart: Invalid path params", http.StatusBadRequest, w)
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
