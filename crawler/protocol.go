package main

import (
	"bufio"
	"errors"
	"net"
	"strconv"
	"strings"
)

type DispatcherResponse struct {
	Status string
	Role   string
	ID     string
	Length int32
	Body   string
}

type DispatcherRequest struct {
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

// crawler connect to dispatcher
func connect(conn net.Conn, role string) error {
	header := "Connect " + role + " crawler0"
	length := 0
	var body string

	content := pack(header, length, body)
	if _, err := conn.Write(content); err != nil {
		return err
	}

	reader := bufio.NewReader(conn)
	if err := checkResponse(reader); err != nil {
		return err
	}

	return nil
}

// crawler upload to dispatcher
func upload(res Response, conn net.Conn) error {
	header := "Upload Crawler crawler0"
	length := 0
	body := res.Body

	content := pack(header, length, body)

	if _, err := conn.Write(content); err != nil {
		return err
	}

	return nil
}


// check dispatcher response
func checkResponse(reader *bufio.Reader) error {
	response, err := parseRes(reader)
	if err != nil {
		return err
	}

	if response.Status == "Cool" {
		return nil
	} else {
		return errors.New("server is not prepared")
	}

}

// parse request from dispatcher which is about sending task
func parseRequest(r *bufio.Reader) (* DispatcherRequest, error) {
	request := &DispatcherRequest{
		Method: "Crawl",
		Role:   "Dispatcher",
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

		if len(line) <= 0 {
			break
		}

		switch lineCount {
		case 0: {
			line = strings.TrimRight(line, "\r\n")
			items := strings.Split(line, " ")
			if len(items) != 3 {
				return nil, errors.New("bad protocol")
			}
			request.Method = items[0]
			request.Role = items[1]
			request.ID = items[2]
		}
		case 1: {
			// no need to do since we havn't use length yet.
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


func parseRes(r *bufio.Reader) (*DispatcherResponse, error) {
	response := &DispatcherResponse{
		Status: "",
		Role: "dispatcher",
		Length: 0,
		Body: "",
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

		if len(line) <= 0 {
			break
		}

		switch lineCount {
		case 0: {
			line = strings.TrimRight(line, "\r\n")
			items := strings.Split(line, " ")
			if len(items) != 3 {
				return nil, errors.New("bad protocol")
			}
			response.Status = items[0]
			response.Role = items[1]
			response.ID = items[2]
		}
		case 1: {
			// no need to do since we havn't use length yet.
		}
		case 2: {
			line = strings.TrimRight(line, "\r\n")
			response.Body = line
		}
		}

		lineCount += 1
	}

	return response, nil
}