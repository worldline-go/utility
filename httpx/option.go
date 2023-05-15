package httpx

import (
	"context"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/rs/zerolog"
	"github.com/worldline-go/logz"
)

type optionClient struct {
	// HttpClient is the http client.
	HttpClient *http.Client
	// PooledClient is generate pooled client if no http client provided.
	// Default is true.
	PooledClient bool
	// TransportWrapper is a function that wraps the default transport.
	TransportWrapper func(context.Context, http.RoundTripper) (http.RoundTripper, error)
	// Ctx for TransportWrapper.
	Ctx context.Context
	// MaxConnections is the maximum number of idle (keep-alive) connections.
	MaxConnections int
	// Logger is the customer logger instance of retryablehttp. Can be either Logger or LeveledLogger
	Logger interface{}
	// InsecureSkipVerify is the flag to skip TLS verification.
	InsecureSkipVerify bool

	// BaseURL is the base URL of the service.
	BaseURL string

	// DisableRetry is the flag to disable retry.
	DisableRetry bool
	// RetryWaitMin is the minimum wait time.
	// Default is 100ms.
	RetryWaitMin time.Duration
	// RetryWaitMax is the maximum wait time.
	RetryWaitMax time.Duration
	// RetryMax is the maximum number of retry.
	RetryMax int
	// RetryPolicy is the retry policy.
	RetryPolicy retryablehttp.CheckRetry
	// Backoff is the backoff policy.
	Backoff retryablehttp.Backoff
}

// Option is a function that configures the client.
type Option func(*optionClient)

// WithHttpClient configures the client to use the provided http client.
func WithHttpClient(httpClient *http.Client) Option {
	return func(o *optionClient) {
		o.HttpClient = httpClient
	}
}

func WithPooledClient(pooledClient bool) Option {
	return func(o *optionClient) {
		o.PooledClient = pooledClient
	}
}

// WithTransportWrapper configures the client to wrap the default transport.
func WithTransportWrapper(f func(context.Context, http.RoundTripper) (http.RoundTripper, error)) Option {
	return func(o *optionClient) {
		o.TransportWrapper = f
	}
}

// WithCtx configures the client to use the provided context.
func WithCtx(ctx context.Context) Option {
	return func(o *optionClient) {
		o.Ctx = ctx
	}
}

// WithMaxConnections configures the client to use the provided maximum number of idle connections.
func WithMaxConnections(maxConnections int) Option {
	return func(o *optionClient) {
		o.MaxConnections = maxConnections
	}
}

// WithLogger configures the client to use the provided logger.
func WithLogger(logger interface{}) Option {
	return func(o *optionClient) {
		o.Logger = logger
	}
}

// WithZerologLogger configures the client to use the provided logger.
func WithZerologLogger(logger zerolog.Logger) Option {
	return func(o *optionClient) {
		o.Logger = logz.AdapterKV{Log: logger}
	}
}

// WithInsecureSkipVerify configures the client to skip TLS verification.
func WithInsecureSkipVerify(insecureSkipVerify bool) Option {
	return func(o *optionClient) {
		o.InsecureSkipVerify = insecureSkipVerify
	}
}

// WithBaseURL configures the client to use the provided base URL.
func WithBaseURL(baseURL string) Option {
	return func(o *optionClient) {
		o.BaseURL = baseURL
	}
}

// WithDisableRetry configures the client to disable retry.
func WithDisableRetry(disableRetry bool) Option {
	return func(options *optionClient) {
		options.DisableRetry = disableRetry
	}
}

// WithRetryWaitMin configures the client to use the provided minimum wait time.
func WithRetryWaitMin(retryWaitMin time.Duration) Option {
	return func(options *optionClient) {
		options.RetryWaitMin = retryWaitMin
	}
}

// WithRetryWaitMax configures the client to use the provided maximum wait time.
func WithRetryWaitMax(retryWaitMax time.Duration) Option {
	return func(options *optionClient) {
		options.RetryWaitMax = retryWaitMax
	}
}

// WithRetryMax configures the client to use the provided maximum number of retry.
func WithRetryMax(retryMax int) Option {
	return func(options *optionClient) {
		options.RetryMax = retryMax
	}
}

// WithRetryPolicy configures the client to use the provided retry policy.
func WithRetryPolicy(retryPolicy retryablehttp.CheckRetry) Option {
	return func(options *optionClient) {
		options.RetryPolicy = retryPolicy
	}
}

// WithBackoff configures the client to use the provided backoff.
func WithBackoff(backoff retryablehttp.Backoff) Option {
	return func(options *optionClient) {
		options.Backoff = backoff
	}
}
