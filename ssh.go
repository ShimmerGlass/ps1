package main

import (
	"os"
)

func pssh() {
	if os.Getenv("SSH_CLIENT") == "" {
		return
	}

	hostname, err := os.Hostname()
	if err != nil {
		return
	}

	pcolor("@", Black, true)
	pcolor(hostname, Cyan, false)
}
