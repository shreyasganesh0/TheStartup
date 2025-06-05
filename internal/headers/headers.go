package headers

import(
	"fmt"
	"bytes"
	"strings"
	"unicode"
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

func (h Headers) Remove(k string) {

	delete(h, k);
	return;
}

func (h Headers) Update(k, v string) {

	h[k] = v;
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
		n = 2
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
		if (c == ' '|| !(unicode.IsUpper(c) || unicode.IsLower(c) || unicode.IsDigit(c) || c == '!' || c == '#' || c == '$' || c == '%' || c == '&' || c == '\'' || c == '*'|| c == '+' || c == '-' || c == '.' || c == '^' || c == '_' ||  c == '`' || c == '|' || c == '~') ) {

			err = fmt.Errorf("Error while parsing field name found whitespace\n");
			return n, done ,err;
		} 

		c = unicode.ToLower(c);

		b.WriteRune(c);
	}

	curr_s = b.String();

	curr_key = curr_s;
	curr_s = "";
	b.Reset();

	s_data = s_data[i + 2:];
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
	val, exists := h[curr_key];
	if (exists == false) {

		h[curr_key] = curr_s;
	} else {

		h[curr_key] = val + ", " + curr_s;
	}
	n = idx+2; // len of consumed bytes +2 is for the crlf

	return n, done, err;
}
