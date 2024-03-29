package integration

import (
	"net/url"
	"os"
	"testing"

	"github.com/agatticelli/form3-client-go/form3"
)

const defaultTestingBaseURL = "http://localhost:8080/v1/"

var client *form3.Client

func initClient(t *testing.T) *form3.Client {
	var err error

	testingBaseURL := os.Getenv("FORM3_API_BASE_URL")
	if testingBaseURL == "" {
		testingBaseURL = defaultTestingBaseURL
	}

	client = form3.NewClient(nil)
	client.BaseURL, err = url.Parse(testingBaseURL)

	if err != nil {
		t.Fatalf("error setting up client: %v", err)
	}

	return client
}
