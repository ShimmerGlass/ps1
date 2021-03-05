package main

import (
	"os"
	"os/user"
)

func prompt() (res []string) {
	arrowColor := White
	if len(os.Args) > 1 && os.Args[1] != "0" {
		arrowColor = Red
	}

	usr, _ := user.Current()

	if usr.Uid == "0" {
		res = append(res, color("#", arrowColor, true))
	} else {
		res = append(res, color("➠", arrowColor, true))
	}

	return
}
