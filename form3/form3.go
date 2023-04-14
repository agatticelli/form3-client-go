package form3

import (
	"net/http"
	"net/url"
)

const (
	// defaultBaseURL is the default base URL for the Form3 API.
	defaultBaseURL = "https://api.form3.tech/v1"
)

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
