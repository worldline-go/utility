package httpx

import (
	"fmt"
	"io"
	"net/http"
)

var (
	ErrValidating      = fmt.Errorf("failed to validate request")
	ErrMarshal         = fmt.Errorf("failed to marshal request body")
	ErrRequest         = fmt.Errorf("failed to do request")
	ErrResponseFuncNil = fmt.Errorf("response function is nil")
)

func UnexpectedResponseError(resp *http.Response) error {
	const truncateBodySize = 256

	limiter := io.LimitedReader{
		R: resp.Body,
		N: truncateBodySize,
	}

	partialBody, err := io.ReadAll(&limiter)
	if err != nil {
		return fmt.Errorf("read response body error: %w", err)
	}

	return fmt.Errorf("unexpected response (%d): %s", resp.StatusCode, string(partialBody))
}
