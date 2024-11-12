package main

import (
	"bufio"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"strings"
)

const maxRoutines = 10

// TODO: concurrency, errors and i'm lost also ports and stuff

/*
Questions: concurrency, errors, queries/params, HTTP versions, content-types, file name(key thing),

	sending files with curl, what are we allowed to use in net/http
*/
func main() {
	ch := make(chan int, maxRoutines)

	args := os.Args[1:]

	//listening
	ls, er := net.Listen("tcp", "localhost:"+args[0])
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
		cn, err := ls.Accept()
		if err != nil {
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
	rq, err := http.ReadRequest(bf)

	if err != nil {
		// If there's an error reading the request, send a 400 Bad Request response
		SendBadRequestResponse(cn, "Error parsing HTTP request")
		return
	}

	md, pr := rq.Method, rq.Proto

	if pr != "HTTP/1.0" && pr != "HTTP/1.1" {
		SendBadRequestResponse(cn, "Invalid HTTP version")
		return
	}

	switch md {
	case "GET":
		bd, err := handleGetRequest(rq)
		if err != nil {
			SendBadRequestResponse(cn, err.Error())
			return
		}
		SendOkResponseWithBody(cn, bd)

	case "POST":
		if err := handlePostRequest(rq); err != nil {
			SendBadRequestResponse(cn, err.Error())
			return
		}
		SendOkResponse(cn)

	case "PUT", "DELETE", "PATCH":
		SendNotImplementedResponse(cn)

	default:
		SendBadRequestResponse(cn, "Invalid HTTP method")
	}
}

// TODO: errors etc
func handleGetRequest(rq *http.Request) ([]byte, error) {
	fn, _ := strings.CutPrefix(rq.URL.Path, "/")

	bd, err := os.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	// create response
	return bd, nil
}

func handlePostRequest(rq *http.Request) error {
	if rq.Header.Get("Content-Type") == "multipart/form-data" {
		return handlePostRequestMultiForm(rq)
	}
	if err := checkContentType(rq.Header.Get("Content-Type")); err != nil {
		return err
	}
	if err := saveToFile(rq.Body, rq.Header.Get("filename")); err != nil {
		return err
	}

	return nil
}

// TODO: probably not good
func handlePostRequestMultiForm(rq *http.Request) error {
	rq.ParseMultipartForm(10 << 20)
	var key string
	for key_ := range rq.MultipartForm.File { //ugly
		key = key_
	}

	fi, ha, err := rq.FormFile(key)
	if err != nil {
		return err
	}
	defer fi.Close()

	if err := checkContentType(ha.Header.Get("Content-Type")); err != nil {
		return err
	}

	if err := saveToFileMultiForm(fi, ha.Filename); err != nil {
		return err
	}

	return nil
}

func checkContentType(ct string) error {
	switch ct {
	case "text/html", "text/plain", "image/gif", "image/jpeg", "image/jpg", "text/css":
		return nil
	default:
		return fmt.Errorf("invalid content type")
	}

}

// Function to save data to a file
func saveToFileMultiForm(fi multipart.File, fn string) error {
	ds, err := os.Create(fn)
	if err != nil {
		return err
	}
	defer ds.Close()
	if _, err := io.Copy(ds, fi); err != nil {
		return err
	}
	return nil
}
func saveToFile(bd io.ReadCloser, fn string) error {
	ds, err := os.Create(fn)
	if err != nil {
		return err
	}
	defer ds.Close()
	if _, err := io.Copy(ds, bd); err != nil {
		return err
	}
	return nil
}
