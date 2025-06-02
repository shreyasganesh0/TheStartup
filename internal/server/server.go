package server

import(
	"fmt"
	"sync/atomic"
	"bytes"
	"net"
	"io"
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
type Handler func(w io.Writer, req *request.Request) *HandlerError

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

func WriteHError(h *HandlerError, w io.Writer) {

	err := response.WriteStatusLine(w, h.StatusCode);
	if (err != nil) {

		fmt.Printf("Failed to write due to %s\n", err);
	}
	header := response.GetDefaultHeaders(len(h.Message));
	err = response.WriteHeaders(w, header);
	if (err != nil) {

		fmt.Printf("Failed to write due to %s\n", err);
	}
	w.Write([]byte(h.Message));
	fmt.Printf("Wrote %v\n",h.Message);
	return;
}

func (s *Server) handle(conn net.Conn) {

	defer conn.Close();
	var statuscode response.StatusCode = 200
	var buf bytes.Buffer

	req, err_req := request.RequestFromReader(conn);
	if (err_req != nil) {

		h_e := HandlerError {

			Message: "Your problem is not my problem\n",
			StatusCode:  400,
		}
		WriteHError(&h_e, conn);
		return;
	}

	h_err := s.handler(&buf, req);
	if (h_err != nil) {

		WriteHError(h_err, conn);
		return;
	}

	err := response.WriteStatusLine(conn, statuscode);
	if (err != nil) {

		fmt.Printf("Failed to write due to %s\n", err);
	}
	h := response.GetDefaultHeaders(buf.Len());
	err = response.WriteHeaders(conn, h);
	if (err != nil) {

		fmt.Printf("Failed to write due to %s\n", err);
	}

	conn.Write(buf.Bytes());
	return;
}
