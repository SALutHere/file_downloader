package downloader

// Chunk represents a range of bytes to be downloaded.
type Chunk struct {
	Index int
	Start int64
	End   int64
}

// SplitIntoChunks divides the file into chunks according the number
// of threads.
func SplitIntoChunks(fileSize int64, threads int) []Chunk {
	if threads < 1 {
		threads = 1
	}

	chunkSize := fileSize / int64(threads)
	chunks := make([]Chunk, 0, threads)

	var start int64

	for i := range threads {
		end := start + chunkSize - 1

		// The last chunk must go to the end of the file
		if i == threads-1 {
			end = fileSize - 1
		}

		chunks = append(chunks, Chunk{
			Index: i,
			Start: start,
			End:   end,
		})

		start = end + 1
	}

	return chunks
}
