package client

import (
	"bytes"
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

func (c *HttpClient) Get(url string) (string, error) {
	resp, err := c.client.Get(url)
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
		message := fmt.Sprintf("Invalid response body")
		return "", errors.New(message)
	}

	return string(body), nil
}

func (c *HttpClient) Post(url string, data []byte) (string, error) {
	resp, err := c.client.Post(url, "application/json", bytes.NewBuffer(data))
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
		return "", err
	}

	return string(body), nil
}
