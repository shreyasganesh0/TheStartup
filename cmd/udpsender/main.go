package main

import (
	"os"
	"fmt"
	"net"
	"bufio"
)

func main() {

	udp_addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:42069");
	if (err != nil) {

		fmt.Printf("Error resolving udp name %z", err);
	}

	udp_conn, errudp := net.DialUDP("udp", nil, udp_addr)
	if (errudp != nil) {

		fmt.Printf("Error resolving udp name %z", errudp);
	}


	buf_reader := bufio.NewReader(os.Stdin);

	for {

		fmt.Print(">");
		s, errs := buf_reader.ReadString('\n');
		if (errs != nil) {

			fmt.Printf("Error resolving udp name %z", errudp);
		}

		s_b := []byte(s);
		_, err_w := udp_conn.Write(s_b);

		if (err_w != nil) {

			fmt.Printf("Failed to write bytes %z\n", err_w);
		}
	}
}




