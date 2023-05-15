package httpx_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/worldline-go/utility/httpx"
)

type Client struct {
	httpx *httpx.Client
}

func (c *Client) CreateX(ctx context.Context, r CreateXRequest) (*CreateXResponse, error) {
	var v CreateXResponse
	if err := c.httpx.DoWithFunc(ctx, r, func(r *http.Response) error {
		if r.StatusCode != http.StatusOK {
			return httpx.UnexpectedResponseError(r)
		}

		if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &v, nil
}

// ----

type CreateXRequest struct {
	ID string `json:"id"`
}

func (CreateXRequest) Method() string {
	return http.MethodPost
}

func (CreateXRequest) Path() string {
	return "/api/v1/x"
}

func (r CreateXRequest) GetJSONBody() interface{} {
	return r
}

type CreateXResponse struct {
	RequestID string `json:"request_id"`
}

// ---

func Example() {
	httpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check request method
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "invalid request method"}`))
			return
		}

		// check request path
		if r.URL.Path != "/api/v1/x" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "invalid request path"}`))
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

	httpxClient, err := httpx.NewClient(
		httpx.WithBaseURL(httpServer.URL),
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	c := &Client{
		httpx: httpxClient,
	}

	v, err := c.CreateX(context.Background(), CreateXRequest{
		ID: "123",
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(v.RequestID)
	// Output:
	// 123+
}
