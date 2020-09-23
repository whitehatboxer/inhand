package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

type Response struct {
	Body   string
	Status int
}

func Crawl(uri string) (res Response) {
	resp, err := http.Get(uri)
	if err != nil {
		log.Printf("request uri fail: %v, uri: %s\n", err, uri)
	}

	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Printf("read body fail: %v, uri: %s\n", err, uri)
	}

	res = Response{
		Body: string(b),
		Status: resp.StatusCode,
	}
	return res
}