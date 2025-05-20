package headers

import(
	"fmt"
	"bytes"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {

	return make(Headers);
}

func skip_ws(s *string) {

	for i, c := range *s {

		if (c == ' ' || c == '\n' || c == '\t') {

			continue;
		} else {

			*s = (*s)[i:];
			return;
		}
	}

	return;
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {

	done = false;

	idx := bytes.Index(data, []byte("\r\n"));
	if (idx == -1) {

		return n, done, err;
	}

	if (idx == 0) {

		done = true;
		return n, done, err;
	}

	tmp_data := data[:idx];

	s_data := string(tmp_data);

	var i int;
	var c rune;
	var curr_s string;
	var curr_key string;
	var b strings.Builder;
	for i, c = range s_data {

		if (c == ':') {

			break;
		}
		if (c == ' ') {

			err = fmt.Errorf("Error while parsing field name found whitespace\n");
			return n, done ,err;
		}
		b.WriteRune(c);
	}

	curr_s = b.String();
	h[curr_s] = "";
	curr_key = curr_s;
	curr_s = "";
	b.Reset();

	fmt.Printf("s data before skip ws is: %s\n",s_data);
	s_data = s_data[i + 2:];
	fmt.Printf("s data before skip ws is: %s\n",s_data);
	skip_ws(&s_data);

	for i, c = range s_data {

		if (c == ' ') {

			skip_ws(&s_data);
			if (len(s_data) > 0) {

				err = fmt.Errorf(" Failed to parse key");
				return n, done, err;
			} else {

				break;
			}
		}

		b.WriteRune(c);
	}

	curr_s = b.String();
	h[curr_key] = curr_s;
	n = idx + 2;

	return n, done, err;
}
