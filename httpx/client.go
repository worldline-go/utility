package httpx

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-retryablehttp"
)

var (
	defaultMaxConnections = 100
	// Default retry configuration
	defaultRetryWaitMin = 1 * time.Second
	defaultRetryWaitMax = 30 * time.Second
	defaultRetryMax     = 4
)

type Client struct {
	HttpClient *http.Client
	BaseURL    *url.URL
}

func NewClient(opts ...Option) (*Client, error) {
	o := optionClient{
		PooledClient:   true,
		MaxConnections: defaultMaxConnections,
		RetryWaitMin:   defaultRetryWaitMin,
		RetryWaitMax:   defaultRetryWaitMax,
		RetryMax:       defaultRetryMax,
		RetryPolicy:    RetryPolicy,
		Backoff:        retryablehttp.DefaultBackoff,
		Logger:         log.New(os.Stderr, "", log.LstdFlags),
	}

	for _, opt := range opts {
		opt(&o)
	}

	if o.BaseURL == "" {
		return nil, fmt.Errorf("base url is required")
	}

	baseUrl, err := url.Parse(o.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base url: %w", err)
	}

	// create client
	client := o.HttpClient
	if client == nil {
		if o.PooledClient {
			client = cleanhttp.DefaultPooledClient()
		} else {
			client = cleanhttp.DefaultClient()
		}
	}

	if o.InsecureSkipVerify {
		//nolint:forcetypeassert // clear
		tlsClientConfig := client.Transport.(*http.Transport).TLSClientConfig
		if tlsClientConfig == nil {
			tlsClientConfig = &tls.Config{
				//nolint:gosec // user defined
				InsecureSkipVerify: true,
			}
		} else {
			tlsClientConfig.InsecureSkipVerify = true
		}

		//nolint:forcetypeassert // clear
		client.Transport.(*http.Transport).TLSClientConfig = tlsClientConfig
	}

	if o.TransportWrapper != nil {
		transport, err := o.TransportWrapper(o.Ctx, client.Transport)
		if err != nil {
			return nil, fmt.Errorf("failed to wrap transport: %w", err)
		}

		client.Transport = transport
	}

	// disable retry
	if !o.DisableRetry {
		// create retry client
		retryClient := retryablehttp.Client{
			HTTPClient:   client,
			Logger:       o.Logger,
			RetryWaitMin: o.RetryWaitMin,
			RetryWaitMax: o.RetryWaitMax,
			RetryMax:     o.RetryMax,
			CheckRetry:   o.RetryPolicy,
			Backoff:      o.Backoff,
		}

		client = retryClient.StandardClient()
	}

	return &Client{
		HttpClient: client,
		BaseURL:    baseUrl,
	}, nil
}
