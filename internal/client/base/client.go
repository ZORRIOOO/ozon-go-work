package httpclient

import (
	"bytes"
	"io/ioutil"
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
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (c *HttpClient) Post(url string, data []byte) (string, error) {
	resp, err := c.client.Post(url, "application/json", bytes.NewBuffer(data))
	if resp.StatusCode != http.StatusOK || err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
