package api

import (
	"fmt"
	"time"
)

func FormatDurationToMinutes(duration time.Duration) string {
	return fmt.Sprintf("%.f minuti", duration.Minutes())
}
