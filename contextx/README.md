# contextx

Context value helper functions.
Differences between context package this holds value in map with mutex lock.

```sh
go get github.com/worldline-go/utility/contextx
```

## Usage

```go
// set value first initialize context value map if not exist
ctx := contextx.WithValue(context.Background(), "secret", "xxx")

// map like access
if v, ok := contextx.Value[string](ctx, "secret"); ok {
    // v is string
}
```
