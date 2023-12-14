// Copyright (c) 2023 Lukas Rist
//
// Copyright 2015 Monmohan Singh. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package bdiff

import (
	"io"
)

// Patch the delta from the source to the destination
func Patch(delta []Block, src io.ReadSeeker, dst io.Writer) error {
	for _, block := range delta {
		if block.HasData {
			if _, err := dst.Write(block.RawBytes); err != nil {
				return err
			}
		} else {
			if _, err := src.Seek(int64(block.Start), io.SeekStart); err != nil {
				return err
			}

			if _, err := io.CopyN(dst, src, int64(block.End-block.Start)); err != nil {
				return err
			}
		}
	}
	return nil
}
