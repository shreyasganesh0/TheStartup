package request

import (
	"io"
	"fmt"
	"bytes"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func parseRequestLine(byts []byte) (RequestLine, error) {

	var req RequestLine;
	var err error;
	var n int;
	
	idx := bytes.Index(byts, []byte("\r\n"));
	if (idx == -1) {

		err = fmt.Errorf("invalid packet %s\n", string(byts));
	}
	
	byts = byts[:idx]
	tmps := string(byts);
	
	n, err = fmt.Sscanf(tmps, "%s %s HTTP/%s", &req.Method, &req.RequestTarget, &req.HttpVersion);
	if (n < 3) {

		err = fmt.Errorf("Failed to parse line\n");
	}


	if (req.Method != "GET" && req.Method != "POST" && req.Method != "PUT" && req.Method != "DELETE") {

		err = fmt.Errorf("Invalid method parsed\n");
	}

	if (req.HttpVersion != "1.1") {

		err = fmt.Errorf("Invalid version: %s %s %s\n", req.Method, req.RequestTarget, req.HttpVersion);
	}

	return req, err

}
func RequestFromReader(reader io.Reader) (*Request, error) {
	
	byts, err := io.ReadAll(reader);	
	if (err != nil) {

		fmt.Printf("Failed to read bytes from reader %z\n", err);
	}

	req, err := parseRequestLine(byts);

	reqs := Request {RequestLine: req};
	return &reqs, err

}
