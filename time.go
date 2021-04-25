package main

import (
	"fmt"
	"time"
)

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
