package formats_test

import (
	"encoding/hex"
	"strings"
	"testing"
)

// block decodes a hex PIN block used as a test vector.
func block(t testing.TB, s string) []byte {
	t.Helper()
	b, err := hex.DecodeString(s)
	if err != nil {
		t.Fatalf("bad test vector %q: %v", s, err)
	}
	return b
}

// hexOf returns the upper-case hex of a PIN block, for comparing with vectors.
func hexOf(b []byte) string {
	return strings.ToUpper(hex.EncodeToString(b))
}
