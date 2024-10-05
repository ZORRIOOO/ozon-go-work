package channel

import (
	"context"
	"golang.org/x/time/rate"
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
		limiter      *rate.Limiter
	}
)

func NewCartChannel(productApi ProductApi, productToken string, rpc int, maxRate int) *CartChannel {
	return &CartChannel{
		productApi:   productApi,
		productToken: productToken,
		limiter:      rate.NewLimiter(rate.Limit(rpc), maxRate),
	}
}

func (channel CartChannel) FetchProductsInParallel(items []model.CartItem, userId model.UserId) ([]model.CartItem, uint32, error) {
	group, ctx := errgroup.NewGroup(context.Background())
	results := make(chan model.CartItem, len(items))

	for _, item := range items {
		productItem := item
		group.Go(func() error {
			return channel.FetchProduct(ctx, productItem, userId, results)
		})
	}

	go func() {
		group.Wait()
		close(results)
	}()

	return channel.CollectResults(results, group.Wait())
}

func (channel CartChannel) FetchProduct(ctx context.Context, item model.CartItem, userId model.UserId, results chan<- model.CartItem) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if err := channel.limiter.Wait(ctx); err != nil {
		return err
	}

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
