package contextx

import (
	"context"
	"reflect"
	"testing"
)

func TestValue(t *testing.T) {
	tests := []struct {
		name   string
		ctx    context.Context
		setFn  func(context.Context) context.Context
		getFn  func(context.Context) (any, bool)
		want   any
		wantOk bool
	}{
		{
			name: "test string",
			setFn: func(ctx context.Context) context.Context {
				return WithValue(ctx, "test", "xxx")
			},
			getFn: func(ctx context.Context) (any, bool) {
				return Value[string](ctx, "test")
			},
			want:   "xxx",
			wantOk: true,
		},
		{
			name: "test int",
			setFn: func(ctx context.Context) context.Context {
				return WithValue(ctx, "test", 1234)
			},
			getFn: func(ctx context.Context) (any, bool) {
				return Value[int](ctx, "test")
			},
			want:   1234,
			wantOk: true,
		},
		{
			name: "test diff",
			setFn: func(ctx context.Context) context.Context {
				return WithValue(ctx, "test", 1234)
			},
			getFn: func(ctx context.Context) (any, bool) {
				return Value[string](ctx, "test")
			},
			want:   "",
			wantOk: false,
		},
		{
			name: "test unset",
			getFn: func(ctx context.Context) (any, bool) {
				return Value[string](ctx, "test")
			},
			want:   "",
			wantOk: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.ctx
			if ctx == nil {
				ctx = context.Background()
			}

			if tt.setFn != nil {
				ctx = tt.setFn(ctx)
			}

			want, wantOk := tt.getFn(ctx)
			if (wantOk != false) != tt.wantOk {
				t.Errorf("GetValue() gotOk = %v, want %v", wantOk, tt.wantOk)
			}

			if !reflect.DeepEqual(want, tt.want) {
				t.Errorf("GetValue() got = %v, want %v", want, tt.want)
			}
		})
	}
}
