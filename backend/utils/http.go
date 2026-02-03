package utils

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPClient is a custom HTTP client with retry logic and timeouts
var HTTPClient = &http.Client{
	Timeout: 3 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
	},
}

// FetchWithRetry performs HTTP GET with automatic retry on failure
func FetchWithRetry(url string, maxRetries int) ([]byte, error) {
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		resp, err := HTTPClient.Get(url)

		// Success case
		if err == nil && resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()
			body, readErr := io.ReadAll(resp.Body)
			if readErr != nil {
				lastErr = fmt.Errorf("failed to read response: %w", readErr)
				continue
			}
			return body, nil
		}

		// Handle rate limiting (429)
		if resp != nil && resp.StatusCode == http.StatusTooManyRequests {
			resp.Body.Close()
			waitTime := time.Duration(attempt+1) * time.Second
			time.Sleep(waitTime)
			lastErr = fmt.Errorf("rate limited, retrying after %v", waitTime)
			continue
		}

		// Handle other errors
		if resp != nil {
			resp.Body.Close()
			lastErr = fmt.Errorf("HTTP %d: %s", resp.StatusCode, http.StatusText(resp.StatusCode))
		} else {
			lastErr = err
		}

		// Exponential backoff
		if attempt < maxRetries-1 {
			time.Sleep(500 * time.Millisecond * time.Duration(attempt+1))
		}
	}

	return nil, fmt.Errorf("max retries (%d) exceeded: %w", maxRetries, lastErr)
}

// FetchWithHeaders performs HTTP GET with custom headers and retry logic
func FetchWithHeaders(url string, headers map[string]string, maxRetries int) ([]byte, error) {
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		// Add common headers
		req.Header.Set("Accept", "application/json")
		req.Header.Set("User-Agent", "AlphaAgent/1.0")

		// Add custom headers
		for key, value := range headers {
			req.Header.Set(key, value)
		}

		resp, err := HTTPClient.Do(req)

		// Success case
		if err == nil && resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()
			body, readErr := io.ReadAll(resp.Body)
			if readErr != nil {
				lastErr = fmt.Errorf("failed to read response: %w", readErr)
				continue
			}
			return body, nil
		}

		// Handle rate limiting (429)
		if resp != nil && resp.StatusCode == http.StatusTooManyRequests {
			resp.Body.Close()
			waitTime := time.Duration(attempt+1) * time.Second
			time.Sleep(waitTime)
			lastErr = fmt.Errorf("rate limited, retrying after %v", waitTime)
			continue
		}

		// Handle other errors
		if resp != nil {
			resp.Body.Close()
			lastErr = fmt.Errorf("HTTP %d: %s", resp.StatusCode, http.StatusText(resp.StatusCode))
		} else {
			lastErr = err
		}

		// Exponential backoff
		if attempt < maxRetries-1 {
			time.Sleep(500 * time.Millisecond * time.Duration(attempt+1))
		}
	}

	return nil, fmt.Errorf("max retries (%d) exceeded: %w", maxRetries, lastErr)
}

// FetchJSON is a convenience wrapper for JSON endpoints
func FetchJSON(url string) ([]byte, error) {
	return FetchWithRetry(url, 3)
}

// FetchJSONWithHeaders is a convenience wrapper for JSON endpoints with headers
func FetchJSONWithHeaders(url string, headers map[string]string) ([]byte, error) {
	return FetchWithHeaders(url, headers, 3)
}
