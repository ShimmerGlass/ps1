package main

import (
	"os"
	"time"
)

func ssh() (res []string) {
	defer measure("ssh", time.Now())

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
