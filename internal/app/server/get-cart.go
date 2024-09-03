package server

import (
	"cart/core/errors"
	"encoding/json"
	"net/http"
	"strconv"
)

func (s *Server) GetCartByUser(w http.ResponseWriter, r *http.Request) {
	rawUserId := r.PathValue("user_id")
	userId, err := strconv.ParseInt(rawUserId, 10, 64)
	if err != nil {
		errors.NewCustomError("POST /user/{user_id}/cart/{sku_id}: Invalid path params", http.StatusBadRequest, w)
		return
	}

	response, err := s.cartService.GetCartByUser(userId)

	rawResponse, err := json.Marshal(response)
	if err != nil {
		errors.NewCustomError("POST /user/{user_id}/cart/{sku_id}: Invalid response", http.StatusBadRequest, w)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(rawResponse)
}
