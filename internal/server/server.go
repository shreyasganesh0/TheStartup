package server

import(
	"fmt"
	"sync/atomic"
	"net"
	"github.com/shreyasganesh0/TheStartup/internal/response"
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

	defer conn.Close();
	var statuscode response.StatusCode = 200

	err := response.WriteStatusLine(conn, statuscode);
	if (err != nil) {

		fmt.Printf("Failed to write due to %s\n", err);
	}
	h := response.GetDefaultHeaders(0);
	err = response.WriteHeaders(conn, h);
	if (err != nil) {

		fmt.Printf("Failed to write due to %s\n", err);
	}
	return;
}
