package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)


func main() {
	address := "localhost:5000"
	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		log.Fatalln(err)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatalln(err)
	}
	msg := "GET /bb HTTP/1.0\r\nHost: 127.0.0.1:5000\r\nUser-Agent: Mozilla/4.0 (compatible; MSIE5.01; Windows NT)\r\n\r\n"
	bytes := []byte(msg)

	conn.Write(bytes)

	//buf := make([]byte, 2000)
	//read, _ := conn.Read(buf)
	//fmt.Println(string(buf[:read]))
	//buf = make([]byte, 2000)
	//read1, _ := conn.Read(buf)
	//fmt.Println(string(buf[:read1]))


	reader := bufio.NewReader(conn)
	buf1 := make([]byte, 10)

	for i := 1; i < 20; i++ {
		buf1, _ = reader.ReadBytes('\n')
		if len(buf1) <= 0 {
			break
		}
		fmt.Println(len(buf1))
		fmt.Println(string(buf1))
	}
}
