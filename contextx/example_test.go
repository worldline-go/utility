package contextx_test

import (
	"context"
	"fmt"

	"github.com/worldline-go/utility/contextx"
)

func Example() {
	ctx := contextx.WithValue(context.Background(), "secret", "xxx")

	if v, ok := contextx.Value[string](ctx, "secret"); ok {
		fmt.Println(v)
	}
	// Output:
	// xxx
}
