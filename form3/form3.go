package form3

import (
	"fmt"
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
