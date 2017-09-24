package main

import "os"

func pprompt() {
	arrowColor := White
	if len(os.Args) > 1 && os.Args[1] != "0" {
		arrowColor = Red
	}

	pcolor(" →  ", arrowColor, true)
}
