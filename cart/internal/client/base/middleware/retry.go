package client

import (
	"fmt"
	"net/http"
	"time"
)

func RetryMiddleware(next http.RoundTripper, retries int, statuses []int) http.RoundTripper {
	return &retryRoundTripper{
		next:     next,
		retries:  retries,
		statuses: statuses,
	}
}

type retryRoundTripper struct {
	next     http.RoundTripper
	retries  int
	statuses []int
}

func (rt *retryRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	for i := 0; i < rt.retries; i++ {
		resp, err = rt.next.RoundTrip(req)
		fmt.Println("RoundTrip", resp)

		if err != nil {
			return nil, err
		}

		if !contains(rt.statuses, resp.StatusCode) {
			return resp, nil
		}

		resp.Body.Close()

		time.Sleep(time.Second * 2)
	}

	message := fmt.Sprintf("Exceeded maximum retries: Last status was %d", rt.retries)
	return resp, fmt.Errorf(message)
}

func contains(statuses []int, status int) bool {
	for _, s := range statuses {
		if s == status {
			return true
		}
	}
	return false
}
