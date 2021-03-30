package main

import (
	"flag"
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

func (c colorCode) String() string {
	return string(c)
}

func (c *colorCode) Set(s string) error {
	v, err := rgbTo256(s)
	if err != nil {
		return err
	}

	*c = v
	return nil
}

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

var (
	Accent  colorCode = Cyan
	Text    colorCode = White
	Neutral colorCode = Black
	Danger  colorCode = Red
	Warning colorCode = Purple
	Success colorCode = Green
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
	j := strings.TrimSpace(flag.Arg(1))

	if j != "0" && j != "" {
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

func mustRgbTo256(hex string) colorCode {
	c, err := rgbTo256(hex)
	if err != nil {
		panic(err)
	}

	return c
}

func rgbTo256(hex string) (colorCode, error) {
	var r, g, b uint8
	_, err := fmt.Sscanf(hex, "#%02x%02x%02x", &r, &g, &b)
	if err != nil {
		return "", err
	}

	r = uint8((float64(r) / 255) * 5)
	g = uint8((float64(g) / 255) * 5)
	b = uint8((float64(b) / 255) * 5)

	n := 16 + 36*r + 6*g + b

	return colorCode("\x1b[%s38;5;"+fmt.Sprint(n)) + "m", nil
}
