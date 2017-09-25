package main

import (
	"fmt"
	"os"
)

type colorCode string

const (
	Black  colorCode = "30m" // Black - Regular
	Red    colorCode = "31m" // Red
	Green  colorCode = "32m" // Green
	Yellow colorCode = "33m" // Yellow
	Blue   colorCode = "34m" // Blue
	Purple colorCode = "35m" // Purple
	Cyan   colorCode = "36m" // Cyan
	White  colorCode = "37m" // White

	rst = "0m"

	escStart = "\x01\x1B["
	escEnd   = "\x02"
)

func color(s string, code colorCode, bold bool) string {
	p := "0;"
	if bold {
		p = "1;"
	}

	return escStart + p + string(code) + escEnd + s
}

func colorRst() string {
	return escStart + string(rst) + escEnd
}

func pcolor(s string, code colorCode, bold bool) {
	os.Stdout.WriteString(color(s, code, bold))
}

func pcolorRst() {
	os.Stdout.WriteString(colorRst())
}

func ptitle(title string) {
	fmt.Printf("\x01\x1B]0;%s\x07\x02", title)
}

func pjobs() {
	if len(os.Args) > 2 && os.Args[2] != "0" {
		pcolor(os.Args[2]+" ", Yellow, false)
	}
}
