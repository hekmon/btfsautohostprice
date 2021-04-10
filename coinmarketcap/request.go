package coinmarketcap

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	baseURLstr   = "https://pro-api.coinmarketcap.com"
	apiKeyHeader = "X-CMC_PRO_API_KEY"
)

var (
	baseURL *url.URL
)

func init() {
	var err error
	baseURL, err = url.Parse(baseURLstr)
	if err != nil {
		panic(err)
	}
}

func (c *Controller) request(ctx context.Context, method, endpoint string, queryParam url.Values, body io.Reader, result interface{}) (creditCount int, err error) {
	// Build URL
	APIURL := *baseURL
	APIURL.Path += endpoint
	if queryParam != nil {
		APIURL.RawQuery = queryParam.Encode()
	}
	// Build Request
	req, err := http.NewRequestWithContext(ctx, method, APIURL.String(), body)
	if err != nil {
		err = fmt.Errorf("can not build HTTP request: %w", err)
		return
	}
	req.Header.Set("Accepts", "application/json")
	req.Header.Add(apiKeyHeader, c.apiKey)
	// Execute request
	resp, err := c.client.Do(req)
	if err != nil {
		err = fmt.Errorf("error while executing HTTP request: %w", err)
		return
	}
	defer resp.Body.Close()
	// Decode response
	payload := responsePayload{
		Data: result,
	}
	switch resp.StatusCode {
	case http.StatusBadRequest:
		fallthrough
	case http.StatusUnauthorized:
		fallthrough
	case http.StatusForbidden:
		fallthrough
	case http.StatusTooManyRequests:
		fallthrough
	case http.StatusInternalServerError:
		fallthrough
	case http.StatusOK:
		if err = json.NewDecoder(resp.Body).Decode(&payload); err != nil {
			err = fmt.Errorf("failed to decode HTTP response as JSON: %w", err)
			return
		}
	default:
		err = fmt.Errorf("unexpected HTTP code: %s", resp.Status)
		return
	}
	// Handle status
	creditCount = payload.Status.CreditCount
	if payload.Status.ErrorCode != 0 {
		err = fmt.Errorf("api error %d: %s", payload.Status.ErrorCode, payload.Status.ErrorMsg)
		return
	}
	return
}

type responsePayload struct {
	Data   interface{}    `json:"data"`
	Status responseStatus `json:"status"`
}

type responseStatus struct {
	Time        time.Time `json:"timestamp"`
	ErrorCode   int       `json:"error_code"`
	ErrorMsg    string    `json:"error_message"`
	Elapsed     int       `json:"elapsed"`
	CreditCount int       `json:"credit_count"`
}
