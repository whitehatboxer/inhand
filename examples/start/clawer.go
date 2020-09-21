package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	url := "http://www.douban.com"

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)

	// set Header
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) "+
		"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.102 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("http get error: ", err)
	}

	bodyInBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read body error: ", err)
	}

	body := string(bodyInBytes)
	fmt.Println(body)
}
