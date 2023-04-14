package form3

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	// defaultBaseURL is the default base URL for the Form3 API.
	defaultBaseURL = "https://api.form3.tech/v1"
)

type Form3BodyRequest[T any] struct {
	Data T `json:"data"`
}

type Form3BodyResponseLinks struct {
	Self  string `json:"self"`
	First string `json:"first,omitempty"`
	Last  string `json:"last,omitempty"`
	Prev  string `json:"prev,omitempty"`
	Next  string `json:"next,omitempty"`
}
type Form3BodyResponse[T any] struct {
	Data  T                      `json:"data,omitempty"`
	Links Form3BodyResponseLinks `json:"links,omitempty"`
}

type Form3BodyResponseError struct {
	ErrorMessage string `json:"error_message"`
}

type Form3APIError struct {
	StatusCode int
	Message    string
}

func (e *Form3APIError) Error() string {
	return fmt.Sprintf("Failed request with status code %d: %s", e.StatusCode, e.Message)
}

// Client handles communication with the Form3 API. It contains the underlying configuration and needed services.
type Client struct {
	// HTTP client used to make requests.
	client *http.Client

	// Base URL of the Form3 API.
	BaseURL *url.URL

	// Form3 services
	Account *AccountService
}

// NewClient returns a new Form3 API client.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, _ := url.Parse(defaultBaseURL)

	client := &Client{client: httpClient, BaseURL: baseURL}

	// attach services
	client.Account = &AccountService{client: client}

	return client
}

// Do sends HTTP API requests and returns the corresponding response or error.
func (c *Client) Do(ctx context.Context, method, url string, body, result interface{}) error {
	req, err := c.newRequest(ctx, method, url, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	res, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()

	err = c.decodeBody(res, result)

	if err != nil {
		return fmt.Errorf("failed to decode response body: %w", err)
	}

	return nil
}

// newRequest creates an HTTP request with the given method, URL and body (if any).
func (c *Client) newRequest(ctx context.Context, method, uri string, body interface{}) (*http.Request, error) {
	// First we parse the uri which includes the path and query parameters.
	parsedUri, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to parse uri: %w", err)
	}

	// We only keep the Path and RawQuery from the parsed uri and resolve it against the base URL.
	// This is needed because the uri might contain a full URL with a different host.s
	uriRef := url.URL{
		Path:     parsedUri.Path,
		RawQuery: parsedUri.RawQuery,
	}
	u := c.BaseURL.ResolveReference(&uriRef)

	var marshalledBody []byte

	// If the body is not nil, we marshal it to JSON.
	if body != nil {
		marshalledBody, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal body: %w", err)
		}
	}

	// We create the request with the given context.
	req, err := http.NewRequestWithContext(ctx, method, u.String(), bytes.NewReader(marshalledBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// We need to set the content-type to application/json if the body is not nil.
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

// decodeBody decodes the response body into the given result taking into account the status code and possible errors.
func (c *Client) decodeBody(res *http.Response, result interface{}) error {
	if res.StatusCode == http.StatusNoContent {
		return nil
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// If the status code is not 2xx, we try to decode the response body as an error.
	if res.StatusCode < 200 || res.StatusCode > 299 {
		if len(resBody) == 0 {
			// If the response body is empty, we return a generic error.
			return &Form3APIError{
				StatusCode: res.StatusCode,
				Message:    http.StatusText(res.StatusCode),
			}
		}

		var errorResult Form3BodyResponseError
		err = json.Unmarshal(resBody, &errorResult)
		if err != nil {
			return fmt.Errorf("failed to decode error response body: %w", err)
		}

		// If the response body is not empty, we return the error message that we received in the response.
		return &Form3APIError{
			StatusCode: res.StatusCode,
			Message:    errorResult.ErrorMessage,
		}
	}

	return json.Unmarshal(resBody, result)
}

// Generic helper to convert a value to a pointer of the same type.
func ToPointer[T any](value T) *T {
	return &value
}
