package service

import (
	"encoding/json"
	"fmt"
	"homework/cart/internal/client/api/loms/types"
	httpclient "homework/cart/internal/client/base/client"
)

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

	var orderId types.OrderId
	if err := json.Unmarshal([]byte(resp), &orderId); err != nil {
		return nil, err
	}
	return &types.OrderCreateResponse{OrderId: orderId}, nil
}

func (l LomsServiceApi) StocksInfo(request types.StocksInfoRequest) (*types.StocksInfoResponse, error) {
	url := fmt.Sprintf("%s/v1/loms/stock/%d/info", l.baseURL, request.Sku)
	resp, err := l.client.Get(url)
	if err != nil {
		return nil, err
	}

	var stocksInfoResponse types.StocksInfoResponse
	if err := json.Unmarshal([]byte(resp), &stocksInfoResponse); err != nil {
		return nil, err
	}
	return &stocksInfoResponse, nil
}
