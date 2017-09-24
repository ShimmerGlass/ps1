package main

import (
	"time"
)

func ptime() {
	pcolor(time.Now().Format("15:04:05.000 "), Black, false)
}
