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


func pack(header string, length int, body string) []byte {
	lengthStr := strconv.Itoa(length)
	content := header + "\r\n" + lengthStr + "\r\n" + body + "\r\n"
	return []byte(content)
}

// seeder connect to dispatcher
func Connect(conn net.Conn, role string) error {
	header := "Connect " + role + " seeder0"
	length := 0
	var body string

	content := pack(header, length, body)
	if _, err := conn.Write(content); err != nil {
		return err
	}

	reader := bufio.NewReader(conn)
	if err := CheckResponse(reader); err != nil {
		return err
	}

	return nil
}

// seed to dispatcher
func seed(body string, conn net.Conn) error {
	header := "Seed Seeder seeder0"
	length := 0

	content := pack(header, length, body)
	if _, err := conn.Write(content); err != nil {
		return err
	}

	return nil
}


// check dispatcher response
func CheckResponse(reader *bufio.Reader) error {
	response, err := ParseRes(reader)
	if err != nil {
		return err
	}

	if response.Status == "Cool" {
		return nil
	} else {
		return errors.New("server is not prepared")
	}

}


func ParseRes(r *bufio.Reader) (*DispatcherResponse, error) {
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