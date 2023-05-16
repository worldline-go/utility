package httpx

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/go-test/deep"
)

type TestRequest struct {
	ID string `json:"id"`
}

func (TestRequest) Method() string {
	return http.MethodPost
}

func (TestRequest) Path() string {
	return "/api/v1/test"
}

func (r TestRequest) BodyJSON() interface{} {
	return r
}

func (r TestRequest) Validate() error {
	if r.ID == "" {
		return fmt.Errorf("id is required")
	}

	return nil
}

func (r TestRequest) Header() http.Header {
	v := http.Header{}
	v.Set("X-Info", "test")

	return v
}

func TestClient_Do(t *testing.T) {
	retryCount := 0
	httpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// retry
		if retryCount > 0 {
			retryCount--
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "internal server error"}`))

			return
		}

		// check request method
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "invalid request method"}`))
			return
		}

		// check request path
		if r.URL.Path != "/api/v1/test" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "invalid request path"}`))
			return
		}

		// check request header
		if r.Header.Get("X-Info") != "test" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "invalid request header"}`))
			return
		}

		// get request body
		var m map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "invalid request body"}`))
			return
		}

		// check request body
		if m["id"] != "123" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "invalid id"}`))
			return
		}

		// write response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"request_id": "123+"}`))
	}))

	defer httpServer.Close()

	httpxClient, err := NewClient(
		WithBaseURL(httpServer.URL),
	)
	if err != nil {
		t.Errorf("NewClient() error = %v", err)
		return
	}

	type fields struct {
		HttpClient *http.Client
		BaseURL    *url.URL
	}
	type args struct {
		ctx  context.Context
		req  Request
		resp interface{}
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		want       interface{}
		wantErr    bool
		retryCount int
		short      bool
		long       bool
	}{
		{
			name: "Do",
			fields: fields{
				HttpClient: httpxClient.HttpClient,
				BaseURL:    httpxClient.BaseURL,
			},
			args: args{
				ctx: context.Background(),
				req: TestRequest{
					ID: "123",
				},
				resp: new(map[string]interface{}),
			},
			want: map[string]interface{}{
				"request_id": "123+",
			},
			wantErr: false,
		},
		{
			name: "Do",
			fields: fields{
				HttpClient: httpxClient.HttpClient,
				BaseURL:    httpxClient.BaseURL,
			},
			args: args{
				ctx: context.Background(),
				req: TestRequest{
					ID: "123",
				},
				resp: new(map[string]interface{}),
			},
			want: map[string]interface{}{
				"request_id": "123+",
			},
			wantErr:    true,
			retryCount: 5,
			long:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.long && testing.Short() {
				t.Skip("skipping test in short mode.")
			}

			c := &Client{
				HttpClient: tt.fields.HttpClient,
				BaseURL:    tt.fields.BaseURL,
			}

			retryCount = tt.retryCount
			if err := c.Do(tt.args.ctx, tt.args.req, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("Client.Do() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if tt.wantErr {
				return
			}

			// tt.args.resp is a pointer
			resp := reflect.ValueOf(tt.args.resp).Elem().Interface()
			if diff := deep.Equal(resp, tt.want); diff != nil {
				t.Errorf("Client.Do() resp diff = %v", diff)
			}
		})
	}
}
