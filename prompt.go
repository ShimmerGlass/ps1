package main

import (
	"flag"
	"os/user"
	"time"
)

func prompt() (res []string) {
	defer measure("prompt", time.Now())

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
