package formats

import (
	"crypto/rand"
	"io"
)

func randomLetters(max int, table []byte) (string, error) {
	b := make([]byte, max)
	_, err := io.ReadAtLeast(rand.Reader, b, max)
	if err != nil {
		return "", io.ErrShortBuffer
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b), nil
}
