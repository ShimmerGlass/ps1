package main

import (
	"flag"
	"os/user"
)

func prompt() (res []string) {
	args := flag.Args()
	arrowColor := Accent
	if len(args) > 1 && args[1] != "0" {
		arrowColor = Danger
	}

	usr, _ := user.Current()

	if usr.Uid == "0" {
		res = append(res, color("#", arrowColor, true))
	} else {
		res = append(res, color("➠", arrowColor, true))
	}

	return
}
