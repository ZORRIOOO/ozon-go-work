package channel

import (
	"homework/cart/internal/client/api/product/types"
	"homework/cart/internal/pkg/cart/model"
	"sync"
)

type (
	ProductApi interface {
		GetProduct(request types.ProductRequest) (*types.ProductResponse, error)
		GetSkuList(request types.SkusRequest) (*types.SkusResponse, error)
	}

	CartChannel struct {
		productApi   ProductApi
		productToken string
	}
)

func NewCartChannel(productApi ProductApi, productToken string) *CartChannel {
	return &CartChannel{
		productApi:   productApi,
		productToken: productToken,
	}
}

func (channel CartChannel) FetchProductsInParallel(items []model.CartItem, userId model.UserId) ([]model.CartItem, uint32, error) {
	var wg sync.WaitGroup
	results := make(chan model.CartItem, len(items))
	errChan := make(chan error, len(items))

	for _, item := range items {
		wg.Add(1)
		go channel.FetchProduct(item, userId, &wg, results, errChan)
	}

	go func() {
		wg.Wait()
		close(results)
		close(errChan)
	}()

	return channel.CollectResults(results, errChan)
}

func (channel CartChannel) FetchProduct(item model.CartItem, userId model.UserId, wg *sync.WaitGroup, results chan<- model.CartItem, errChan chan<- error) {
	defer wg.Done()

	request := types.ProductRequest{
		Sku:   item.SKU,
		Token: channel.productToken,
	}
	product, err := channel.productApi.GetProduct(request)
	if err != nil {
		errChan <- err
		return
	}

	cartItem := model.CartItem{
		SKU:    item.SKU,
		Count:  item.Count,
		Name:   product.Name,
		Price:  product.Price,
		UserId: userId,
	}

	results <- cartItem
}

func (channel CartChannel) CollectResults(results <-chan model.CartItem, errChan <-chan error) ([]model.CartItem, uint32, error) {
	var responseItems []model.CartItem
	totalPrice := uint32(0)

	for {
		select {
		case result, ok := <-results:
			if !ok {
				results = nil
			} else {
				totalPrice += result.Price * uint32(result.Count)
				responseItems = append(responseItems, result)
			}
		case err := <-errChan:
			if err != nil {
				return nil, 0, err
			}
		}

		if results == nil && len(errChan) == 0 {
			break
		}
	}

	return responseItems, totalPrice, nil
}
