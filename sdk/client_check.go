package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Principal string
type Action string
type Resource string

// CheckRequest - Provides a principal, action, and resource for Cedar Agent to
// evaluate an authorization decision.
type CheckRequest struct {
	// Principal is typically the user.
	Principal Principal `json:"principal"`
	// Action is the verb (e.g., read, write, update, view, list)
	Action Action `json:"action"`
	// Resource is the thing being accessed.
	Resource Resource `json:"resource"`
}

// Decision - An authorization decision.
type Decision struct {
	// Allowed - Whether the principal is allowed to perform the action on the
	// resource.
	Allowed bool
	// Diagnostics - Insight into how the decision was reached.
	Diagnostics Diagnostics
}

// IsAuthorizedResponse - HTTP Response from Cedar Agent.
type IsAuthorizedResponse struct {
	// Decision - Whether the principal is allowed to perform the action on the
	// resource.
	Decision string `json:"decision"`
	// Diagnostics - Insight into how the decision was reached.
	Diagnostics Diagnostics `json:"diagnostics"`
}

// Diagnostics - Sheds light into the evaluation decision.
type Diagnostics struct {
	// Errors that occurred during the evaluation decisions.
	Errors []interface{} `json:"errors"`
	// Reason that the decision was made (e.g., policies involved).
	Reason []string `json:"reason"`
}

// Check - Performs an authorization request using Cedar Agent.
func (c Client) Check(_ context.Context, payload CheckRequest) (*Decision, error) {
	// Build request
	req, err := buildCheckRequest(payload, c.cfg.baseURL)
	if err != nil {
		return nil, err
	}

	// Execute request
	res, err := c.c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to execute http request: %w", err)
	}

	var response IsAuthorizedResponse

	err = unmarshalResponse(res, &response)
	if err != nil {
		return nil, err
	}

	return &Decision{
		Allowed:     response.Decision == "Allow",
		Diagnostics: response.Diagnostics,
	}, nil
}

func unmarshalResponse(res *http.Response, target any) error {
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("unable to read http response: %w", err)
	}

	err = json.Unmarshal(resBody, &target)
	if err != nil {
		return fmt.Errorf("unable to unmarshal http response: %w", err)
	}

	return nil
}

func buildCheckRequest(payload CheckRequest, baseURL string) (*http.Request, error) {
	reqBody, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal http request body: %w", err)
	}

	url := fmt.Sprintf("%s/v1/is_authorized", baseURL)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("unable to build http request: %w", err)
	}

	return req, nil
}
