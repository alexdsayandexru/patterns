package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type Result struct {
	result string
	error  error
}

func chanTest() {
	var wg sync.WaitGroup
	c := make(chan Result)
	wg.Add(1)
	go func(c chan Result) {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			c <- Result{
				result: fmt.Sprintf("VAL:%d", i),
				error:  errors.New(fmt.Sprintf("ERR:%d", i)),
			}
		}
		close(c)
	}(c)
	wg.Add(1)
	go func(c chan Result) {
		defer wg.Done()
		for v := range c {
			fmt.Println("1:" + v.result)
			time.Sleep(time.Second)
		}
	}(c)
	wg.Add(1)
	go func(c chan Result) {
		defer wg.Done()
		for v := range c {
			fmt.Println("2:" + v.result)
			time.Sleep(time.Second)
		}
	}(c)
	wg.Wait()
}

func selectTest() {
	c1 := make(chan string)

	go func() {
		for i := 0; i < 8; i++ {
			c1 <- fmt.Sprintf("c1:[%d]", i)
			time.Sleep(4 * time.Second)
		}
	}()

	for {
		select {
		case msg := <-c1:
			fmt.Println(msg)
		case <-time.After(time.Second * 5):
			fmt.Println("timeout")
			return
		}
	}
}
