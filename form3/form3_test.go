package form3

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"
)

func Test_ToPointer(t *testing.T) {
	// Test with an integer value
	var intValue int = 42
	var intPointer *int = ToPointer(intValue)
	if *intPointer != intValue {
		t.Fatalf("ToPointer did not return a pointer to the correct integer value. Expected %v but got %v", intValue, *intPointer)
	}

	// Test with a string value
	var strValue string = "Hello, world!"
	var strPointer *string = ToPointer(strValue)
	if *strPointer != strValue {
		t.Fatalf("ToPointer did not return a pointer to the correct string value. Expected %v but got %v", strValue, *strPointer)
	}

	// Test with a boolean value
	var boolValue bool = true
	var boolPointer *bool = ToPointer(boolValue)
	if *boolPointer != boolValue {
		t.Fatalf("ToPointer did not return a pointer to the correct boolean value. Expected %v but got %v", boolValue, *boolPointer)
	}

	// Test with a nil value
	var errValue error
	var errPointer *error = ToPointer(errValue)
	if *errPointer != nil {
		t.Fatalf("ToPointer did not return a nil pointer for a nil value. Expected %v but got %v", nil, errPointer)
	}
}

func TestClient_newRequest(t *testing.T) {
	type fields struct {
		BaseURL *url.URL
	}
	type args struct {
		method string
		uri    string
		body   interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *http.Request
		wantErr bool
	}{
		{
			name: "Test simple GET with path suffix and query params",
			fields: fields{
				BaseURL: &url.URL{Scheme: "https", Host: "api.form3.tech", Path: "/v1/"},
			},
			args: args{
				method: "GET",
				uri:    "accounts?version=42&organisation_id=1234",
				body:   nil,
			},
			want: &http.Request{
				Method: "GET",
				URL: &url.URL{
					Scheme:   "https",
					Host:     "api.form3.tech",
					Path:     "/v1/accounts",
					RawQuery: "version=42&organisation_id=1234",
				},
				Header: http.Header{},
			},
		},
		{
			name: "Test simple GET with absolute path",
			fields: fields{
				BaseURL: &url.URL{Scheme: "https", Host: "api.form3.tech", Path: "/v1/"},
			},
			args: args{
				method: "GET",
				uri:    "/v2/accounts",
				body:   nil,
			},
			want: &http.Request{
				Method: "GET",
				URL: &url.URL{
					Scheme:   "https",
					Host:     "api.form3.tech",
					Path:     "/v2/accounts",
					RawQuery: "",
				},
				Header: http.Header{},
			},
		},
		{
			name: "Test simple GET with a new absolute URL as path",
			fields: fields{
				BaseURL: &url.URL{Scheme: "https", Host: "api.form3.tech", Path: "/v1/"},
			},
			args: args{
				method: "GET",
				uri:    "https://fake.domain/v3/accounts",
				body:   nil,
			},
			want: &http.Request{
				Method: "GET",
				URL: &url.URL{
					Scheme:   "https",
					Host:     "api.form3.tech",
					Path:     "/v3/accounts",
					RawQuery: "",
				},
				Header: http.Header{},
			},
		},
		{
			name: "Test simple POST with nil body",
			fields: fields{
				BaseURL: &url.URL{Scheme: "https", Host: "api.form3.tech", Path: "/v1/"},
			},
			args: args{
				method: "POST",
				uri:    "accounts",
				body:   nil,
			},
			want: &http.Request{
				Method: "POST",
				URL: &url.URL{
					Scheme: "https",
					Host:   "api.form3.tech",
					Path:   "/v1/accounts",
				},
				Header: http.Header{},
			},
		},
		{
			name: "Test simple POST with body",
			fields: fields{
				BaseURL: &url.URL{Scheme: "https", Host: "api.form3.tech", Path: "/v1/"},
			},
			args: args{
				method: "POST",
				uri:    "accounts",
				body:   map[string]interface{}{"organisationID": "1234"},
			},
			want: &http.Request{
				Method: "POST",
				URL: &url.URL{
					Scheme: "https",
					Host:   "api.form3.tech",
					Path:   "/v1/accounts",
				},
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
				Body: io.NopCloser(strings.NewReader(`{"organisationID":"1234"}`)),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				BaseURL: tt.fields.BaseURL,
			}
			got, err := c.newRequest(context.Background(), tt.args.method, tt.args.uri, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Client.newRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Compare URL
			if !reflect.DeepEqual(got.URL, tt.want.URL) {
				t.Fatalf("Client.newRequest() - URL - got = %v, want %v", got.URL, tt.want.URL)
			}

			// Compare Method
			if got.Method != tt.want.Method {
				t.Fatalf("Client.newRequest() - Method - got = %v, want %v", got.Method, tt.want.Method)
			}

			// Compare Headers
			if !reflect.DeepEqual(got.Header, tt.want.Header) {
				t.Fatalf("Client.newRequest() - Headers - got = %v, want %v", got.Header, tt.want.Header)
			}

			// Compare Raw Query
			if got.URL.RawQuery != tt.want.URL.RawQuery {
				t.Fatalf("Client.newRequest() - RawQuery - got = %v, want %v", got.URL.RawQuery, tt.want.URL.RawQuery)
			}

			// Compare Body
			if tt.want.Body != nil {
				gotBody, err := io.ReadAll(got.Body)
				if err != nil {
					t.Fatalf("Client.newRequest() - got.Body read - got = %v, want %v", err, nil)
				}

				wantBody, err := io.ReadAll(tt.want.Body)
				if err != nil {
					t.Fatalf("Client.newRequest() - want.Body read - got = %v, want %v", err, nil)
				}

				if !reflect.DeepEqual(gotBody, wantBody) {
					t.Fatalf("Client.newRequest() - Body - got = %v, want %v", string(gotBody), string(wantBody))
				}
			}
		})
	}
}

func TestClient_decodeBody(t *testing.T) {
	type successResponse struct {
		Version *int `json:"version,omitempty"`
	}
	type args struct {
		res    *http.Response
		result interface{}
	}
	tests := []struct {
		name       string
		args       args
		wantResult interface{}
		wantErr    interface{}
	}{
		{
			name: "Test decodeBody with statusCode 204",
			args: args{
				res: &http.Response{
					StatusCode: 204,
				},
				result: nil,
			},
			wantResult: nil,
			wantErr:    nil,
		},
		{
			name: "Test decodeBody with success code and data",
			args: args{
				res: &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(`{"version": 1}`)),
				},
				result: &successResponse{},
			},
			wantResult: successResponse{Version: ToPointer(1)},
			wantErr:    nil,
		},
		{
			name: "Test decodeBody with error code and body data",
			args: args{
				res: &http.Response{
					StatusCode: 409,
					Body:       io.NopCloser(strings.NewReader(`{"error_message": "Account cannot be created as it violates a duplicate constraint"}`)),
				},
				result: nil,
			},
			wantResult: nil,
			wantErr:    Form3APIError{StatusCode: 409, Message: "Account cannot be created as it violates a duplicate constraint"},
		},
		{
			name: "Test decodeBody with error code without data",
			args: args{
				res: &http.Response{
					StatusCode: 400,
					Body:       io.NopCloser(strings.NewReader("")),
				},
				result: nil,
			},
			wantResult: nil,
			wantErr:    Form3APIError{StatusCode: 400, Message: http.StatusText(400)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{}
			err := c.decodeBody(tt.args.res, tt.args.result)
			if err != nil && tt.wantErr == nil {
				t.Fatalf("Client.decodeBody() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantResult != nil {
				res, ok := tt.args.result.(*successResponse)
				if !ok {
					t.Fatalf("Client.decodeBody() - result type - got = %v, want %v", reflect.TypeOf(tt.args.result), reflect.TypeOf(&successResponse{}))
				}

				wantRes, ok := tt.wantResult.(successResponse)
				if !ok {
					t.Fatalf("Client.decodeBody() - wantResult type - got = %v, want %v", reflect.TypeOf(tt.wantResult), reflect.TypeOf(&successResponse{}))
				}

				if *res.Version != *wantRes.Version {
					t.Fatalf("Client.decodeBody() - result - got = %v, want %v", *res.Version, *wantRes.Version)
				}
			}

			if tt.wantErr != nil {
				res, ok := err.(*Form3APIError)
				if !ok {
					t.Fatalf("Client.decodeBody() - result type - got = %v, want %v", reflect.TypeOf(tt.args.result), reflect.TypeOf(&successResponse{}))
				}

				wantRes, ok := tt.wantErr.(Form3APIError)
				if !ok {
					t.Fatalf("Client.decodeBody() - wantResult type - got = %v, want %v", reflect.TypeOf(tt.wantResult), reflect.TypeOf(&successResponse{}))
				}

				if res.Message != wantRes.Message {
					t.Fatalf("Client.decodeBody() - result - got = %v, want %v", res.Message, wantRes.Message)
				}

				if res.StatusCode != wantRes.StatusCode {
					t.Fatalf("Client.decodeBody() - result - got = %v, want %v", res.StatusCode, wantRes.StatusCode)
				}
			}
		})
	}
}

func TestClient_Do(t *testing.T) {
	testCases := []struct {
		name         string
		method       string
		path         string
		body         interface{}
		handleFunc   http.HandlerFunc
		expectedBody string
		wantErr      *Form3APIError
	}{
		{
			name:   "Test Client.Do() with GET method and success response",
			method: http.MethodGet,
			path:   "/v1/accounts/10",
			handleFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"data": {"id": "10"}}`))
			},
			expectedBody: `{"data":{"id":"10"}}`,
			wantErr:      nil,
		},
		{
			name:   "Test Client.Do() with GET method and 404 response",
			method: http.MethodGet,
			path:   "/v1/accounts/10",
			handleFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"error_message": "resource not found"}`))
			},
			wantErr: &Form3APIError{StatusCode: 404, Message: "resource not found"},
		},
		{
			name:   "Test Client.Do() with POST method and data",
			method: http.MethodPost,
			path:   "/v1/accounts",
			body:   struct{ ID string }{ID: "10000"},
			handleFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)

				bodyStr, _ := io.ReadAll(r.Body)
				b := struct{ ID string }{}
				json.Unmarshal(bodyStr, &b)

				w.Write([]byte(fmt.Sprintf(`{"data": {"id": "%s"}}`, b.ID)))
			},
			wantErr:      nil,
			expectedBody: `{"data":{"id":"10000"}}`,
		},
		{
			name:   "Test Client.Do() with timeout",
			method: http.MethodGet,
			path:   "/v1/accounts/10",
			handleFunc: func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(2 * time.Second)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"data": {"id": "10"}}`))
			},
			wantErr: &Form3APIError{StatusCode: http.StatusGatewayTimeout, Message: http.StatusText(http.StatusGatewayTimeout)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup server
			server := httptest.NewServer(http.HandlerFunc(tc.handleFunc))
			defer server.Close()

			// Setup client
			serverUrl, _ := url.Parse(server.URL)
			client := &Client{BaseURL: serverUrl, client: &http.Client{Timeout: 1 * time.Second}}

			// Make request
			var result map[string]interface{}
			err := client.Do(context.Background(), tc.method, tc.path, tc.body, &result)
			if err != nil && tc.wantErr == nil || err == nil && tc.wantErr != nil {
				t.Fatalf("Client.Do() error = %v, wantErr %v", err, tc.wantErr)
			}

			if tc.wantErr != nil {
				form3Error, ok := err.(*Form3APIError)
				if !ok {
					t.Fatalf("Client.Do() - result type - got = %v, want %v", reflect.TypeOf(err), reflect.TypeOf(&Form3APIError{}))
				}

				if form3Error.StatusCode != tc.wantErr.StatusCode {
					t.Fatalf("Client.Do() - result - got = %v, want %v", form3Error.StatusCode, tc.wantErr.StatusCode)
				}

				if form3Error.Message != tc.wantErr.Message {
					t.Fatalf("Client.Do() - result - got = %v, want %v", form3Error.Message, tc.wantErr.Message)
				}
			} else {
				body, err := json.Marshal(result)
				if err != nil {
					t.Fatalf("Client.Do() Marshal error = %v", err)
				}

				if string(body) != tc.expectedBody {
					t.Fatalf("Client.Do() - result - got = %v, want %v", string(body), tc.expectedBody)
				}
			}
		})
	}
}
