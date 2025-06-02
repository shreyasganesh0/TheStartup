package request

import (
	"io"
	"fmt"
	"bytes"
	"bufio"
	"strconv"
	"github.com/shreyasganesh0/TheStartup/internal/headers"
)

const READ_BYTES int = 8;

type Request struct {
	RequestLine RequestLine
	Headers headers.Headers
	Body []byte
	state int
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func (r *Request) parse(data []byte) (int, error) {

	var req RequestLine;
	var err error;
	var n int;
	var done bool;

	if r.state == 0 {

		req, n, err = parseRequestLine(data);
		if n > 0 {

			r.RequestLine = req;
			r.state = 2; //switch to header
		}
	} else if r.state == 1 {

		fmt.Printf("ended\n");
		err = fmt.Errorf("Trying to read in done state\n");

	} else if r.state == 2 {
		fmt.Printf("state2 data is : %s\n",data );

		n, done, err = r.Headers.Parse(data)
		fmt.Printf("n: %v done: %v err: %v\n",n, done, err);
		if (done == true) {

			_, exists := r.Headers["content-length"]
			if (exists == false) {
				
				r.state = 1;
			} else {

				r.state = 3;
			}
		}

	} else if r.state == 3 {

		v, _ := r.Headers["content-length"];


		n, err_conv := strconv.Atoi(v);
		if err_conv != nil {

			err = err_conv;
		}

		if (len(data) == n) {

			r.Body = append(r.Body, data...);
			
		} else if (len(data) > n){

			n = 0
			err = fmt.Errorf("Length %d and content length header %s dont match\n", len(r.Body), v);
		} 
		r.state = 1


	} else {
		
		err = fmt.Errorf("Unknown state\n");
	}

	return n, err

}

func parseRequestLine(byts []byte) (RequestLine, int, error) {

	var req RequestLine;
	var err error;
	var num_bytes int;
	var n int;
	
	idx := bytes.Index(byts, []byte("\r\n"));
	if (idx == -1) {

		err = fmt.Errorf("invalid packet %s\n", string(byts));
		return req, 0, nil
	}
	
	byts = byts[:idx];
	num_bytes = len(byts) + 2;
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

	return req, num_bytes, err

}

func RequestFromReader(reader io.Reader) (*Request, error) {
	
	var err error;
	var n int;
	var r Request;

	r.Headers = headers.NewHeaders()
	byts := make([]byte, READ_BYTES, READ_BYTES);
	send_byts := make([]byte, 0);
	buf := bufio.NewReader(reader);

	for (r.state != 1) {

		n, err = buf.Read(byts);	
		if (err == io.EOF) {

			r.state = 1
			break;
		}
		fmt.Printf("Checking...\n");
		if (n > 0) {

			send_byts = append(send_byts, byts[:n]...);
			fmt.Printf("Raw buffer contents: %q\n", send_byts)
		    n_sub, err_sub := r.parse(send_byts);
			if err_sub != nil {

				fmt.Printf("Error  parsing is %v", err_sub);
				err = err_sub;
				break;
			}

			if n_sub > 0 {

				send_byts = send_byts[n_sub:]
				if bytes.Equal(send_byts, []byte("\r\n")) {

					r.parse(send_byts);
				}

			}
		}
	}


	return &r, err

}
