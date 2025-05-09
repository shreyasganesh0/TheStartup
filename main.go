package main

import( 
	"fmt"
	"os"
	"io"
	"bytes"
)

func getLinesChannel(f io.ReadCloser) <- chan string {

	ch := make(chan string);
	go func() {

		defer close(ch);
		chars := make([]byte, 8);
		tmpb := make([]byte, 0);
		updated_before_print := 1;

		for {
			n, err := f.Read(chars);

			if (n > 0) {

				idx := bytes.IndexByte(chars, '\n');
				if (idx != -1) {

					tmpb = append(tmpb, chars[:idx]...);
					ch <- string(tmpb)
					tmpb = tmpb[:0];
					tmpb = append(tmpb, chars[idx+1:]...);
					updated_before_print = 0;
				} else {

					tmpb = append(tmpb, chars...);
					updated_before_print = 1;
				}
			}

			if (err == io.EOF) {
				break;
			}

		}
		if (updated_before_print == 1) {

			ch <- string(tmpb);
		}
	}()

	return ch;
}

func main() {


	fd, err := os.Open("messages.txt");
	defer fd.Close();
	if (err != nil) {

		fmt.Printf("Error opening file %z\n", err);
	}

	ch := getLinesChannel(fd);

	for msg := range ch {

		fmt.Printf("read: %s\n", msg);
	}


	return;

}
