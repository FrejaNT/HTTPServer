package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
)

func HandlePostRequest(cn net.Conn, bf *bufio.Reader) {
	for {
		ln, err := bf.ReadString('\n')
		if err != nil || ln == "\r\n" || ln == "\n" {
			break
		}
	}

	// create response
	res := http.Response{
		Status:        "200 OK",
		StatusCode:    http.StatusOK,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        make(http.Header),
		ContentLength: 0,
		Body:          nil,
	}

	res.Header.Set("Content-Type", "text/plain")

	if err := res.Write(cn); err != nil {
		fmt.Printf("error writing response: %s\n", err)
		return
	}
}
