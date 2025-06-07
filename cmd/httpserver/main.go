package main

import(
	"log"
	"os/signal"
	"syscall"
	"os"
	"io"
	"strings"
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"crypto/sha256"

	"github.com/shreyasganesh0/TheStartup/internal/server"
	"github.com/shreyasganesh0/TheStartup/internal/request"
	"github.com/shreyasganesh0/TheStartup/internal/response"
)

func proxyHandle(url string, w *response.Writer) {

	var statuscode response.StatusCode = 200
	var x_buf bytes.Buffer

	resp, err := http.Get(url);
	fmt.Printf("Url is %s", url);
	if (err != nil) {

		h_err := &server.HandlerError {
			Message: err.Error(), 
			StatusCode: 400,
		}
		h_err.WriteHError(w);
		return;
	}
	defer resp.Body.Close();
	
	buf := make([]byte,1024);
	err = w.WriteStatusLine(statuscode);
	if (err != nil) {

		fmt.Printf("Failed to write due to %s\n", err);
	}
	h := response.GetDefaultHeaders(len(buf));
	h.Update("Transfer-Encoding", "chunked");
	h.Update("Trailers", "X-Content-SHA256, X-Content-Length");
	h.Remove("Content-Length");

	err = w.WriteHeaders(h);
	if (err != nil) {

		fmt.Printf("Failed to write due to %s\n", err);
	}
	for {

		n, err_rd := resp.Body.Read(buf);

		if (err_rd == io.EOF) {

			_, errw := w.WriteChunkedBodyDone();
			if (errw != nil) {

				fmt.Printf("Failed to write due to %v\n", errw);
			}

			h.Update("X-Content-SHA256", fmt.Sprintf("%x", sha256.Sum256(x_buf.Bytes())));
			h.Update("X-Content-Length",strconv.Itoa(x_buf.Len()));

			err_w := w.WriteTrailers(h);
			if (err_w != nil) {

				fmt.Printf("Failed to write due to %v\n", err_w);
			}

			break;
		}
		if (err_rd != nil) {

			fmt.Printf("Error found while reading %v", err_rd);
		}

		if (n > 0) {

			_, err_w := w.WriteChunkedBody(buf[:n]);
			x_buf.Write(buf[:n]);

			if (err_w != nil) {
				fmt.Printf("Error found while reading %v", err_w);
			}
		}
	}

	fmt.Printf("Done writing\n");
}

func myhandler(w *response.Writer, req *request.Request) {

	var badreq_message string = `<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`

	var serverr_message string = `<html>
 <head>
	<title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`

	var ok_message string = `<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`

	var h_err *server.HandlerError;

	if (strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin")) {

		url := "https://httpbin.org" + strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin");
		proxyHandle(url, w);
		return;
	} else if (req.RequestLine.RequestTarget == "/yourproblem") {

		h_err = &server.HandlerError {
				Message: badreq_message, 
				StatusCode: 400,
			}
		h_err.WriteHError(w);
		return;
	} else if (req.RequestLine.RequestTarget == "/myproblem") {

		h_err = &server.HandlerError {
			Message: serverr_message, 
			StatusCode: 500,
		}
		h_err.WriteHError(w);
		return;
	} 

//default message	
	var statuscode response.StatusCode = 200
	err := w.WriteStatusLine(statuscode);
	if (err != nil) {

		fmt.Printf("Failed to write due to %s\n", err);
	}
	h := response.GetDefaultHeaders(len(ok_message));
	h.Update("Content-Type", "text/html")

	err = w.WriteHeaders(h);
	if (err != nil) {

		fmt.Printf("Failed to write due to %s\n", err);
	}

	w.Writer.Write([]byte(ok_message));
	return; 
}

func main() {
	port := 42069
	server, err := server.Serve(port, myhandler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
