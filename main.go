package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/http"
	"os"
	"strings"
)

// max number of routines allowed
const maxRoutines = 10

func main() {

	args := os.Args[1:]

	// second argument is ip if it exists otherwise listens on all interfaces
	var ip string
	if len(args) > 1 {
		ip = args[1]
	} else {
		ip = "0.0.0.0"
	}

	//listening
	ls, er := net.Listen("tcp", ip+":"+args[0])
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

	// Channel for limiting number of routines
	ch := make(chan int, maxRoutines)

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
	rq, err := http.ReadRequest(bf)

	if err != nil {
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
			if errors.Is(err, fs.ErrNotExist) {
				SendFileNotFoundResponse(cn)
				return
			}
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

func handleGetRequest(rq *http.Request) ([]byte, error) {
	fn, err := getURLFileName(rq)
	if err != nil {
		return nil, err
	}

	bd, err := os.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	return bd, nil
}

func handlePostRequest(rq *http.Request) error {
	fn, err := getURLFileName(rq)

	if err != nil {
		return err
	}

	if strings.HasPrefix(rq.Header.Get("Content-Type"), "multipart/form-data") {
		return handlePostRequestMultiForm(rq, fn)
	}

	// if there is a valid destination url
	if fn != "" {
		if err := saveToFile(rq.Body, fn); err != nil {
			return err
		}
		return nil
	}

	// otherwise check headers for content-type/name
	if err := checkContentType(rq.Header.Get("Content-Type")); err != nil {
		return err
	}

	if err := saveToFile(rq.Body, rq.Header.Get("file")); err != nil {
		return err
	}

	return nil
}

func handlePostRequestMultiForm(rq *http.Request, fname string) error {
	rq.ParseMultipartForm(10 << 20)

	//gives http: no such file when "file" key not used in request
	//maybe give more explicit error?
	fi, hd, err := rq.FormFile("file")
	if err != nil {
		return err
	}
	defer fi.Close()

	if err := checkContentType(hd.Header.Get("Content-Type")); err != nil {
		return err
	}

	var fn string
	if fname != "" {
		fn = fname
	} else {
		fn = hd.Filename
	}

	if err := saveToFile(fi, fn); err != nil {
		return err
	}

	return nil
}

func getURLFileName(rq *http.Request) (string, error) {
	fn, _ := strings.CutPrefix(rq.URL.Path, "/")
	if fn == "" {
		return "", nil
	}
	ex := strings.Split(fn, ".")
	if len(ex) != 2 {
		return "", fmt.Errorf("invalid URL")
	}
	if err := checkExtension(ex[1]); err != nil {
		return "", err
	}
	return fn, nil
}

func checkExtension(ex string) error {
	switch ex {
	case "html", "txt", "gif", "jpeg", "jpg", "css":
		return nil
	default:
		return fmt.Errorf("invalid content type")
	}
}

func checkContentType(ct string) error {
	switch ct {
	case "text/html", "text/plain", "image/gif", "image/jpeg", "image/jpg", "text/css":
		return nil
	default:
		return fmt.Errorf("invalid content type")
	}
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
