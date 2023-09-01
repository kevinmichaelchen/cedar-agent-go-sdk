package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// CheckRequest - Provides a principal, action, and resource for Cedar Agent to
// evaluate an authorization decision.
type CheckRequest struct {
	Principal string `json:"principal"`
	Action    string `json:"action"`
	Resource  string `json:"resource"`
}

// CheckResponse - An authorization decision.
type CheckResponse struct {
	Decision    string `json:"decision"`
	Diagnostics struct {
		Errors []interface{} `json:"errors"`
		Reason []string      `json:"reason"`
	} `json:"diagnostics"`
}

// Check - Performs an authorization request using Cedar Agent.
func (c Client) Check(_ context.Context, payload CheckRequest) (*CheckResponse, error) {
	reqBody, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal http request body: %w", err)
	}

	url := fmt.Sprintf("%s/v1/is_authorized", c.cfg.baseURL)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("unable to build http request: %w", err)
	}

	res, err := c.c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to execute http request: %w", err)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read http response: %w", err)
	}

	var response CheckResponse

	err = json.Unmarshal(resBody, &response)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal http response: %w", err)
	}

	return &response, nil
}
