package server

import (
	"encoding/json"
	"fmt"
	"homework/cart/core/errors"
	"homework/cart/internal/pkg/cart/model"
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
	if userId <= 0 || skuId <= 0 || err != nil {
		errors.NewCustomError("POST /user/{user_id}/cart/{sku_id}: Invalid path params", http.StatusBadRequest, w)
		return
	}

	body, err := io.ReadAll(r.Body)
	var createRequest AddCartItemRequest
	err = json.Unmarshal(body, &createRequest)
	if err != nil {
		errors.NewCustomError("POST /user/{user_id}/cart/{sku_id}: Invalid request body", http.StatusBadRequest, w)
		return
	}

	if createRequest.Count <= 0 {
		errors.NewCustomError("POST /user/{user_id}/cart/{sku_id}: Invalid request body", http.StatusBadRequest, w)
		return
	}

	cartParams := model.CartParameters{
		SKU:    skuId,
		UserId: userId,
		Count:  createRequest.Count,
	}
	item, err := s.cartService.AddItem(cartParams)
	if err != nil {
		message := err.Error()
		errors.NewCustomError(message, http.StatusInternalServerError, w)
		return
	}

	if item.SKU == skuId {
		fmt.Fprint(w, http.StatusOK)
	}
}
