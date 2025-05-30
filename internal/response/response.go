package response

import(
	"strconv"
	"fmt"
	"io"
	"github.com/shreyasganesh0/TheStartup/internal/headers"
)

type StatusCode int

const (
	StatusOk StatusCode = 200
	StatusNotFound StatusCode = 400
	StatusServerError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statuscode StatusCode) error {


	var status_line string

	switch statuscode {

		case 200:

			status_line = "HTTP/1.1 200 OK\r\n";
		case 400:

			status_line = "HTTP/1.1 400 Not Found\r\n";
	
		case 500:

			status_line = "HTTP/1.1 500 Internal Server Error\r\n";

		default:

			return fmt.Errorf("Unsupported status code\n");
	}

	w.Write([]byte(status_line));
	return nil
} 

func GetDefaultHeaders(contentLen int) headers.Headers {

	v :=strconv.Itoa(contentLen);
	h := headers.Headers{
			"Content-Length":  v,
			"Connection": "closed",
			"Content-Type": "text/plain",
		}
	return h;
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {

	for k, v := range headers {

		w.Write([]byte(fmt.Sprintf("%s: %s\r\n", k, v)));
	}
	w.Write([]byte("\r\n"));
	return nil;
}
