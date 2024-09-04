package server

import (
	"cart/core/errors"
	"cart/internal/pkg/cart/model"
	"fmt"
	"net/http"
	"strconv"
)

func (s *Server) DeleteCartItem(w http.ResponseWriter, r *http.Request) {
	rawUserId := r.PathValue("user_id")
	userId, err := strconv.ParseInt(rawUserId, 10, 64)
	rawSkuId := r.PathValue("sku_id")
	skuId, err := strconv.ParseInt(rawSkuId, 10, 64)
	if userId <= 0 || skuId <= 0 || err != nil {
		errors.NewCustomError("DELETE /user/{user_id}/cart/{sku_id}: Invalid path params", http.StatusBadRequest, w)
		return
	}

	cartParams := model.CartParameters{
		SKU:    skuId,
		UserId: userId,
	}
	item, err := s.cartService.DeleteItem(cartParams)
	if err != nil {
		message := fmt.Sprintf("DELETE /user/{user_id}/cart/{sku_id}: %s", err.Error())
		errors.NewCustomError(message, http.StatusInternalServerError, w)
		return
	}

	if item.SKU == skuId {
		fmt.Fprint(w, http.StatusNoContent)
	}
}
