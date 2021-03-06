package main

import (
	"fmt"
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

func title(title string) string {
	return fmt.Sprintf("%s\x1B]0;%s\x07%s", escStart(), title, escEnd())
}

func jobs() (res []string) {
	if len(os.Args) < 3 {
		return
	}
	j := strings.TrimSpace(os.Args[2])

	if j != "0" {
		res = append(res, color(j, Yellow, false))
	}

	return
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
