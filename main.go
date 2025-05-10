package main

import( 
	"fmt"
	"io"
//	"log"
	"net"
	"bytes"
)

func getLinesChannel(f io.ReadCloser) <- chan string {

	ch := make(chan string);
	go func() {

		defer close(ch);
		chars := make([]byte, 8);
		tmpb := make([]byte, 0);

		for {
			n, err := f.Read(chars);

			if (n > 0) {

				idx := bytes.IndexByte(chars, '\n');
				chunk := chars[:n]
				if (idx != -1) {

					tmpb = append(tmpb, chunk[:idx]...);
					ch <- string(tmpb)
					tmpb = tmpb[:0];
					tmpb = append(tmpb, chunk[idx+1:]...);
				} else {

					tmpb = append(tmpb, chunk...);
				}
			}

			if (err == io.EOF) {
				break;
			}

		}
		if len(tmpb) > 0 {

			ch <- string(tmpb);
		}
	}()

	return ch;
}

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

		ch := getLinesChannel(tcp_conn);
		for msg := range ch {

			fmt.Printf("%s\n", msg);
		}
		//log.Printf("closing connection %z\n", tcp_conn);
	}


	return;

}
