package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// APIClient is a thin wrapper around net/http for calling the
// order-management REST API. It injects the caller's JWT (read from the
// browser cookie by the handler) as a Bearer token on every request.
type APIClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// APIError represents a non-2xx response from the backend, carrying the
// status code and whatever error message the API returned so handlers can
// show something meaningful to the user.
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("api: %d %s", e.StatusCode, e.Message)
}

type errorBody struct {
	Error string `json:"error"`
}

// do executes an HTTP request against the API and decodes a JSON response
// body into out (if out is non-nil and the response has a body).
func (c *APIClient) do(method, path, token string, body interface{}, out interface{}) error {
	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, c.BaseURL+path, reqBody)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		msg := resp.Status
		var eb errorBody
		if json.Unmarshal(respBytes, &eb) == nil && eb.Error != "" {
			msg = eb.Error
		}
		return &APIError{StatusCode: resp.StatusCode, Message: msg}
	}

	if out != nil && len(respBytes) > 0 {
		if err := json.Unmarshal(respBytes, out); err != nil {
			return err
		}
	}

	return nil
}

func (c *APIClient) Get(path, token string, out interface{}) error {
	return c.do(http.MethodGet, path, token, nil, out)
}

func (c *APIClient) Post(path, token string, body interface{}, out interface{}) error {
	return c.do(http.MethodPost, path, token, body, out)
}

func (c *APIClient) Put(path, token string, body interface{}, out interface{}) error {
	return c.do(http.MethodPut, path, token, body, out)
}

func (c *APIClient) Delete(path, token string) error {
	return c.do(http.MethodDelete, path, token, nil, nil)
}

// IsUnauthorized reports whether err came back as a 401 from the API,
// meaning the caller's session token is missing/expired.
func IsUnauthorized(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusUnauthorized
	}
	return false
}
