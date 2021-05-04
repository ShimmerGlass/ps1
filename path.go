package main

import (
	"os"
	"os/user"
	"strings"
	"time"
)

type prettyPath struct {
	ok      []string
	missing []string
}

func (p *prettyPath) print() (res []string) {
	parts := []string{}
	for _, p := range p.ok {
		parts = append(parts, color(p, Cyan, false))
	}
	for _, p := range p.missing {
		parts = append(parts, color(p, Red, false))
	}

	res = append(res, strings.Join(parts, color("/", Black, false)))
	return
}

func (p *prettyPath) string() string {
	parts := append(p.ok, p.missing...)
	return strings.Join(parts, "/")
}

func pathExists(p string) bool {
	_, err := os.Stat(p)
	if os.IsNotExist(err) {
		return false
	}

	return true
}

func newPrettyPath(path, from string) prettyPath {
	defer measure("path", time.Now())

	prettyPath := prettyPath{}

	pathParts := strings.Split(path, "/")
	fromParts := strings.Split(from, "/")

	isRel := strings.HasPrefix(path, from)

	prettyPath.missing = pathParts
	for i := len(pathParts); i > 0; i-- {
		path := strings.Join(pathParts[:i], "/")
		if pathExists(path) {
			prettyPath.ok = pathParts[:i]
			prettyPath.missing = pathParts[i:]
			break
		}
	}

	if isRel && len(prettyPath.ok) >= len(fromParts) {
		prettyPath.ok = append([]string{"~"}, prettyPath.ok[len(fromParts):]...)
	}

	return prettyPath
}

func getCwd() string {
	defer measure("cwd", time.Now())

	cwd, err := os.Getwd()
	if err != nil {
		errorAdd(err)
		return ""
	}

	return cwd
}

func home() string {
	defer measure("home", time.Now())

	usr, _ := user.Current()
	return usr.HomeDir
}
