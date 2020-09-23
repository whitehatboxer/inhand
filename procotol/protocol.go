package main

import (
	"net"
	"strconv"
)

func pack(header string, length int, body string) []byte {
	lengthStr := strconv.Itoa(length)
	content := header + "\n" + lengthStr + "\n" + body
	return []byte(content)
}

// crawler/parser/seeder connect to dispatch
func Connect(conn net.Conn, role string) error {
	header := "Connect " + role
	length := 0
	var body string

	content := pack(header, length, body)
	if _, err := conn.Write(content); err != nil {
		return err
	}

	return nil
}

func Seed(body string, conn net.Conn) error {
	header := "Seed Seeder"
	length := 0

	content := pack(header, length, body)
	if _, err := conn.Write(content); err != nil {
		return err
	}

	return nil
}

func Crawl(body string, conn net.Conn) error {
	header := "Crawl Dispatcher"
	length := 0

	content := pack(header, length, body)
	if _, err := conn.Write(content); err != nil {
		return err
	}

	return nil
}

func Parse(body string, conn net.Conn) error {
	header := "Parse Dispatcher"
	length := 0

	content := pack(header, length, body)
	if _, err := conn.Write(content); err != nil {
		return err
	}

	return nil
}

func Store(body string, conn net.Conn) error {
	header := "Parse Dispatcher"
	length := 0

	content := pack(header, length, body)
	if _, err := conn.Write(content); err != nil {
		return err
	}

	return nil
}



