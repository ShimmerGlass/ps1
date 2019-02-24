package main

import (
	"fmt"
	"math"
	"os"
	"strings"
)

type term int

const (
	termBash term = iota
	termZsh
)

var currentTerm term

func init() {
	switch {
	case os.Getenv("ZSH_NAME") != "":
		currentTerm = termZsh
	default:
		currentTerm = termBash
	}
}

type colorCode string

const (
	Black  colorCode = "\x1B[%s30m"
	Red    colorCode = "\x1B[%s31m"
	Green  colorCode = "\x1B[%s32m"
	Yellow colorCode = "\x1B[%s33m"
	Blue   colorCode = "\x1B[%s34m"
	Purple colorCode = "\x1B[%s35m"
	Cyan   colorCode = "\x1B[%s36m"
	White  colorCode = "\x1B[%s37m"

	rst = "\x1B[0m"

	escBashStart = "\x01"
	escBashEnd   = "\x02"

	escZshStart = "%{"
	escZshEnd   = "%}"
)

func color(s string, code colorCode, bold bool) string {
	p := "0;"
	if bold {
		p = "1;"
	}

	return escStart() + fmt.Sprintf(string(code), p) + escEnd() + s
}

func colorRst() string {
	return escStart() + string(rst) + escEnd()
}

func pcolor(s string, code colorCode, bold bool) {
	os.Stdout.WriteString(color(s, code, bold))
}

func pcolorRst() {
	os.Stdout.WriteString(colorRst())
}

func ptitle(title string) {
	fmt.Printf("%s\x1B]0;%s\x07%s", escStart(), title, escEnd())
}

func pjobs() {
	if len(os.Args) < 3 {
		return
	}
	j := strings.TrimSpace(os.Args[2])

	if j != "0" {
		pcolor(j+" ", Yellow, false)
	}
}

func escStart() string {
	switch currentTerm {
	case termBash:
		return escBashStart
	case termZsh:
		return escZshStart
	}
	return ""
}

func escEnd() string {
	switch currentTerm {
	case termBash:
		return escBashEnd
	case termZsh:
		return escZshEnd
	}
	return ""
}

func rainbow(s string, l, hstart, hend float64) string {
	span := 0.0
	if hend > hstart {
		span = hend - hstart
	} else {
		span = 360 - hstart + hend
	}
	step := span / float64(len(s))

	res := ""
	for i, c := range s {
		cid := 0
		cscore := math.MaxFloat64
		wantedH := (step*float64(i) + hstart)
		if wantedH > 360 {
			wantedH -= 360
		}
		for currentID, color := range TermCodes256 {
			score := (math.Abs(color.H-wantedH)/360)*2 + (math.Abs(float64(color.L)-l)/100)*3 + (float64(100-color.S)/100)*3
			if score < cscore {
				cid = currentID
				cscore = score
			}
		}

		res += fmt.Sprintf("%s\x1b[38;5;%dm%s%s", escStart(), cid, escEnd(), string(c))
	}

	return res
}
