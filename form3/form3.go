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

type Client struct {
	// HTTP client used to make requests.
	client *http.Client

	// Base URL of the Form3 API.
	BaseURL *url.URL
}

func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, _ := url.Parse(defaultBaseURL)

	return &Client{client: httpClient, BaseURL: baseURL}
}

func (c *Client) Do(ctx context.Context, method, url string, body, result interface{}) error {
	req, err := c.newRequest(ctx, method, url, body)
	if err != nil {
		return err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return c.decodeBody(res, result)
}

func (c *Client) newRequest(ctx context.Context, method, uri string, body interface{}) (*http.Request, error) {
	parsedUri, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	uriRef := url.URL{
		Path:     parsedUri.Path,
		RawQuery: parsedUri.RawQuery,
	}
	u := c.BaseURL.ResolveReference(&uriRef)

	var marshalledBody []byte

	if body != nil {
		marshalledBody, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), bytes.NewReader(marshalledBody))
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

func (c *Client) decodeBody(res *http.Response, result interface{}) error {
	if res.StatusCode == http.StatusNoContent {
		return nil
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		if len(resBody) == 0 {
			return &Form3APIError{
				StatusCode: res.StatusCode,
				Message:    http.StatusText(res.StatusCode),
			}
		}

		var errorResult Form3BodyResponseError
		err = json.Unmarshal(resBody, &errorResult)
		if err != nil {
			return err
		}

		return &Form3APIError{
			StatusCode: res.StatusCode,
			Message:    errorResult.ErrorMessage,
		}
	}

	return json.Unmarshal(resBody, result)
}
