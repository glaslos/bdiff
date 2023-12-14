package bdiff

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDiff(t *testing.T) {
	p1 := []byte("the lazy dog sleeps in the sun")
	p2 := []byte("the lazy cat sleeps in the sun")

	src1 := bytes.NewBuffer(p1)
	fp, err := NewFingerprint(src1, 4)
	require.NoError(t, err)

	src2 := bytes.NewBuffer(p2)
	diff, err := Diff(src2, uint32(len(p2)), fp)
	require.NoError(t, err)
	require.NotEmpty(t, diff)
}
