package abfallkalender

import (
	"fmt"
	"net/http"
	"time"
)

// httpClient is shared by all requests and carries a timeout so a stalled
// insert-it.de endpoint cannot hang a caller indefinitely.
var httpClient = &http.Client{Timeout: 30 * time.Second}

// httpGet performs a GET request and returns the response only on a 2xx
// status. The caller owns resp.Body and must close it; on a non-2xx status
// httpGet closes the body itself before returning the error.
func httpGet(url string) (*http.Response, error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		resp.Body.Close()
		return nil, fmt.Errorf("insert-it.de returned HTTP %s for %s", resp.Status, url)
	}
	return resp, nil
}
