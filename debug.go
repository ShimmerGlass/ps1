package main

import (
	"fmt"
	"sync"
	"time"
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

type debugTime struct {
	Label    string
	Duration time.Duration
}

var debugTimes []debugTime

func measure(label string, since time.Time) {
	debugTimes = append(debugTimes, debugTime{
		Label:    label,
		Duration: time.Since(since),
	})
}

func printDebugTimes() {
	for _, d := range debugTimes {
		fmt.Println(d.Label, "\t", d.Duration)
	}
}
