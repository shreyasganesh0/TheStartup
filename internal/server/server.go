package server

import(
	"fmt"
	"sync/atomic"
	"net"

	"github.com/shreyasganesh0/TheStartup/internal/response"
	"github.com/shreyasganesh0/TheStartup/internal/request"
)

type Server struct {

	Listener net.Listener
	closed atomic.Bool
	handler Handler
}

type HandlerError struct {

	StatusCode response.StatusCode
	Message string
}
type Handler func(w *response.Writer, req *request.Request) 

func Serve(port int, h Handler) (*Server, error) {

	var server Server
	addr := fmt.Sprintf("127.0.0.1:%d", port);
	listener, err := net.Listen("tcp", addr)
	if (err != nil) {

		fmt.Printf("Error creating file %v\n", err);
	}

	server.Listener = listener;
	server.handler = h;
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

func (h *HandlerError) WriteHError(w *response.Writer) {

	fmt.Printf("got stat code %v",h.StatusCode);
	err := w.WriteStatusLine(h.StatusCode);
	if (err != nil) {

		fmt.Printf("Failed to write due to %s\n", err);
	}
	header := response.GetDefaultHeaders(len(h.Message));
	header.Update("Content-Type", "text/html")
	err = w.WriteHeaders(header);
	if (err != nil) {

		fmt.Printf("Failed to write due to %s\n", err);
	}
	w.Writer.Write([]byte(h.Message));
	fmt.Printf("Wrote %v",h.Message);
	return;
}

func (s *Server) handle(conn net.Conn) {

	defer conn.Close();

	var w response.Writer
	w.Writer = conn

	req, err_req := request.RequestFromReader(conn);
	if (err_req != nil) {

		fmt.Printf("Error found is %v\n", err_req);
		h_e := HandlerError {

			Message: "Your problem is not my problem\n",
			StatusCode:  400,
		}
		h_e.WriteHError(&w);
		return;
	}

	s.handler(&w, req);
	return;
}
