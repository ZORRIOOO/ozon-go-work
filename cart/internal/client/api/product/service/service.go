package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"homework/cart/internal/client/api/product/types"
	httpclient "homework/cart/internal/client/base"
)

type ProductService interface {
	GetProduct(request types.ProductRequest) (*types.ProductResponse, error)
	GetSkuList(request types.SkusRequest) (*types.SkusResponse, error)
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

func (s *productService) GetProduct(request types.ProductRequest) (*types.ProductResponse, error) {
	url := fmt.Sprintf("%s/get_product", s.baseURL)
	requestBody, err := json.Marshal(request)
	if err != nil {
		message := fmt.Sprintf("POST /get_product: Invalid request body")
		return nil, errors.New(message)
	}

	resp, err := s.client.Post(url, requestBody)
	if err != nil {
		message := fmt.Sprintf("POST /get_product: %s", err.Error())
		return nil, errors.New(message)
	}

	var productResponse types.ProductResponse
	if err := json.Unmarshal([]byte(resp), &productResponse); err != nil {
		message := fmt.Sprintf("POST /get_product: Invalid response body")
		return nil, errors.New(message)
	}

	return &productResponse, nil
}

func (s *productService) GetSkuList(request types.SkusRequest) (*types.SkusResponse, error) {
	url := fmt.Sprintf("%s/list_skus", s.baseURL)

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Post(url, requestBody)
	if err != nil {
		return nil, err
	}

	var skusResponse types.SkusResponse
	if err := json.Unmarshal([]byte(resp), &skusResponse); err != nil {
		return nil, err
	}

	return &skusResponse, nil
}
