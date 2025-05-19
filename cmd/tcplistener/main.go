package main

import( 
	"fmt"
	"net"
	"github.com/shreyasganesh0/TheStartup/request"
)

func main() {

	listener, err := net.Listen("tcp4", "127.0.0.1:42069")
	if (err != nil) {

		fmt.Printf("Error creating file %z\n", err);
	}

	defer listener.Close()
	for {

		tcp_conn, err_tcp := listener.Accept()
		if (err_tcp != nil) {

			fmt.Printf("Error creating file %z\n", err_tcp);
		}

		//log.Printf("created connection %z\n", tcp_conn);


		req, err_req := request.RequestFromReader(tcp_conn);
		if (err_req != nil) {

			fmt.Printf("Error creating file %z\n", err_tcp);
		}

		fmt.Printf("Request ine:\n");
		fmt.Printf("- Method: %s\n", req.RequestLine.Method);
		fmt.Printf("- Target: %s\n", req.RequestLine.RequestTarget);
		fmt.Printf("- Version: %s\n", req.RequestLine.HttpVersion);

	}


	return;

}
