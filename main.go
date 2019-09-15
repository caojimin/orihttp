package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
)

var addr string

func init() {
	flag.StringVar(&addr, "addr", "127.0.0.1:9300", "the http server address")
	flag.Parse()
}

func handleConnection(conn net.Conn) {
	defer func() {
		_ = conn.Close()
	}()
	br := bufio.NewReader(conn)
	reader := make([]byte, br.Size())
	resp := []string{"HTTP/1.1 200 OK"}
	req := ""
	for {
		n, err := br.Read(reader)
		if err != nil || n == 0 {
			break
		}
		req += fmt.Sprintf("%s", reader)[:n]
		reader = reader[:0]
	}
	resp = append(resp, fmt.Sprintf("Content-Length: %d", len(req)))
	resp = append(resp, "")
	resp = append(resp, req)
	_, err := fmt.Fprintf(conn, strings.Join(resp, "\n"))
	if err != nil {
		log.Printf("write to client error: %s\n", err)
	}
}
func main() {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = listener.Close()
	}()
	log.Printf("the server is running on %s!\n", addr)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("get client connection error: %s\n", err)
		} else {
			go handleConnection(conn)
		}
	}

}
