package downloader

import (
	"fmt"
	"io"
	"net/http"
)

// DownloadChunk downloads a single range of bytes and writes it to a file.
func DownloadChunk(url string, ch Chunk, file io.WriterAt) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	rangeHeader := fmt.Sprintf("bytes=%d-%d", ch.Start, ch.End)
	req.Header.Set("Range", rangeHeader)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("chunk #%d request execution error: %w", ch.Index, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusPartialContent {
		return fmt.Errorf("the server returned status: %s, expected: %d", resp.Status, http.StatusPartialContent)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	if _, err = file.WriteAt(data, ch.Start); err != nil {
		return fmt.Errorf("error writing chunk #%d to file: %w", ch.Index, err)
	}

	return nil
}
