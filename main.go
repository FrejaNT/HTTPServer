package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"
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
		ch <- 1
		go func() {
			defer func() { <-ch }()
			parseRequest(cn)
		}()
	}
}

func parseRequest(cn net.Conn) {
	defer cn.Close()
	// parse request

	bf := bufio.NewReader(cn)

	// Parse the HTTP request using http.ReadRequest
	rq, er := http.ReadRequest(bf)
	if er != nil {
		// If there's an error reading the request, send a 400 Bad Request response
		sendBadRequestResponse(cn, "Bad Request")
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

func checkRequestType(rq *http.Request) RequestType {
	method, proto := rq.Method, rq.Proto
	if proto != "HTTP/1.0" && proto != "HTTP/1.1" {
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
