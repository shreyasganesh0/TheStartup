package main

import(
	"log"
	"os/signal"
	"syscall"
	"os"
	"io"
	"github.com/shreyasganesh0/TheStartup/internal/server"
	"github.com/shreyasganesh0/TheStartup/internal/request"
)


func myhandler(w io.Writer, req *request.Request) *server.HandlerError {

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

	if (req.RequestLine.RequestTarget == "/yourproblem") {

		return &server.HandlerError {
			Message: badreq_message, 
			StatusCode: 400,
		}
	} 
	if (req.RequestLine.RequestTarget == "/myproblem") {

		return &server.HandlerError {
			Message: serverr_message, 
			StatusCode: 500,
		}
	} 

	w.Write([]byte(ok_message));
	return nil;
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
