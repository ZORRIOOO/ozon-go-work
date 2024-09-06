package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"homework/cart/internal/client/api/product/types"
	httpclient "homework/cart/internal/client/base/client"
	"net/http"
)

type ProductService interface {
	GetProduct(request types.ProductRequest) (*types.ProductResponse, error, int)
	GetSkuList(request types.SkusRequest) (*types.SkusResponse, error, int)
}

type productService struct {
	client  *httpclient.HttpClient
	baseURL string
}

func NewProductService(client *httpclient.HttpClient, baseURL string) ProductService {
	return &productService{
		client:  client,
		baseURL: baseURL,
	}
}

func (s *productService) GetProduct(request types.ProductRequest) (*types.ProductResponse, error, int) {
	url := fmt.Sprintf("%s/get_product", s.baseURL)
	requestBody, err := json.Marshal(request)
	if err != nil {
		message := fmt.Sprintf("POST /get_product: Invalid request body")
		return nil, errors.New(message), http.StatusBadRequest
	}

	resp, err, status := s.client.Post(url, requestBody)
	if err != nil {
		message := fmt.Sprintf("POST /get_product: %s", err.Error())
		return nil, errors.New(message), status
	}

	var productResponse types.ProductResponse
	if err := json.Unmarshal([]byte(resp), &productResponse); err != nil {
		message := fmt.Sprintf("POST /get_product: Invalid response body")
		return nil, errors.New(message), http.StatusInternalServerError
	}

	return &productResponse, nil, http.StatusOK
}

func (s *productService) GetSkuList(request types.SkusRequest) (*types.SkusResponse, error, int) {
	url := fmt.Sprintf("%s/list_skus", s.baseURL)

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, err, http.StatusBadRequest
	}

	resp, err, status := s.client.Post(url, requestBody)
	if err != nil {
		return nil, err, status
	}

	var skusResponse types.SkusResponse
	if err := json.Unmarshal([]byte(resp), &skusResponse); err != nil {
		return nil, err, http.StatusInternalServerError
	}

	return &skusResponse, nil, http.StatusOK
}
