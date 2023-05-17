package contextx

import (
	"context"
)

type ctxInternalKey string

const ctxInternal ctxInternalKey = "internal"

var unUsedValue = newContextValue()

// Value gets value from context's map.
//
// Returns same as map's value and ok.
func Value[T any](ctx context.Context, key any) (T, bool) {
	var o *contextValue
	if ctx == nil {
		o = unUsedValue
	} else {
		var ok bool
		if o, ok = getValue(ctx); !ok {
			o = unUsedValue
		}
	}

	o.m.RLock()
	defer o.m.RUnlock()

	ret, ok := o.v[key].(T)

	return ret, ok
}

// WithValue sets value to context's map.
//
// If context is nil, it use context.Background().
// If context's map is nil, it will init context's map and add value to ctx.
func WithValue(ctx context.Context, key any, value any) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	o, ok := getValue(ctx)
	if !ok {
		return WithValue(Init(ctx), key, value)
	}

	o.m.Lock()
	defer o.m.Unlock()

	o.v[key] = value

	return ctx
}

// Init context with internal value.
//
// If context is nil, it use context.Background().
func Init(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	return context.WithValue(ctx, ctxInternal, newContextValue())
}

func getValue(ctx context.Context) (*contextValue, bool) {
	if ctx == nil {
		return nil, false
	}

	v, ok := ctx.Value(ctxInternal).(*contextValue)
	return v, ok
}
