package main

import (
	"fmt"
	"time"
)

func doRequests(cb *CircuitBreaker, count int) {
	for i := 0; i < count; i++ {
		result, err := cb.call()
		fmt.Println(result, err)
		time.Sleep(time.Second)
	}
}

func main() {
	var cb = NewCircuitBreaker(
		func() (string, error) {
			return httpGet("http://localhost:8083/")
		},
		func(s string, err error) {
			fmt.Println(s, err)
		})

	doRequests(cb, 200)
}
