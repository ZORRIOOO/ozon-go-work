package server

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"homework/cart/core/errors"
	"homework/cart/internal/pkg/cart/model"
	"io"
	"net/http"
	"strconv"
)

type AddCartItemRequest struct {
	SKU    int64  `json:"sku_id" validate:"required,min=1"`
	UserId int64  `json:"user_id" validate:"required,min=1"`
	Count  uint16 `json:"count" validate:"required,min=1"`
}

func (s *Server) AddCartItem(w http.ResponseWriter, r *http.Request) {
	rawUserId := r.PathValue("user_id")
	userId, err := strconv.ParseInt(rawUserId, 10, 64)
	rawSkuId := r.PathValue("sku_id")
	skuId, err := strconv.ParseInt(rawSkuId, 10, 64)

	createRequest := AddCartItemRequest{
		SKU:    skuId,
		UserId: userId,
	}
	body, err := io.ReadAll(r.Body)
	err = json.Unmarshal(body, &createRequest)
	if err != nil {
		errors.NewCustomError("POST /user/{user_id}/cart/{sku_id}: Invalid request body", http.StatusBadRequest, w)
		return
	}

	validate := validator.New()
	err = validate.Struct(createRequest)
	if err != nil {
		errors.NewCustomError("POST /user/{user_id}/cart/{sku_id}: Invalid request body", http.StatusBadRequest, w)
		return
	}

	cartParams := model.CartParameters{
		SKU:    createRequest.SKU,
		UserId: createRequest.UserId,
		Count:  createRequest.Count,
	}
	item, err, status := s.cartService.AddItem(cartParams)
	if err != nil {
		message := err.Error()
		errors.NewCustomError(message, status, w)
		return
	}

	if item.SKU == skuId {
		fmt.Fprint(w, http.StatusOK)
	}
}
