package main

import (
	"os"
	"os/user"
	"strings"
)

type prettyPath struct {
	ok      []string
	missing []string
}

func (p *prettyPath) print() {
	parts := p.ok
	for _, p := range p.missing {
		parts = append(parts, color(p, Red, false))
	}

	pcolor(strings.Join(parts, color("/", Cyan, false)), Cyan, false)
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
	cwd, ok := os.LookupEnv("PWD")
	if ok {
		return cwd
	}

	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}

	return cwd
}

func home() string {
	usr, _ := user.Current()
	return usr.HomeDir
}
