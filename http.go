package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

func httpGet(requestURL string) (string, error) {
	resp, err := http.Get(requestURL)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("failed to read body: %s", err)
	}

	return string(body), nil
}
