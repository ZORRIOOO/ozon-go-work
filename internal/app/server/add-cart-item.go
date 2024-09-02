package server

import (
	"cart/core/errors"
	"cart/internal/pkg/cart/model"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

type AddCartItemRequest struct {
	Count uint16 `json:"count"`
}

func (s *Server) AddCartItem(w http.ResponseWriter, r *http.Request) {
	rawUserId := r.PathValue("user_id")
	userId, err := strconv.ParseInt(rawUserId, 10, 64)

	rawSkuId := r.PathValue("sku_id")
	skuId, err := strconv.ParseInt(rawSkuId, 10, 64)
	if err != nil {
		errors.NewCustomError(w, "POST /user/{user_id}/cart/{sku_id}", http.StatusBadRequest, err)
		return
	}

	body, err := io.ReadAll(r.Body)
	var createRequest AddCartItemRequest
	err = json.Unmarshal(body, &createRequest)
	if err != nil {
		errors.NewCustomError(w, "POST /user/{user_id}/cart/{sku_id}", http.StatusBadRequest, err)
		return
	}

	if createRequest.Count != 0 {
		errors.NewCustomError(w, "POST /user/{user_id}/cart/{sku_id}", http.StatusBadRequest, err)
		return
	}

	cartItem := model.CartItem{
		SKU:    skuId,
		Name:   "",
		Count:  createRequest.Count,
		Price:  0,
		UserId: userId,
	}

	if createRequest.Count != 0 {
		errors.NewCustomError(w, "POST /user/{user_id}/cart/{sku_id}", http.StatusInternalServerError, "Invalid arguments")
	}

	_, err = s.cartService.AddItem(r.Context(), cartItem)
	w.WriteHeader(http.StatusOK)
}
