package tool

import (
	"fmt"
	"time"

	"ktt/backend/utils/log"
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
		log.Info("exec [%s] elapsed: %s\n", "funcName", funcName, "duration", duration)
	}
}

type set map[string]struct{}

var (
	daysWithoutIssue set
	daysTotal        set
)

func InitRunningDays() (int, int) {
	// Get the current date in a human-readable format
	currentDate := time.Now().Format("2006-01-02")
	// Initialize the set if not already done
	if daysWithoutIssue == nil {
		daysWithoutIssue = make(set)
	}
	if daysTotal == nil {
		daysTotal = make(set)
	}
	// if !checkErrorExists() {
	// 	daysWithoutIssue[currentDate] = struct{}{}
	// }
	daysTotal[currentDate] = struct{}{}
	return len(daysWithoutIssue), len(daysTotal)
}
