package main

import (
	"fmt"
	"sync"
	"time"
)

func doSync(cb *CircuitBreaker) {
	for i := 0; i < 1000; i++ {
		cb.call(i)
	}
}

func doAsync(cb *CircuitBreaker) {
	var waitGroup sync.WaitGroup
	for i := 0; i < 1000; i++ {
		waitGroup.Add(1)
		go func(i int) {
			defer waitGroup.Done()
			cb.call(i)
		}(i)
	}
	waitGroup.Wait()
}

func main() {
	var cb = NewCircuitBreaker(
		func() (string, error) {
			return httpGet("http://localhost:8083/")
		},
		func(s string, err error) {
			fmt.Println(s, err)
			time.Sleep(time.Millisecond * 1000)
		})
	doAsync(cb)
}
