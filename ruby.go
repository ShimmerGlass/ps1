package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func prubyVersion(base string) {
	_, err := os.Stat(filepath.Join(base, "Gemfile"))
	if err != nil {
		return
	}

	out, err := exec.Command("rvm", "current").Output()
	if err != nil {
		return
	}

	v := strings.Replace(string(out), "ruby-", "", -1)

	pcolor("(rb ", Red, false)
	pcolor(strings.TrimSpace(v), Red, true)
	pcolor(") ", Red, false)
}
