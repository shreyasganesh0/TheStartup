package headers

import(
	"fmt"
	"bytes"
)

type Headers map[string]string


func skip_ws(s *string) {

	for i, c := range s {

		if (c == ' ' || c == '\n' || c == '\t') {

			continue;
		} else {

			s = s[i + 1:];
			return;
		}
	}

	return;
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {

	var n int = 0;
	var done bool = false;
	var err error;

	idx := bytes.Index(data, []byte('\r\n'));
	if (idx == -1) {

		return n, done, err;
	}

	if (idx == 0) {

		done = true;
		return n, done, err;
	}

	s_data := string(data);

	var i int;
	var c rune;
	var curr_s string;
	for i, c = range s_data {

		if (c == ':') {

			break;
		}
		if (c == ' ') {

			err = fmt.Errorf("Error while parsing field name found whitespace\n");
			done = true;
			return n, done ,err;
		}
		curr_s += c;
	}

	h[curr_s] = "";
	curr_key = curr_s;
	curr_s = "";

	s_data = s_data[i + 1:];
	skip_ws(&s_data);

	for i, c = range s_data {

		if (c == ' ') {

			skip_ws(&s_data);
			if (s_data[i + 1] != '\r' || s_data[i + 2] != '\n') {

				err = fmt.Errorf(" Failed to parse key");
				done = true;
				return n, done, err;
			} else {

				done = true;
				n = len(data);
				break;
			}
		}

		curr_s += c;
	}

	h[curr_key] = curr_s;
	n = len(data);
	done = true;

	return n, done, err;
}
