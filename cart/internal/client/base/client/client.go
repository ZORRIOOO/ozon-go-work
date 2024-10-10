package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	middleware "homework/cart/internal/client/base/middleware"
	"io"
	"net/http"
	"time"
)

type HttpClient struct {
	client *http.Client
}

func NewHttpClient(timeout time.Duration, retries int, statusList []int) *HttpClient {
	client := &http.Client{
		Timeout:   timeout,
		Transport: middleware.RetryMiddleware(http.DefaultTransport, retries, statusList),
	}
	return &HttpClient{client: client}
}

func (c *HttpClient) Get(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating GET request: %v", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		message := fmt.Sprintf("API=ProductService, Method=%s, URL=%s, Status=%v", resp.Request.Method, resp.Request.URL, resp.StatusCode)
		return "", errors.New(message)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("invalid response body: %v", err)
	}

	return string(body), nil
}

func (c *HttpClient) Post(ctx context.Context, url string, data []byte) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return "", fmt.Errorf("error creating POST request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making POST request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		message := fmt.Sprintf("API=ProductService, Method=%s, URL=%s, Status=%v", resp.Request.Method, resp.Request.URL, resp.StatusCode)
		return "", errors.New(message)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("invalid response body: %v", err)
	}

	return string(body), nil
}
