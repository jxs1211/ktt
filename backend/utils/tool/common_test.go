package tool

import (
	"fmt"
	"testing"
	"time"

	"golang.org/x/sync/singleflight"
)

func FetchExpensiveData() (int64, error) {
	fmt.Println("FetchExpensiveData called", time.Now())
	time.Sleep(1 * time.Second)
	return time.Now().Unix(), nil
}

var group singleflight.Group

func UsingSingleFlight(key string) {
	v, _, _ := group.Do(key, func() (interface{}, error) {
		return FetchExpensiveData()
	})
	fmt.Println(v)
}

func TestUsingSingleFlight(t *testing.T) {
	go UsingSingleFlight("key")
	go UsingSingleFlight("key")
	go UsingSingleFlight("key")

	time.Sleep(2 * time.Second)

	go UsingSingleFlight("key")
	go UsingSingleFlight("key")
	go UsingSingleFlight("key")

	time.Sleep(2 * time.Second)
}

func TestIsPortOpen(t *testing.T) {
	tests := []struct {
		name    string
		host    string
		port    int
		want    bool
		wantErr bool
	}{
		{
			name: "base",
			host: "127.0.0.1",
			port: 12111,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsPortOpen(tt.host, tt.port)
			t.Log(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsPortOpen() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsPortOpen() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestXxx(t *testing.T) {
	t.Log("testest")
}
