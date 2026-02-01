package accrual

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Response struct {
	Order   string   `json:"order"`
	Status  string   `json:"status"`
	Accrual *float64 `json:"accrual,omitempty"`
}

type Client struct {
	baseURL string
	http    *http.Client
}

func NewClient(addr string) *Client {
	return &Client{
		baseURL: addr,
		http: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

var (
	ErrNotFound   = errors.New("not registered")
	ErrTooManyReq = errors.New("rate limit")
)

func (c *Client) Get(ctx context.Context, number string) (*Response, error) {
	url := fmt.Sprintf("%s/api/orders/%s", c.baseURL, number)

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var r Response
		return &r, json.NewDecoder(resp.Body).Decode(&r)

	case http.StatusNoContent:
		return nil, ErrNotFound

	case http.StatusTooManyRequests:
		return nil, ErrTooManyReq

	default:
		return nil, errors.New("accrual error")
	}
}
