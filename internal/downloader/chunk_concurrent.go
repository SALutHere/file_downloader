package downloader

import (
	"fmt"
	"io"
	"sync"
)

// DownloadChunksConcurrently downloads chunks using goroutines
func DownloadChunksConcurrently(url string, chunks []Chunk, file io.WriterAt, state *State) error {
	var wg sync.WaitGroup

	errCh := make(chan error, len(chunks))

	for _, ch := range chunks {
		wg.Add(1)

		go func(ch Chunk) {
			defer wg.Done()

			if err := DownloadChunkWithRetry(url, ch, file); err != nil {
				errCh <- fmt.Errorf("error in downloading chunk #%d: %w", ch.Index, err)
				return
			}

			state.Chunks[ch.Index].Done = true

			if err := state.Save(); err != nil {
				errCh <- fmt.Errorf("error in saving state: %w", err)
				return
			}
		}(ch)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		return err
	}

	return nil
}
