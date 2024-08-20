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
