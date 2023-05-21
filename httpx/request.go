package httpx

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type request struct {
	method string
	path   string
	header http.Header
	body   io.Reader
}

func (r request) Method() string {
	return r.method
}

func (r request) Path() string {
	return r.path
}

func (r request) Header() http.Header {
	return r.header
}

func (r request) Body() io.Reader {
	return r.body
}

// Request sends an HTTP request and calls the response function.
func (c *Client) Request(ctx context.Context, method, path string, body io.Reader, header http.Header, fn func(*http.Response) error) error {
	return c.DoWithFunc(ctx, request{
		method: method,
		path:   path,
		header: header,
		body:   body,
	}, fn)
}

// RequestWithURL sends an HTTP request and calls the response function.
func (c *Client) RequestWithURL(ctx context.Context, method, URL string, body io.Reader, header http.Header, fn func(*http.Response) error) error {
	baseUrl, err := url.Parse(URL)
	if err != nil {
		return fmt.Errorf("parse base url: %w", err)
	}

	return c.DoWithFunc(
		RequestCtx(ctx, CtxWithBaseURL(baseUrl)),
		request{
			method: method,
			header: header,
			body:   body,
		},
		fn,
	)
}
