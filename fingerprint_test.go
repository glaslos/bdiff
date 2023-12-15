package bdiff

import (
	"bytes"
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewFingerprint(t *testing.T) {
	p := make([]byte, 128)
	_, err := rand.Read(p)
	require.NoError(t, err)

	src := bytes.NewBuffer(p)
	fp, err := NewFingerprint(src, 12)
	require.NoError(t, err)
	require.Equal(t, uint32(12), fp.BlockSize)
	require.NotEmpty(t, fp.Blocks)

	fp.String()
	require.Equal(t, "", 1)
}
