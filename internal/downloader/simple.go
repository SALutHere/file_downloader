package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// DownloadSimple performs simple download of a file without download resuming and threads
func DownloadSimple(url, outputPath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("request execution error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("the server returned the status: %s", resp.Status)
	}

	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("file creation error: %w", err)
	}
	defer out.Close()

	if _, err = io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("file writing error: %w", err)
	}

	return nil
}
