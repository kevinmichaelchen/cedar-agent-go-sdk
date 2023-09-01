package sdk

import (
	"net/http"
)

const (
	defaultBaseURL               = "http://localhost:8180"
	defaultParallelizationFactor = 3
)

// Client - A client for interacting with Cedar Agent.
type Client struct {
	c   *http.Client
	cfg config
}

type config struct {
	baseURL    string
	numWorkers int
}

// NewClient - Creates a new client for interacting with Cedar Agent.
func NewClient(c *http.Client, opts ...Option) *Client {
	cfg := config{
		baseURL:    defaultBaseURL,
		numWorkers: defaultParallelizationFactor,
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	return &Client{c: c, cfg: cfg}
}

// Option - A functional option for configuring your Cedar Agent client.
type Option func(*config)

// WithBaseURL - Configures Cedar Agent's base URL.
// By default, we assume it runs on http://localhost:8180.
func WithBaseURL(baseURL string) Option {
	return func(c *config) {
		c.baseURL = baseURL
	}
}

// WithParallelizationFactor - Configures the client's parallelization factor
// when performing a batch of authorization requests in parallel.
func WithParallelizationFactor(numWorkers int) Option {
	return func(c *config) {
		c.numWorkers = numWorkers
	}
}
