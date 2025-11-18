package progress

import "sync"

type ChunkProgress struct {
	Index      int
	Downloaded int64
	Total      int64
}

type State struct {
	Chunks     []ChunkProgress
	TotalBytes int64
	mu         sync.Mutex
}

func NewState(chunkCount int) *State {
	return &State{
		Chunks: make([]ChunkProgress, chunkCount),
	}
}

func (s *State) Update(cp ChunkProgress) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Chunks[cp.Index] = cp

	var sum int64
	for _, ch := range s.Chunks {
		sum += ch.Downloaded
	}
	s.TotalBytes = sum
}

func (s *State) Snapshot() (chunks []ChunkProgress, total int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cp := make([]ChunkProgress, len(s.Chunks))
	copy(cp, s.Chunks)

	return cp, s.TotalBytes
}
