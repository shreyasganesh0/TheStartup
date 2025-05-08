package main

import( 
	"fmt"
	"os"
	"io"
	"bytes"
)

func main() {


	fd, err := os.Open("messages.txt");
	if (err != nil) {

		fmt.Printf("Error opening file %z\n", err);
	}

	chars := make([]byte, 8);
	tmpb := make([]byte, 0);
	updated_before_print := 1;

	for {
		_, err = fd.Read(chars);

		idx := bytes.IndexByte(chars, '\n');
		if (idx != -1) {

			tmpb = append(tmpb, chars[:idx]...);
			fmt.Printf("read: %s\n", string(tmpb));
			tmpb = tmpb[:0];
			tmpb = append(tmpb, chars[idx+1:]...);
			updated_before_print = 0;
		} else {

			tmpb = append(tmpb, chars...);
			updated_before_print = 1;
		}
		if (err == io.EOF) {

			break;
		}

	}
	if (updated_before_print == 1) {

		fmt.Printf("read: %s\n", string(tmpb));
	}

	return;

}
