package main

import (
	"bufio"
	"errors"
	"net"
	"strconv"
	"strings"
)


// struct described connection from other module
type ModuleRequest struct {
	Method string
	Role   string
	ID     string
	Length int32
	Body   string
}


func pack(header string, length int, body string) []byte {
	lengthStr := strconv.Itoa(length)
	content := header + "\r\n" + lengthStr + "\r\n" + body + "\r\n"
	return []byte(content)
}


// Send Crawl
func Crawl(url, id string, conn net.Conn) error {
	header := "Crawl Dispatcher " + id
	length := 0
	body := url

	content := pack(header, length, body)
	if _, err := conn.Write(content); err != nil {
		return err
	}

	return nil
}

// Accept module connection
func AcceptConnect(conn net.Conn) (*ModuleRequest, error) {
	// parse request and say cool if request is valid
	reader := bufio.NewReader(conn)
	request, err := ParseReq(reader)
	if err != nil {
		return nil, err
	}

	header := "Cool Dispatcher Dispatcher0"
	length := 0
	body := ""
	content := pack(header, length, body)
	if _, err := conn.Write(content); err != nil {
		return nil, err
	}

	return request, nil
}

// Accept seeder seed
func acceptSeed(conn net.Conn) (*Seed, error) {
	reader := bufio.NewReader(conn)
	request, err := ParseReq(reader)
	if err != nil {
		return nil, err
	}

	seed := &Seed{
		MaxDepth: 1,
		Depth:    1,
		ID:       request.ID,
		Urls:     []string{request.Body},
	}

	return seed, nil
}

// Parse module request
func ParseReq(r *bufio.Reader) (*ModuleRequest, error) {
	request := &ModuleRequest{
		Method: "Connect",
		Role:   "",
		ID:     "",
		Length: 0,
		Body:   "",
	}

	lineCount := 0
	for {
		if lineCount > 2 {
			break
		}

		line, err := r.ReadString('\n')
		if err != nil {
			return nil, err
		}

		if len(line) <= 0{
			break
		}

		switch lineCount {
		case 0: {
			line = strings.TrimRight(line, "\r\n")
			items := strings.Split(line, " ")
			if len(items) != 3 {
				return nil, errors.New("bad protocol header1")
			}
			request.Method = items[0]
			request.Role = items[1]
			request.ID = items[2]
		}
		case 1: {
			// no need to do since we haven't use length yet.
		}
		case 2: {
			line = strings.TrimRight(line, "\r\n")
			request.Body = line
		}
		}

		lineCount += 1
	}

	return request, nil
}

