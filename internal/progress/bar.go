package progress

import (
	"fmt"
	"math"
	"strings"
)

func DrawBar(width int, fraction float64) string {
	if math.IsNaN(fraction) || fraction < 0 {
		fraction = 0
	}
	if fraction > 1 {
		fraction = 1
	}

	filled := int(float64(width) * fraction)
	if filled < 0 {
		filled = 0
	}
	if filled > width {
		filled = width
	}

	empty := width - filled

	return fmt.Sprintf("[%s%s]",
		strings.Repeat("█", filled),
		strings.Repeat("░", empty),
	)
}
