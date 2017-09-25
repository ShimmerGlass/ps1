package main

import (
	"os"
	"os/user"
	"runtime"
)

func pprompt() {
	arrowColor := White
	if len(os.Args) > 1 && os.Args[1] != "0" {
		arrowColor = Red
	}

	usr, _ := user.Current()

	if usr.Uid == "0" {
		pcolor(" # ", arrowColor, true)
	} else if runtime.GOOS == "darwin" {
		pcolor(" → ", arrowColor, true)
	} else {
		pcolor(" →  ", arrowColor, true)
	}
}
