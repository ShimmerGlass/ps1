package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func rubyVersion(base string) (res []string) {
	defer measure("ruby", time.Now())

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

	res = append(res, color(strings.TrimSpace(p[1]), mustRgbTo256("#A51401"), false))
	return
}
