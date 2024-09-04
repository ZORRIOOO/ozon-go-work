package httpclient

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type HttpClient struct {
	client *http.Client
}

func NewHttpClient(timeout time.Duration) *HttpClient {
	return &HttpClient{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *HttpClient) Get(url string) (string, error) {
	resp, err := c.client.Get(url)
	if resp.StatusCode != http.StatusOK || err != nil {
		body, _ := io.ReadAll(resp.Body)
		message := fmt.Sprintf(string(body))
		return "", errors.New(message)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		message := fmt.Sprintf("Invalid response body")
		return "", errors.New(message)
	}

	return string(body), nil
}

func (c *HttpClient) Post(url string, data []byte) (string, error) {
	resp, err := c.client.Post(url, "application/json", bytes.NewBuffer(data))
	if resp.StatusCode != http.StatusOK || err != nil {
		body, _ := io.ReadAll(resp.Body)
		message := fmt.Sprintf(string(body))
		return "", errors.New(message)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		message := fmt.Sprintf("Invalid response body")
		return "", errors.New(message)
	}

	return string(body), nil
}
