package main

import (
	"fmt"
	"sync"
	"time"
	"errors"
)

func doSync(cb *CircuitBreaker) {
	for i := 0; i < 100; i++ {
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

type Result struct {
	result string
	error error
}

func main() {
	
	return
	var cb = NewCircuitBreaker(
		func() (string, error) {
			return httpGet("http://localhost:8888/")
		},
		func(s string, err error) {
			fmt.Println(s, err)
			time.Sleep(time.Millisecond * 1000)
		})
	doSync(cb)
}

func chanTest () {
	var wg sync.WaitGroup 
	c := make(chan Result) 
	wg.Add(1)
	go func(c chan Result) {
		defer wg.Done()
		for i:=0; i<10; i++ {
			c<-Result {
				result: fmt.Sprintf("VAL:%d", i),
				error: errors.New(fmt.Sprintf("ERR:%d", i)),
			}
		}
		close(c)
	}(c)
	wg.Add(1)
	go func(c chan Result) {
		defer wg.Done()
		for v := range c {
			fmt.Println("1:" + v.result)
			//time.Sleep(time.Second)
		}
	}(c)
	wg.Add(1)
	go func(c chan Result) {
		defer wg.Done()
		for v := range c {
			fmt.Println("2:" + v.result)
			//time.Sleep(time.Second)
		}
	}(c)
	wg.Wait() 
}