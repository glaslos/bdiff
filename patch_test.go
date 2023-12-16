package bdiff

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPatch(t *testing.T) {
	p1 := []byte("the lazy dog sleeps in the sun")
	p2 := []byte("the lazy cat sleeps in the sun")

	src := bytes.NewReader(p1)
	fp, err := NewFingerprint(src, 4)
	require.NoError(t, err)

	dst := bytes.NewBuffer(p2)
	diff, err := Diff(dst, len(p2), fp)
	require.NoError(t, err)
	require.NotEmpty(t, diff)

	err = Patch(diff, src, dst)
	require.NoError(t, err)
	require.Equal(t, string(p2), dst.String())
}
