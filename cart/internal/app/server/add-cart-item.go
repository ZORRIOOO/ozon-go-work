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
	userId, _ := strconv.ParseInt(rawUserId, 10, 64)
	rawSkuId := r.PathValue("sku_id")
	skuId, _ := strconv.ParseInt(rawSkuId, 10, 64)

	createRequest := AddCartItemRequest{
		SKU:    skuId,
		UserId: userId,
	}
	body, err := io.ReadAll(r.Body)
	err = json.Unmarshal(body, &createRequest)
	if err != nil {
		message := fmt.Sprintf("POST /user/{user_id}/cart/{sku_id}: %s", err.Error())
		errors.NewCustomError(message, http.StatusBadRequest, w)
		return
	}

	validate := validator.New()
	err = validate.Struct(createRequest)
	if err != nil {
		message := fmt.Sprintf("POST /user/{user_id}/cart/{sku_id}: %s", errors.GetValidationErrMsg(err))
		errors.NewCustomError(message, http.StatusBadRequest, w)
		return
	}

	cartParams := model.CartParameters{
		SKU:    createRequest.SKU,
		UserId: createRequest.UserId,
		Count:  createRequest.Count,
	}
	_, err = s.addItemHandler.AddItem(cartParams)
	if err != nil {
		message := fmt.Sprintf("POST /user/{user_id}/cart/{sku_id}: %s", err.Error())
		statusCode := errors.GetStatusCode(message)
		errors.NewCustomError(message, statusCode, w)
		return
	}

	fmt.Fprint(w, http.StatusOK)
}
