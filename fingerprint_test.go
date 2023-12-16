package bdiff

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewFingerprint(t *testing.T) {
	p := []byte("the lazy cat sleeps in the sun")
	src := bytes.NewBuffer(p)
	fp, err := NewFingerprint(src, 12)
	require.NoError(t, err)
	require.Equal(t, uint32(12), fp.BlockSize)
	require.NotEmpty(t, fp.Blocks)
}
