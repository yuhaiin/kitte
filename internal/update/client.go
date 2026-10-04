package update

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

const maxDownloadSize = 64 << 20

type httpClient struct {
	client *http.Client
}

func newHTTPClient() *httpClient {
	return &httpClient{client: &http.Client{Timeout: 2 * time.Minute}}
}

func (c *httpClient) get(ctx context.Context, url string) ([]byte, error) {
	var lastErr error

	for attempt := 1; attempt <= 3; attempt++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("User-Agent", "yuhaiin-kitte-updater/1")

		resp, err := c.client.Do(req)
		if err == nil {
			body, readErr := io.ReadAll(io.LimitReader(resp.Body, maxDownloadSize+1))
			resp.Body.Close()

			if readErr == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
				if len(body) > maxDownloadSize {
					return nil, fmt.Errorf("download %s exceeds %d bytes", url, maxDownloadSize)
				}
				return body, nil
			}

			if readErr != nil {
				lastErr = fmt.Errorf("read %s: %w", url, readErr)
			} else {
				lastErr = fmt.Errorf("download %s: HTTP %s", url, resp.Status)
				if resp.StatusCode < 500 && resp.StatusCode != http.StatusTooManyRequests {
					return nil, lastErr
				}
			}
		} else {
			lastErr = fmt.Errorf("download %s: %w", url, err)
		}

		if attempt == 3 {
			break
		}

		timer := time.NewTimer(time.Duration(attempt) * time.Second)
		select {
		case <-ctx.Done():
			timer.Stop()
			return nil, ctx.Err()
		case <-timer.C:
		}
	}

	return nil, lastErr
}
