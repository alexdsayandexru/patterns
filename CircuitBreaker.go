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
const MAX_FAIL_REQUESTS = 5

type RequestHandler = func() (string, error)
type ResponseCallback = func(string, error)
type BreakerState = int

type CircuitBreaker struct {
	breakerState     BreakerState
	requestHandler   RequestHandler
	responseCallback ResponseCallback
	successCount     int
	failCount        int
	error            error
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
	cb.breakerState = OPEN
	cb.failCount = MAX_FAIL_REQUESTS
	cb.error = error
	cb.call(i)
}

func (cb *CircuitBreaker) handleHalf(i int) {
	cb.breakerState = HALF
	cb.successCount = MAX_SUCCESS_REQUESTS
	cb.call(i)
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
	var err error = nil
	for j := 0; j < MAX_ATTEMPT_REQUESTS; j++ {
		response, err := cb.requestHandler()
		if err == nil {
			cb.responseHandler(i, "CLOSE", response, nil)
			return
		} else {
			err = err
			time.Sleep(time.Second)
		}
	}
	cb.handleOpen(i, err)
}

func (cb *CircuitBreaker) halfRequestHandler(i int) {
	if cb.successCount > 0 {
		cb.successCount--
		response, err := cb.requestHandler()
		if err == nil {
			cb.responseHandler(i, "HALF", response, nil)
		} else {
			cb.handleOpen(i, err)
		}
	} else {
		cb.handleClose(i)
	}
}

func (cb *CircuitBreaker) openRequestHandler(i int) {
	if cb.failCount > 0 {
		cb.failCount--
		cb.responseHandler(i, "OPEN", "", cb.error)
	} else {
		cb.handleHalf(i)
	}
}
