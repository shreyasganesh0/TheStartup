package headers

import(
	"testing"
	"github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	headers = NewHeaders()
	data = []byte("Set-Person: lane-loves-go\r\n");

	n, done, err = headers.Parse(data)
	require.NotNil(t, headers)
	assert.Equal(t, "lane-loves-go", headers["set-person"])

	data = []byte("Set-Person: prime-loves-zig\r\n");
	n, done, err = headers.Parse(data)
	assert.Equal(t, "lane-loves-go, prime-loves-zig", headers["set-person"])

	data = []byte("Set-Person: tj-loves-ocaml\r\n");
	n, done, err = headers.Parse(data)
	assert.Equal(t, "lane-loves-go, prime-loves-zig, tj-loves-ocaml", headers["set-person"])
}
