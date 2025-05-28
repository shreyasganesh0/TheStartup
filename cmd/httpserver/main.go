package main

import(
	"net"
	"fmt"
	"log"
	"sync/atomic"
	"os/signal"
	"syscall"
	"os"
)

type Server struct {

	Listener net.Listener
	closed atomic.Bool
}

func Serve(port int) (*Server, error) {

	var server Server
	addr := fmt.Sprintf("127.0.0.1:%d", port);
	listener, err := net.Listen("tcp", addr)
	if (err != nil) {

		fmt.Printf("Error creating file %v\n", err);
	}

	server.Listener = listener;
	go server.listen();
	return &server, err;
}

func (s *Server) Close() error {

	s.closed.Store(true)
	if s.Listener != nil {
		return s.Listener.Close();
	}

	return nil;

}
func (s *Server) listen() {
	for {

		tcp_conn, err_tcp := s.Listener.Accept()
		if (err_tcp != nil) {
			if s.closed.Load() {
				return //skip handled ones
			}

			fmt.Printf("Error creating file %z\n", err_tcp);
			continue
		}
		go s.handle(tcp_conn)
	}
}
func (s *Server) handle(conn net.Conn) {

	defer conn.Close()
	response := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nHello World!"
	conn.Write([]byte(response));
	return;
}

func main() {
	port := 42069
	server, err := Serve(port)
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
