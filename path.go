package main

import (
	"os"
	"os/user"
	"strings"
)

type prettyPath struct {
	abs      string
	relative string
	ok       []string
	missing  []string
}

func (p *prettyPath) print() {
	var res string
	if p.relative == "" {
		res = ""
	} else {
		res = "~"
		if len(p.ok) > 0 || len(p.missing) > 0 {
			res += "/"
		}
	}

	parts := p.ok
	for _, p := range p.missing {
		parts = append(parts, color(p, Red, false))
	}

	res += strings.Join(parts, color("/", Cyan, false))

	pcolor(res, Cyan, false)
}

func (p *prettyPath) string() string {
	if len(p.relative) > 0 {
		return p.relative
	}

	return p.abs
}

func pathExists(p string) bool {
	_, err := os.Stat(p)
	if os.IsNotExist(err) {
		return false
	}

	return true
}

func newPrettyPath(path, from string) prettyPath {
	prettyPath := prettyPath{
		abs: path,
	}

	pathParts := strings.Split(path, "/")
	fromParts := strings.Split(from, "/")

	if len(pathParts) < len(fromParts) {
		return prettyPath
	}

	for i := range fromParts {
		if pathParts[i] != fromParts[i] {
			fromParts = []string{}
			break
		}
	}

	if len(fromParts) > 0 {
		prettyPath.relative = "~"
		if len(pathParts[len(fromParts):]) > 0 {
			prettyPath.relative += "/" + strings.Join(pathParts[len(fromParts):], "/")
		}
	}

	prettyPath.missing = pathParts
	for i := len(pathParts); i > 0; i-- {
		path := strings.Join(pathParts[:i], "/")
		if pathExists(path) {
			prettyPath.ok = pathParts[len(fromParts):i]
			prettyPath.missing = pathParts[i:]
			break
		}
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
