package service

import (
	"encoding/json"
	"fmt"
	"homework/cart/internal/client/api/product/types"
	httpclient "homework/cart/internal/client/base/client"
)

type ProductService interface {
	GetProduct(request types.ProductRequest) (*types.ProductResponse, error)
	GetSkuList(request types.SkusRequest) (*types.SkusResponse, error)
}

type ProductServiceApi struct {
	client  *httpclient.HttpClient
	baseURL string
}

func NewProductServiceApi(client *httpclient.HttpClient, baseURL string) ProductService {
	return &ProductServiceApi{
		client:  client,
		baseURL: baseURL,
	}
}

func (s *ProductServiceApi) GetProduct(request types.ProductRequest) (*types.ProductResponse, error) {
	url := fmt.Sprintf("%s/get_product", s.baseURL)
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Post(url, requestBody)
	if err != nil {
		return nil, err
	}

	var productResponse types.ProductResponse
	if err := json.Unmarshal([]byte(resp), &productResponse); err != nil {
		return nil, err
	}

	return &productResponse, nil
}

func (s *ProductServiceApi) GetSkuList(request types.SkusRequest) (*types.SkusResponse, error) {
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
