package main

import (
	"fmt"
	"time"
)

const CLOSE = 0
const HALF = 1
const OPEN = 2

const MAX_ATTEMPT_REQUESTS = 5
const MAX_SUCCESS_REQUESTS = 10
const OPEN_STATE_SECONDS = 5

type RequestHandler = func() (string, error)
type ResponseCallback = func(string, error)
type BreakerState = int

type CircuitBreaker struct {
	breakerState         BreakerState
	requestHandler       RequestHandler
	responseCallback     ResponseCallback
	successRequestsCount int
	error                error
}

func NewCircuitBreaker(handler RequestHandler, callback ResponseCallback) *CircuitBreaker {
	cb := CircuitBreaker{
		requestHandler:   handler,
		responseCallback: callback,
		breakerState:     CLOSE,
	}
	return &cb
}

func (cb *CircuitBreaker) responseHandler(i int, state string, response string, err error) {
	cb.responseCallback(fmt.Sprintf("[%d] %s %s", i, state, response), err)
}

func (cb *CircuitBreaker) handleClose(i int) {
	cb.breakerState = CLOSE
	cb.call(i)
}

func (cb *CircuitBreaker) handleOpen(i int, error error) {
	timer := time.NewTimer(time.Second * OPEN_STATE_SECONDS)
	go func() {
		<-timer.C
		cb.handleHalf(i)
	}()
	cb.error = error
	cb.breakerState = OPEN
	cb.call(i)
}

func (cb *CircuitBreaker) handleHalf(i int) {
	cb.successRequestsCount = MAX_SUCCESS_REQUESTS
	cb.breakerState = HALF
}

func (cb *CircuitBreaker) call(i int) {
	switch cb.breakerState {
	case CLOSE:
		cb.closeRequestHandler(i)
	case HALF:
		cb.halfRequestHandler(i)
	case OPEN:
		cb.openRequestHandler(i)
	default:
		panic(fmt.Sprintf("Unknown BreakerState:[%d]", cb.breakerState))
	}
}

func (cb *CircuitBreaker) closeRequestHandler(i int) {
	var error error
	for j := 0; j < MAX_ATTEMPT_REQUESTS; j++ {
		response, e := cb.requestHandler()
		if e == nil {
			cb.responseHandler(i, "CLOSE:", response, nil)
			return
		} else {
			error = e
			time.Sleep(time.Second)
		}
	}
	cb.handleOpen(i, error)
}

func (cb *CircuitBreaker) halfRequestHandler(i int) {
	if cb.successRequestsCount > 0 {
		cb.successRequestsCount--
		response, err := cb.requestHandler()
		if err == nil {
			cb.responseHandler(i, "HALF:", response, nil)
		} else {
			cb.handleOpen(i, err)
		}
	} else {
		cb.handleClose(i)
	}
}

func (cb *CircuitBreaker) openRequestHandler(i int) {
	cb.responseHandler(i, "OPEN:", "", cb.error)
}
