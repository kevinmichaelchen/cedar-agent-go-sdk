package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	c   *http.Client
	cfg config
}

type config struct {
	baseURL string
}

type Option func(*config)

func WithBaseURL(baseURL string) Option {
	return func(c *config) {
		c.baseURL = baseURL
	}
}

func NewClient(c *http.Client, opts ...Option) *Client {
	cfg := config{
		baseURL: "http://localhost:8180",
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	return &Client{c: c, cfg: cfg}
}

type CheckRequest struct {
	Principal string `json:"principal"`
	Action    string `json:"action"`
	Resource  string `json:"resource"`
}

type CheckResponse struct {
	Decision    string `json:"decision"`
	Diagnostics struct {
		Errors []interface{} `json:"errors"`
		Reason []string      `json:"reason"`
	} `json:"diagnostics"`
}

func (c Client) Check(_ context.Context, payload CheckRequest) (bool, error) {
	reqBody, err := json.Marshal(payload)
	if err != nil {
		return false, fmt.Errorf("unable to marshal http request body: %w", err)
	}

	url := fmt.Sprintf("%s/v1/is_authorized", c.cfg.baseURL)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(reqBody))
	if err != nil {
		return false, fmt.Errorf("unable to build http request: %w", err)
	}

	res, err := c.c.Do(req)
	if err != nil {
		return false, fmt.Errorf("unable to execute http request: %w", err)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return false, fmt.Errorf("unable to read http response: %w", err)
	}

	var response CheckResponse

	err = json.Unmarshal(resBody, &response)
	if err != nil {
		return false, fmt.Errorf("unable to unmarshal http response: %w", err)
	}

	return response.Decision == "Allow", nil
}
