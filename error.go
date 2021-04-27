package main

import (
	"fmt"
	"sync"
)

var errors []error
var errorsLock sync.Mutex

func errorAdd(err error) {
	errorsLock.Lock()
	defer errorsLock.Unlock()

	errors = append(errors, err)
}

func errorCount() (res []string) {
	if len(errors) > 0 {
		res = append(res, color(fmt.Sprintf("⚠ %d", len(errors)), Danger, false))
	}
	return
}

func printErrors() {
	for _, err := range errors {
		fmt.Println(err.Error())
	}
}
