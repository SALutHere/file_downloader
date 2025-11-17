package downloader

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type FileInfo struct {
	URL          string
	Size         int
	AcceptRanges bool
	ContentType  string
}

// GetRemoteFileInfo gets information about a file.
// Firstly, it tries to get info from HEAD-request.
// Secondly, if server does not support HEAD-requests, it tries to
// send GET Range: bytes=0-0
func GetRemoteFileInfo(url string) (*FileInfo, error) {
	resp, err := http.Head(url)
	if err == nil && resp.StatusCode < 400 {
		defer resp.Body.Close()
		return parseInfoFromHeaders(url, resp)
	}

	rangeReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("GET-request praparing error: %w", err)
	}
	rangeReq.Header.Set("Range", "bytes=0-0")

	rangeResp, err := http.DefaultClient.Do(rangeReq)
	if err != nil {
		return nil, fmt.Errorf("error executing range-request: %w", err)
	}
	defer rangeResp.Body.Close()

	return parseInfoFromRange(url, rangeResp)
}

// parseInfoFromHeaders extracts file information from the response of the HEAD request
func parseInfoFromHeaders(url string, resp *http.Response) (*FileInfo, error) {
	info := &FileInfo{
		URL:          url,
		ContentType:  resp.Header.Get("Content-Type"),
		AcceptRanges: resp.Header.Get("Accept-Ranges") == "bytes",
	}

	if cl := resp.Header.Get("Content-Length"); cl != "" {
		size, err := strconv.Atoi(cl)
		if err != nil {
			return nil, fmt.Errorf("Content-Length read error: %w", err)
		}
		info.Size = size
	}

	return info, nil
}

// parseInfoFromRange extracts file information from the response of the GET Range request
func parseInfoFromRange(url string, resp *http.Response) (*FileInfo, error) {
	cr := resp.Header.Get("Content-Range")
	if cr == "" {
		return nil, fmt.Errorf("the server does not support range queries")
	}

	parts := strings.Split(cr, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("incorrect Content-Range: %s", cr)
	}

	size, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("error extracting file size: %w", err)
	}

	info := &FileInfo{
		URL:          url,
		Size:         size,
		AcceptRanges: true,
		ContentType:  resp.Header.Get("Content-Type"),
	}

	return info, nil
}
