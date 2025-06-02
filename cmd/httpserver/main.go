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

	var h_err server.HandlerError;

	if (req.RequestLine.RequestTarget == "/yourproblem") {

		h_err.Message = "Your problem is not my problem\n";
		h_err.StatusCode = 400;
	} else if (req.RequestLine.RequestTarget == "/myproblem") {

		h_err.Message = "Woopsie, my bad\n";
		h_err.StatusCode = 500;
	} else {

		w.Write([]byte("All good, frfr\n"));
		return nil;
	}

	return &h_err;
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
