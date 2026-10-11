package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// durationThreshold is the minimum command runtime before we surface how long
// it took in the prompt. Short commands stay out of the way.
const durationThreshold = 10 * time.Second

// cmdDuration renders the previous command's wall-clock runtime when it exceeds
// durationThreshold. The elapsed time (in fractional seconds) is passed as the
// third positional arg by the shell's preexec/precmd timing hooks.
func cmdDuration() (res []string) {
	raw := strings.TrimSpace(flag.Arg(2))
	if raw == "" {
		return
	}

	secs, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		errorAdd(fmt.Errorf("parse duration %q: %s", raw, err))
		return
	}

	d := time.Duration(secs * float64(time.Second))
	if d < durationThreshold {
		return
	}

	res = append(res, color("⧗ "+fmtDuration(d), Warning, false))
	return
}

func fmtDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60

	switch {
	case h > 0:
		return fmt.Sprintf("%dh%dm%ds", h, m, s)
	case m > 0:
		return fmt.Sprintf("%dm%ds", m, s)
	default:
		return fmt.Sprintf("%ds", s)
	}
}
