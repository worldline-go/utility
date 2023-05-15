package httpx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Request interface {
	Method() string
	Path() string
}

type RequestValidator interface {
	Validate() error
}

type QueryStringGenerator interface {
	ToQuery() url.Values
}

type JSONBodyFlag interface {
	GetJSONBody() interface{}
}

func (c *Client) DoWithFunc(ctx context.Context, req Request, fn func(*http.Response) error) error {
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
		header = make(http.Header)
		body   io.Reader
	)

	if jb, ok := req.(JSONBodyFlag); ok {
		bodyGet := jb.GetJSONBody()

		bodyData, err := json.Marshal(bodyGet)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrMarshal, err)
		}

		body = bytes.NewReader(bodyData)

		header.Set("Content-Type", "application/json")
	}

	if c.BaseURL == nil {
		return fmt.Errorf("base url is required")
	}

	uSend := c.BaseURL.ResolveReference(u)

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

func (c *Client) Do(ctx context.Context, req Request, resp interface{}) error {
	fn := func(r *http.Response) error { return defaultResponseFunc(r, resp) }

	return c.DoWithFunc(ctx, req, fn)
}

func defaultResponseFunc(resp *http.Response, data interface{}) error {
	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		return UnexpectedResponseError(resp)
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
