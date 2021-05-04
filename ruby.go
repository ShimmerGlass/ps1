package main

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

func rubyVersion(base string) (res []string) {
	defer measure("ruby", time.Now())

	_, err := os.Stat(filepath.Join(base, "Gemfile"))
	if os.IsNotExist(err) {
		return
	}
	if err != nil {
		errorAdd(err)
		return
	}

	out, err := run("ruby", "-v")
	if err != nil {
		errorAdd(err)
		return
	}

	p := strings.Split(out, " ")
	if len(p) < 2 {
		return
	}

	res = append(res, color(strings.TrimSpace(p[1]), mustRgbTo256("#A51401"), false))
	return
}
