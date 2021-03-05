package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func rubyVersion(base string) (res []string) {
	_, err := os.Stat(filepath.Join(base, "Gemfile"))
	if err != nil {
		return
	}

	out, err := exec.Command("ruby", "-v").Output()
	if err != nil {
		return
	}

	p := strings.Split(string(out), " ")
	if len(p) < 2 {
		return
	}

	res = append(res, color(strings.TrimSpace(p[1]), Red, false))
	return
}
