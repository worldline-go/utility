package httpx

import (
	"context"
	"net/http"

	"github.com/worldline-go/utility/contextx"
)

type rValueType string

type (
	rValueHeaderType = http.Header
	rValueRetryType  = Retry
)

const (
	rValueHeader rValueType = "header"
	rValueRetry  rValueType = "retry"
)

type optionContext struct {
	ctx context.Context
}

type OptionContext func(*optionContext)

// WithHeader sets the header to be sent.
func CtxWithHeader(key, value string) OptionContext {
	return func(o *optionContext) {
		if v, ok := contextx.Value[rValueHeaderType](o.ctx, rValueHeader); ok {
			v.Set(key, value)

			return
		}

		header := http.Header{}
		header.Set(key, value)

		o.ctx = contextx.WithValue(o.ctx, rValueHeader, header)
	}
}

// WithRetry sets the retry to be sent.
//
// Just work with our RetryPolicy.
func CtxWithRetry(retry Retry) OptionContext {
	return func(o *optionContext) {
		o.ctx = contextx.WithValue(o.ctx, rValueRetry, retry)
	}
}

// RequestCtx adds the request options to the context.
func RequestCtx(ctx context.Context, opts ...OptionContext) context.Context {
	o := optionContext{
		ctx: ctx,
	}

	for _, opt := range opts {
		opt(&o)
	}

	return o.ctx
}
