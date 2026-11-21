package formats_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// requireDebugNote checks that the debug writer received exactly the note
// SetDebugWriter promises: the operation and the format, and nothing else.
func requireDebugNote(t *testing.T, got, operation, format string) {
	t.Helper()
	require.Equal(t, "PIN block "+operation+" operation finished\nFormat  : "+format+"\n\n", got)
}
