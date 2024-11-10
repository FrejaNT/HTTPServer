package main

import (
	"fmt"
	"net"
	"net/http"
)

func SendBadRequestResponse(cn net.Conn, message string) {
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

func SendNotImplementedResponse(cn net.Conn) {
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

func SendOkResponse(cn net.Conn) {
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

func SendOkResponseWithBody(cn net.Conn, bd []byte) {
	res := http.Response{
		Status:        "200 OK",
		StatusCode:    http.StatusOK,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        make(http.Header),
		ContentLength: int64(len(bd)),
		Body:          nil,
	}

	res.Header.Set("Content-Type", "text/plain")

	if err := res.Write(cn); err != nil {
		fmt.Printf("error writing response: %s\n", err)
		return
	}
	cn.Write(bd)
}
