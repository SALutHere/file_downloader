package downloader

import (
	"fmt"
	"io"
	"net/http"
)

const bufferSize = 64 * 1024 // 64 KB

type ChunkProgress struct {
	Index      int
	Downloaded int64
	Total      int64
}

// DownloadChunk downloads (and resumes) a chunk.
// Uses incremental writes and updates chunk progress.
func DownloadChunk(
	url string,
	ch Chunk,
	file io.WriterAt,
	state *State,
	progress chan<- ChunkProgress,
) error {
	st := &state.Chunks[ch.Index]

	if st.Done {
		return nil
	}

	state.mu.Lock()
	offset := st.Start + st.Downloaded
	state.mu.Unlock()

	totalSize := st.End - st.Start + 1

	if st.Downloaded >= totalSize {
		st.Done = true
		return nil
	}

	rangeHeader := fmt.Sprintf("bytes=%d-%d", offset, ch.End)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Range", rangeHeader)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("chunk #%d request execution error: %w", ch.Index, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusPartialContent {
		return fmt.Errorf("the server returned status: %s, expected: %d", resp.Status, http.StatusPartialContent)
	}

	buf := make([]byte, bufferSize)

	for {
		n, readErr := resp.Body.Read(buf)
		if n > 0 {
			if _, err = file.WriteAt(buf[:n], offset); err != nil {
				return fmt.Errorf("error writing chunk #%d: %w", ch.Index, err)
			}

			state.mu.Lock()
			st.Downloaded += int64(n)
			offset += int64(n)
			state.mu.Unlock()

			if err = state.Save(); err != nil {
				return fmt.Errorf("error saving state: %w", err)
			}

			if progress != nil {
				progress <- ChunkProgress{
					Index:      ch.Index,
					Downloaded: st.Downloaded,
					Total:      totalSize,
				}
			}

			if st.Downloaded >= totalSize {
				state.mu.Lock()
				st.Done = true
				state.mu.Unlock()

				if err = state.Save(); err != nil {
					return fmt.Errorf("final state save error: %w", err)
				}
				return nil
			}
		}

		if readErr == io.EOF {
			return nil
		}
		if readErr != nil {
			return readErr
		}
	}
}
