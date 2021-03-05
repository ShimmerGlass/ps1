package main

import (
	"os"
)

func ssh() (res []string) {
	if os.Getenv("SSH_CLIENT") == "" {
		return
	}

	hostname, err := os.Hostname()
	if err != nil {
		return
	}

	res = append(res,
		color("@", Black, true)+color(hostname, Purple, false),
	)
	return
}
