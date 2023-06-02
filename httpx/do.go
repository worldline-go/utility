package httpx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/worldline-go/utility/contextx"
)

type Request interface {
	// Method returns the HTTP method.
	Method() string
	// Path returns the path to be sent.
	Path() string
}

type Header interface {
	// Header returns the header to be sent.
	Header() http.Header
}

type RequestValidator interface {
	// Validate returns error if request is invalid.
	Validate() error
}

type QueryStringGenerator interface {
	// ToQuery returns the query string to be sent.
	ToQuery() url.Values
}

type BodyJSON interface {
	// BodyJSON can return any type that can be marshaled to JSON.
	// Automatically sets Content-Type to application/json.
	BodyJSON() interface{}
}

type Body interface {
	// Body returns the body to be sent.
	Body() io.Reader
}

// DoWithFunc sends an HTTP request and calls the response function.
//
// Request additional implements RequestValidator, QueryStringGenerator, Header, Body, BodyJSON
func (c *Client) DoWithFunc(ctx context.Context, req Request, fn func(*http.Response) error) error {
	baseURL := c.BaseURL
	if rURL, ok := contextx.Value[rValueURLType](ctx, rValueURL); ok {
		baseURL = rURL
	}

	if baseURL == nil {
		return fmt.Errorf("base url is required")
	}

	if v, ok := req.(RequestValidator); ok {
		if err := v.Validate(); err != nil {
			return fmt.Errorf("%w: %v", ErrValidating, err)
		}
	}

	u := &url.URL{
		Path: req.Path(),
	}

	if g, ok := req.(QueryStringGenerator); ok {
		u.RawQuery = g.ToQuery().Encode()
	}

	var (
		header http.Header
		body   io.Reader
	)

	// set default header
	if h, ok := req.(Header); ok {
		header = h.Header()
	}

	if header == nil {
		header = make(http.Header)
	}

	if b, ok := req.(Body); ok {
		body = b.Body()
	} else if jb, ok := req.(BodyJSON); ok {
		bodyGet := jb.BodyJSON()

		bodyData, err := json.Marshal(bodyGet)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrMarshal, err)
		}

		body = bytes.NewReader(bodyData)

		header.Set("Content-Type", "application/json")
	}

	// add context values
	if rHeader, ok := contextx.Value[rValueHeaderType](ctx, rValueHeader); ok {
		for k := range rHeader {
			header.Set(k, rHeader.Get(k))
		}
	}

	uSend := baseURL.ResolveReference(u)

	httpReq, _ := http.NewRequestWithContext(ctx, req.Method(), uSend.String(), body)
	httpReq.Header = header

	httpResp, err := c.HttpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrRequest, err)
	}

	defer func() {
		_, _ = io.Copy(io.Discard, httpResp.Body)
		_ = httpResp.Body.Close()
	}()

	if fn == nil {
		return ErrResponseFuncNil
	}

	return fn(httpResp)
}

// Do sends an HTTP request and json unmarshals the response body to data.
//
// Do work same as DoWithFunc with defaultResponseFunc.
func (c *Client) Do(ctx context.Context, req Request, resp interface{}) error {
	fn := func(r *http.Response) error { return defaultResponseFunc(r, resp) }

	return c.DoWithFunc(ctx, req, fn)
}

func defaultResponseFunc(resp *http.Response, data interface{}) error {
	if err := UnexpectedResponse(resp); err != nil {
		return err
	}

	// 204s, for example
	if resp.ContentLength == 0 {
		return nil
	}

	if data == nil {
		return nil
	}

	if err := json.NewDecoder(resp.Body).Decode(data); err != nil {
		return fmt.Errorf("decode response body: %w", err)
	}

	return nil
}
