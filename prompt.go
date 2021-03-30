package main

import (
	"flag"
	"os/user"
)

func prompt() (res []string) {
	arrowColor := Accent
	if flag.Arg(0) != "0" {
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
