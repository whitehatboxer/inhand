package clawer

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Response struct {
	Body string
	Status int
}

func Claw(uri string) (res Response) {
	resp, err := http.Get(uri)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
		os.Exit(1)
	}

	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", uri, err)
		os.Exit(1)
	}
	//fmt.Printf("%s", b)

	res = Response{
		Body: string(b),
		Status: resp.StatusCode,
	}
	return res
}
