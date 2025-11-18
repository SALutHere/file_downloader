package downloader

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// ChunkState stores the state of the single chunk
type ChunkState struct {
	Index      int   `json:"index"`
	Start      int64 `json:"start"`
	End        int64 `json:"end"`
	Downloaded int64 `json:"downloaded"`
	Done       bool  `json:"done"`
}

// State stores information about the file downloading
type State struct {
	URL    string       `json:"url"`
	Output string       `json:"output"`
	Size   int64        `json:"size"`
	Chunks []ChunkState `json:"chunks"`

	mu sync.Mutex `json:"-"`
}

// Save saves the state to a JSON file (atomically)
func (s *State) Save() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	tmp := StateFileName(s.Output) + ".tmp"

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("state serialization error: %w", err)
	}

	if err = os.WriteFile(tmp, data, 0644); err != nil {
		return fmt.Errorf("temporary state-file writing error: %w", err)
	}

	if err = os.Rename(tmp, StateFileName(s.Output)); err != nil {
		return fmt.Errorf("error replacing state-file: %w", err)
	}

	return nil
}

// StateFileName returns the name of the state file
func StateFileName(output string) string {
	return output + ".state.json"
}

// LoadState reads the state from a file
func LoadState(output string) (*State, error) {
	data, err := os.ReadFile(StateFileName(output))
	if err != nil {
		return nil, fmt.Errorf("error reading state-file: %w", err)
	}

	var s State
	if err = json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("error parsing state-file: %w", err)
	}

	return &s, nil
}
