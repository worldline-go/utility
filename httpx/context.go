package httpx

import (
	"context"
	"net/http"
)

type ctxValueType string

const (
	ctxValue ctxValueType = "requestValue"
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
	values map[rValueType]interface{}
}

func (o *optionContext) SetValue(key rValueType, value interface{}) {
	o.values[key] = value
}

type OptionContext func(*optionContext)

// WithHeader sets the header to be sent.
func CtxWithHeader(key, value string) OptionContext {
	return func(o *optionContext) {
		if v, ok := o.values[rValueHeader]; ok {
			v.(http.Header).Add(key, value)

			return
		}

		header := http.Header{}
		header.Add(key, value)

		o.SetValue(rValueHeader, header)
	}
}

// WithRetry sets the retry to be sent.
//
// Just work with our RetryPolicy.
func CtxWithRetry(retry Retry) OptionContext {
	return func(o *optionContext) {
		o.SetValue(rValueRetry, retry)
	}
}

// ---

var unUsedValue = make(map[rValueType]interface{})

func requestCtxGet[T any](ctx context.Context, key rValueType) (T, bool) {
	if ctx == nil {
		return unUsedValue[key].(T), true
	}

	o, ok := ctx.Value(ctxValue).(map[rValueType]interface{})
	if !ok {
		o = unUsedValue
	}

	ret, ok := o[key].(T)

	return ret, ok
}

// ---

// RequestCtx adds the request options to the context.
func RequestCtx(ctx context.Context, opts ...OptionContext) context.Context {
	o := optionContext{
		values: make(map[rValueType]interface{}),
	}

	for _, opt := range opts {
		opt(&o)
	}

	return context.WithValue(ctx, ctxValue, o.values)
}
