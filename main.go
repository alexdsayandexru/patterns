package main

import (
	"fmt"
	"strings"
	"time"
)

func doSync() {
	var cb = NewCircuitBreaker(
		func() (string, error) {
			return httpGet("http://localhost:8080/")
		},
		func(s string, err error) {
			if err != nil {
				fmt.Println(s, err)
			} else {
				fmt.Println(strings.TrimRight(s, "\n"))
			}
			time.Sleep(time.Millisecond * 1000)
		})

	for i := 0; i < 1000; i++ {
		cb.call(i)
	}
}

func main() {
	//doSync()
	selectTest()
}
