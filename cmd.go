package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func run(bin string, args ...string) (string, error) {
	res, err := exec.Command(bin, args...).Output()
	if err != nil {
		return "", fmt.Errorf("exec: %s %s: %s", bin, strings.Join(args, " "), err)
	}

	return strings.TrimSuffix(string(res), "\n"), nil
}
