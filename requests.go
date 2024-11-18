package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

func SendBadRequestResponse(cn net.Conn, message string) {
	body := io.NopCloser(strings.NewReader(message))
	rs := http.Response{
		Status:        "400 Bad Request",
		StatusCode:    http.StatusBadRequest,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        make(http.Header),
		ContentLength: int64(len(message)),
		Body:          body,
	}

	rs.Header.Add("Connection", "close")
	rs.Header.Add("Content-Type", "text/plain")

	if err := rs.Write(cn); err != nil {
		fmt.Printf("error writing bad request response: %s\n", err)
	}

	if tcpConn, ok := cn.(*net.TCPConn); ok {
		tcpConn.CloseWrite()
	}
}

func SendNotImplementedResponse(cn net.Conn) {
	rs := http.Response{
		Status:     "501 Not Implemented",
		StatusCode: http.StatusNotImplemented,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
	}

	rs.Header.Add("Connection", "close")

	if err := rs.Write(cn); err != nil {
		fmt.Printf("error writing bad request response: %s\n", err)
	}
}

func SendOkResponse(cn net.Conn) {
	rs := http.Response{
		Status:        "200 OK",
		StatusCode:    http.StatusOK,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        make(http.Header),
		ContentLength: 0,
		Body:          nil,
	}

	rs.Header.Add("Connection", "close")

	if err := rs.Write(cn); err != nil {
		fmt.Printf("error writing response: %s\n", err)
		return
	}
}

func SendOkResponseWithBody(cn net.Conn, bd []byte) {
	body := io.NopCloser(bytes.NewReader(bd))
	rs := http.Response{
		Status:        "200 OK",
		StatusCode:    http.StatusOK,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        make(http.Header),
		ContentLength: int64(len(bd)),
		Body:          body,
	}

	rs.Header.Add("Connection", "close")

	if err := rs.Write(cn); err != nil {
		fmt.Printf("error writing response: %s\n", err)
		return
	}
}
func SendFileNotFoundResponse(cn net.Conn) {
	rs := http.Response{
		Status:        "404 Not Found",
		StatusCode:    http.StatusNotFound,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        make(http.Header),
		ContentLength: 0,
		Body:          nil,
	}

	rs.Header.Add("Connection", "close")

	if err := rs.Write(cn); err != nil {
		fmt.Printf("error writing response: %s\n", err)
		return
	}
}
