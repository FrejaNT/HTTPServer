package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
)

type RequestType int64

const (
	BAD  RequestType = 0
	NYI  RequestType = 1
	GET  RequestType = 2
	POST RequestType = 3
)

const maxRoutines = 10

func main() {
	var wg sync.WaitGroup
	ch := make(chan int, maxRoutines)

	//listening
	ls, er := net.Listen("tcp", "localhost:3333")
	if er != nil {
		fmt.Printf("error starting server: %s\n", er)
		os.Exit(1)
	}
	defer ls.Close()

	// checking host, port etc
	host, port, err := net.SplitHostPort(ls.Addr().String())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Listening on host: %s, port: %s\n", host, port)

	// connections
	for {
		cn, er := ls.Accept()
		if er != nil {
			fmt.Printf("error starting connection: %s\n", er)
		}
		wg.Add(1)
		ch <- 1
		go func() {
			defer func() { wg.Done(); <-ch }()
			parseRequest(cn)
		}()
	}
}

func parseRequest(cn net.Conn) {
	defer cn.Close()
	// parse request

	bf := bufio.NewReader(cn)
	rq, er := bf.ReadString('\n')
	if er != nil {
		fmt.Printf("error writing response: %s\n", er)
		return
	}
	rType := checkRequestType(rq)

	switch rType {
	case GET:
		HandleGetRequest(cn, bf)
	case POST:
		HandlePostRequest(cn, bf)
	case NYI:
		sendNotImplementedResponse(cn)
	case BAD:
		sendBadRequestResponse(cn, "goof")
	}
}

func checkRequestType(rq string) RequestType {
	parts := strings.Fields(rq)
	if len(parts) != 3 {
		return BAD
	}

	method, path, proto := parts[0], parts[1], parts[2]

	if !strings.HasPrefix(path, "/") || proto != "HTTP/1.0" && proto != "HTTP/1.1" {
		return BAD
	}

	switch method {
	case "GET":
		return GET
	case "POST":
		return POST
	case "PUT":
		return NYI
	case "DELETE":
		return NYI
	case "PATCH":
		return NYI
	}

	return BAD
}

func sendBadRequestResponse(cn net.Conn, message string) {
	res := http.Response{
		Status:        "400 Bad Request",
		StatusCode:    http.StatusBadRequest,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        make(http.Header),
		ContentLength: int64(len(message)),
		Body:          nil,
	}

	if err := res.Write(cn); err != nil {
		fmt.Printf("error writing bad request response: %s\n", err)
		return
	}
	cn.Write([]byte(message))
}

func sendNotImplementedResponse(cn net.Conn) {
	res := http.Response{
		Status:     "501 Not Implemented",
		StatusCode: http.StatusNotImplemented,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
	}

	if err := res.Write(cn); err != nil {
		fmt.Printf("error writing bad request response: %s\n", err)
		return
	}
}

/* 	// Wrap the connection in a bufio.Reader
reader := bufio.NewReader(c)

// Parse the HTTP request using http.ReadRequest
req, err := http.ReadRequest(reader)
if err != nil {
	// If there's an error reading the request, send a 400 Bad Request response
	sendBadRequestResponse(c)
	return */
