package progress

import (
	"fmt"
	"time"
)

type Renderer struct {
	state      *State
	chunkCount int
	width      int
	interval   time.Duration
}

func NewRenderer(chunkCount int) *Renderer {
	return &Renderer{
		state:      NewState(chunkCount),
		chunkCount: chunkCount,
		width:      30,
		interval:   150 * time.Millisecond,
	}
}

func (r *Renderer) Run(in <-chan ChunkProgress) {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	for {
		select {
		case cp, ok := <-in:
			if !ok {
				r.renderFinal()
				return
			}
			r.state.Update(cp)

		case <-ticker.C:
			r.render()
		}
	}
}

func (r *Renderer) render() {
	chunks, total := r.state.Snapshot()

	for i := range r.chunkCount + 2 {
		fmt.Print("\r\033[K") // clear string

		if i > 0 {
			fmt.Print("\033[F")
		}
	}

	// Render chunks
	for _, ch := range chunks {
		frac := float64(ch.Downloaded) / float64(ch.Total)
		bar := DrawBar(r.width, frac)
		fmt.Printf("Chunk #%d: %s %3.0f%%  (%.2f MB / %.2fMB)\n",
			ch.Index,
			bar,
			frac*100,
			float64(ch.Downloaded)/1024/1024,
			float64(ch.Total)/1024/1024,
		)
	}

	var fullTotal int64
	for _, ch := range chunks {
		fullTotal += ch.Total
	}

	totalFrac := float64(total) / float64(fullTotal)
	totalBar := DrawBar(r.width, totalFrac)

	fmt.Printf("Total:    %s %3.0f%%  (%.2f MB / %.2f MB)\n",
		totalBar,
		totalFrac*100,
		float64(total)/1024/1024,
		float64(fullTotal)/1024/1024,
	)
}

func (r *Renderer) renderFinal() {
	r.render()
	fmt.Println()
	fmt.Println()
}
