package service

import (
	"encoding/json"
	"fmt"
	"homework/cart/internal/client/api/loms/types"
	httpclient "homework/cart/internal/client/base/client"
)

func (l LomsServiceApi) CreateOrder(request types.OrderCreateRequest) (*types.OrderCreateResponse, error) {
	items := request.Items
	url := fmt.Sprintf("%s/v1/loms/user/%d/order/create", l.baseURL, request.User)
	requestBody, err := json.Marshal(items)
	if err != nil {
		return nil, err
	}

	resp, err := l.client.Post(url, requestBody)
	if err != nil {
		return nil, err
	}

	var orderCreateResponse types.OrderCreateResponse
	if err := json.Unmarshal([]byte(resp), &orderCreateResponse); err != nil {
		return nil, err
	}
	return &orderCreateResponse, nil
}

type LomsServiceApi struct {
	client  *httpclient.HttpClient
	baseURL string
}

func NewLomsServiceApi(client *httpclient.HttpClient, baseURL string) *LomsServiceApi {
	return &LomsServiceApi{
		client:  client,
		baseURL: baseURL,
	}
}
