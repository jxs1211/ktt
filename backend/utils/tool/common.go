package tool

import (
	"fmt"
	"time"
)

func TrackTime(funcName string) func() {
	pre := time.Now()
	return func() {
		elapsed := time.Since(pre)
		var duration string
		switch {
		case elapsed < time.Millisecond:
			duration = fmt.Sprintf("%.2f µs", float64(elapsed.Nanoseconds())/1000)
		case elapsed < time.Second:
			duration = fmt.Sprintf("%.2f ms", float64(elapsed.Nanoseconds())/1000000)
		default:
			duration = fmt.Sprintf("%.2f s", elapsed.Seconds())
		}
		fmt.Printf("exec [%s] elapsed: %s\n", funcName, duration)
	}
}
