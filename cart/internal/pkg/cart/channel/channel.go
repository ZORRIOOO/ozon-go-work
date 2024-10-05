package channel

import (
	"homework/cart/core/errgroup"
	"homework/cart/internal/client/api/product/types"
	"homework/cart/internal/pkg/cart/model"
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
	results := make(chan model.CartItem, len(items))

	var g errgroup.Group

	for _, item := range items {
		productItem := item
		g.Go(func() error {
			return channel.FetchProduct(productItem, userId, results)
		})
	}

	go func() {
		g.Wait()
		close(results)
	}()

	return channel.CollectResults(results, g.Wait())
}

func (channel CartChannel) FetchProduct(item model.CartItem, userId model.UserId, results chan<- model.CartItem) error {
	request := types.ProductRequest{
		Sku:   item.SKU,
		Token: channel.productToken,
	}
	product, err := channel.productApi.GetProduct(request)
	if err != nil {
		return err
	}

	cartItem := model.CartItem{
		SKU:    item.SKU,
		Count:  item.Count,
		Name:   product.Name,
		Price:  product.Price,
		UserId: userId,
	}

	results <- cartItem
	return nil
}

func (channel CartChannel) CollectResults(results <-chan model.CartItem, err error) ([]model.CartItem, uint32, error) {
	var responseItems []model.CartItem
	totalPrice := uint32(0)

	for result := range results {
		totalPrice += result.Price * uint32(result.Count)
		responseItems = append(responseItems, result)
	}

	if err != nil {
		return nil, 0, err
	}

	return responseItems, totalPrice, nil
}
