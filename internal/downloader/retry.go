package downloader

import (
	"fmt"
	"io"
	"time"

	"github.com/SALutHere/file_downloader/internal/progress"
)

const (
	maxRetries     = 3
	baseRetryDelay = time.Second
)

// DownloadChunkWithRetry downloads a single range of bytes and writes it to a file.
// It is like DownloadChunk, but it retries till the maxRetries and uses exponential delay
// between attempts.
func DownloadChunkWithRetry(
	url string,
	ch Chunk,
	file io.WriterAt,
	state *State,
	progress chan<- progress.ChunkProgress,
) error {
	var err error

	for att := 1; att <= maxRetries; att++ {
		err = DownloadChunk(url, ch, file, state, progress)
		if err == nil {
			return nil
		}

		if att < maxRetries {
			delay := baseRetryDelay << (att - 1) // exponential delay
			time.Sleep(delay)
		}
	}

	return fmt.Errorf(
		"chunk #%d failed to download after %d attempts: %w",
		ch.Index,
		maxRetries,
		err,
	)
}
