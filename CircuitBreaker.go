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

func (cb *CircuitBreaker) handleClose() (string, error) {
	cb.breakerState = CLOSE
	return cb.call()
}

func (cb *CircuitBreaker) handleOpen(error error) (string, error) {
	cb.breakerState = OPEN
	cb.failCount = MAX_FAIL_REQUESTS
	cb.error = error
	return cb.call()
}

func (cb *CircuitBreaker) handleHalf() (string, error) {
	cb.breakerState = HALF
	cb.successCount = MAX_SUCCESS_REQUESTS
	return cb.call()
}

func (cb *CircuitBreaker) call() (string, error) {
	if cb.breakerState == CLOSE {
		return cb.closeRequestHandler()
	} else if cb.breakerState == HALF {
		return cb.halfRequestHandler()
	} else if cb.breakerState == OPEN {
		return cb.openRequestHandler()
	}
	panic(fmt.Sprintf("Unknown BreakerState:[%d]", cb.breakerState))
}

func (cb *CircuitBreaker) closeRequestHandler() (string, error) {
	var error error = nil
	for i := 0; i < MAX_ATTEMPT_REQUESTS; i++ {
		respond, err := cb.requestHandler()
		if err == nil {
			return "CLOSE:" + respond, nil
		} else {
			error = err
			time.Sleep(time.Second)
		}
	}
	return cb.handleOpen(error)
}

func (cb *CircuitBreaker) halfRequestHandler() (string, error) {
	if cb.successCount > 0 {
		cb.successCount--
		respond, err := cb.requestHandler()
		if err == nil {
			return "HALF:" + respond, nil
		} else {
			return cb.handleOpen(err)
		}
	} else {
		return cb.handleClose()
	}
}

func (cb *CircuitBreaker) openRequestHandler() (string, error) {
	if cb.failCount > 0 {
		cb.failCount--
		return "OPEN:", cb.error
	} else {
		return cb.handleHalf()
	}
}
