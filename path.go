package main

import (
	"os"
	"os/user"
	"strings"
)

func pathExists(p string) bool {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return false
	}

	return true
}

func pcwd(path, from string) {
	pathParts := strings.Split(path, "/")
	fromParts := strings.Split(from, "/")

	if len(pathParts) < len(fromParts) {
		pcolor(path, Cyan, false)
	}

	for i := range fromParts {
		if pathParts[i] != fromParts[i] {
			pcolor(path, Cyan, false)
		}
	}

	var missingParts []string
	var okParts []string

	for i := len(pathParts); i > 0; i-- {
		path := strings.Join(pathParts[:i], "/")
		if pathExists(path) {
			okParts = pathParts[len(fromParts):i]
			missingParts = pathParts[i:]
			break
		}
	}

	res := "~"
	if len(okParts) > 0 {
		res += "/" + strings.Join(okParts, "/")
	}
	if len(missingParts) > 0 {
		res += "/" + color(strings.Join(missingParts, "/"), Red, false)
	}

	pcolor(res, Cyan, false)
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
