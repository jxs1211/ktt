package tool

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
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

// IsPortOpen checks if a port is open on the specified host.
func IsPortOpen(host string, port int) (bool, error) {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", address, 2*time.Second) // Set a timeout for the connection
	if err != nil {
		// If the error is a timeout or connection refused, the port is not open
		if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
			log.Info("IsPortOpen", "err", opErr)
			return false, nil
		}
		log.Info("IsPortOpen", "other err", err)
		return false, err
	}
	defer conn.Close()
	return true, nil
}

func DBHomePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	dbPath := filepath.Join(homeDir, ".KTT", "db", "sqlite.db")
	return dbPath, nil
}

func GetDBFilePath() (string, error) {
	dbPath, err := DBHomePath()
	if err != nil {
		return "", err
	}
	// Ensure the directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory %s: %w", dir, err)
	}
	// Check if the file exists and is writable
	if _, err := os.Stat(dbPath); err == nil {
		// File exists, check if it's writable
		file, err := os.OpenFile(dbPath, os.O_WRONLY, 0666)
		if err != nil {
			return "", fmt.Errorf("database file exists but is not writable: %w", err)
		}
		file.Close()
	} else if os.IsNotExist(err) {
		// File doesn't exist, try to create it
		file, err := os.Create(dbPath)
		if err != nil {
			return "", fmt.Errorf("failed to create database file: %w", err)
		}
		file.Close()
	} else {
		return "", fmt.Errorf("error checking database file: %w", err)
	}

	return dbPath, nil
}
