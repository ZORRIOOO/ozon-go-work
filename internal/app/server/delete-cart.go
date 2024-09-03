package server

import (
	"cart/core/errors"
	"fmt"
	"net/http"
	"strconv"
)

func (s *Server) DeleteCart(w http.ResponseWriter, r *http.Request) {
	rawUserId := r.PathValue("user_id")
	userId, err := strconv.ParseInt(rawUserId, 10, 64)
	if err != nil {
		errors.NewCustomError("POST /user/{user_id}/cart/{sku_id}: Invalid path params", http.StatusBadRequest, w)
		return
	}

	s.cartService.DeleteItemsByUser(userId)
	fmt.Fprint(w, http.StatusNoContent)
}
