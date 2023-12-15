package bdiff

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPatch(t *testing.T) {
	p1 := []byte("the lazy dog sleeps in the sun")
	p2 := []byte("the lazy cat sleeps in the sun")

	src := bytes.NewReader(p1)
	fp, err := NewFingerprint(src, 1)
	require.NoError(t, err)

	dst := bytes.NewBuffer(p2)
	diff, err := Diff(dst, uint32(len(p2)), fp)
	require.NoError(t, err)
	require.NotEmpty(t, diff)

	_, err = src.Seek(0, io.SeekStart)
	require.NoError(t, err)

	err = Patch(diff, src, dst)
	require.NoError(t, err)
	require.Equal(t, string(p2), dst.String())
}
