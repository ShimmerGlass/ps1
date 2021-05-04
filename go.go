package main

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

func goVersion(base string) (res []string) {
	defer measure("go", time.Now())

	_, err := os.Stat(filepath.Join(base, "go.mod"))
	if os.IsNotExist(err) {
		return
	}
	if err != nil {
		errorAdd(err)
		return
	}

	out, err := run("go", "version")
	if err != nil {
		errorAdd(err)
		return
	}

	p := strings.Split(out, " ")
	if len(p) < 3 {
		return
	}

	v := strings.Replace(p[2], "go", "", 1)

	res = append(res, color(strings.TrimSpace(v), mustRgbTo256("#67D0DE"), false))
	return
}
