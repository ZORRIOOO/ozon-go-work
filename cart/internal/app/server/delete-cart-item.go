package server

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"homework/cart/core/errors"
	"homework/cart/internal/pkg/cart/model"
	"net/http"
	"strconv"
)

type DeleteCartItemRequest struct {
	SKU    int64 `json:"sku_id" validate:"required,min=1"`
	UserId int64 `json:"user_id" validate:"required,min=1"`
}

func (s *Server) DeleteCartItem(w http.ResponseWriter, r *http.Request) {
	rawUserId := r.PathValue("user_id")
	userId, _ := strconv.ParseInt(rawUserId, 10, 64)
	rawSkuId := r.PathValue("sku_id")
	skuId, _ := strconv.ParseInt(rawSkuId, 10, 64)

	deleteRequest := DeleteCartItemRequest{
		SKU:    skuId,
		UserId: userId,
	}
	validate := validator.New()
	err := validate.Struct(deleteRequest)
	if err != nil {
		message := fmt.Sprintf("DELETE /user/{user_id}/cart/{sku_id}: %s", errors.GetValidationErrMsg(err))
		errors.NewCustomError(message, http.StatusBadRequest, w)
		return
	}

	cartParams := model.DeleteCartParameters{
		SKU:    skuId,
		UserId: userId,
	}
	_, err = s.cartService.DeleteItem(cartParams)
	if err != nil {
		message := fmt.Sprintf("DELETE /user/{user_id}/cart/{sku_id}: %s", err.Error())
		errors.NewCustomError(message, http.StatusInternalServerError, w)
		return
	}

	fmt.Fprint(w, http.StatusOK)
}
