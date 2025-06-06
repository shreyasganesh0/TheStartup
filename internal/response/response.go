package response

import(
	"strconv"
	"fmt"
	"net"
	"github.com/shreyasganesh0/TheStartup/internal/headers"
)

type StatusCode int

const (
	StatusOk StatusCode = 200
	StatusNotFound StatusCode = 400
	StatusServerError StatusCode = 500
)

type Writer struct {

	Writer net.Conn
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {

	var n int
	var err error
	var n_tot int

	n, err = fmt.Fprintf(w.Writer, "%X\r\n", len(p));
	if err != nil {

		return n, err
	}
	n_tot += n;
	fmt.Printf("Writing %X\n", len(p));

	n, err = w.Writer.Write(p);
	if err != nil {

		return n_tot, err
	}
	n_tot += n;
	fmt.Printf("Writing %q\n", p);

	n, err = w.Writer.Write([]byte("\r\n"));
	if err != nil {

		return n_tot, err
	}
	n_tot += n;
	fmt.Printf("Writing \r\n");

	return n_tot, nil;
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {

	var n int
	var err error
	var n_tot int

	n, err = w.Writer.Write([]byte("0\r\n\r\n"));
	if err != nil {

		return n_tot, err
	}
	n_tot += n;
	fmt.Printf("Writing 0\r\n\r\n");

	return n_tot, nil;
}

func (w *Writer) WriteStatusLine( statuscode StatusCode) error {


	var status_line string

	switch statuscode {

		case 200:

			status_line = "HTTP/1.1 200 OK\r\n";
		case 400:

			status_line = "HTTP/1.1 400 Bad Request\r\n";
		case 500:

			status_line = "HTTP/1.1 500 Internal Server Error\r\n";
		default:

			return fmt.Errorf("Unsupported status code\n");
	}

	w.Writer.Write([]byte(status_line));
	fmt.Printf("Wrote %v",status_line);
	return nil
} 

func GetDefaultHeaders(contentLen int) headers.Headers {

	v :=strconv.Itoa(contentLen);
	h := headers.Headers{
			"Content-Length":  v,
			"Connection": "close",
			"Content-Type": "text/plain",
		}
	return h;
}

func (w *Writer) WriteHeaders( headers headers.Headers) error {

	for k, v := range headers {

		w.Writer.Write([]byte(fmt.Sprintf("%s: %s\r\n", k, v)));
		fmt.Printf("Wrote %v", fmt.Sprintf("%s: %s\r\n", k, v));
	}
	w.Writer.Write([]byte("\r\n"));
	fmt.Printf("Wrote %v", "\r\n");
	return nil;
}


func (w *Writer) WriteBody(p []byte) error {

	w.Writer.Write(p);
	return nil;
}
