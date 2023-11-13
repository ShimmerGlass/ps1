package main

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

var versionned = []func(string) []string{
	rubyVersion,
	goVersion,
	pythonVersion,
}

func versions(root string) (res []string) {
	for _, ver := range versionned {
		res = append(res, ver(root)...)
	}

	return
}

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

	version := strings.TrimSpace(p[1])
	version, _, _ = strings.Cut(version, "p")

	res = append(res, color(version, mustRgbTo256("#A51401"), false))
	return
}

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

func pythonVersion(base string) (res []string) {
	defer measure("python", time.Now())

	_, err := os.Stat(filepath.Join(base, "setup.py"))
	if os.IsNotExist(err) {
		return
	}
	if err != nil {
		errorAdd(err)
		return
	}

	out, err := run("python3", "--version")
	if err != nil {
		errorAdd(err)
		return
	}

	_, version, _ := strings.Cut(out, " ")

	res = append(res, color(strings.TrimSpace(version), mustRgbTo256("#ffd343"), false))
	return
}
