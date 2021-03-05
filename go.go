package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func goVersion(base string) (res []string) {
	_, err := os.Stat(filepath.Join(base, "go.mod"))
	if err != nil {
		return
	}

	out, err := exec.Command("go", "version").Output()
	if err != nil {
		return
	}

	p := strings.Split(string(out), " ")
	if len(p) < 3 {
		return
	}

	v := strings.Replace(p[2], "go", "", 1)

	res = append(res, color(strings.TrimSpace(v), Cyan, false))
	return
}
