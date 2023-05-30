package httpx

import (
	"fmt"
	"net/http"
)

var (
	ErrValidating      = fmt.Errorf("failed to validate request")
	ErrMarshal         = fmt.Errorf("failed to marshal request body")
	ErrRequest         = fmt.Errorf("failed to do request")
	ErrResponseFuncNil = fmt.Errorf("response function is nil")
)

func UnexpectedResponseError(resp *http.Response) error {
	partialBody := limitedResponse(resp)

	return fmt.Errorf("unexpected response (%d): %s", resp.StatusCode, string(partialBody))
}
