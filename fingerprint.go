// Copyright (c) 2023 Lukas Rist
//
// Copyright 2015 Monmohan Singh. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package bdiff

import (
	"crypto/sha256"
	"hash/adler32"
	"io"
)

func NewFingerprint(src io.Reader, blockSize uint32) (Fingerprint, error) {
	var (
		n     int
		start uint32
		err   error
		block Block
	)

	fingerprint := Fingerprint{
		BlockSize: blockSize,
		Blocks:    make(map[uint32]map[[sha256.Size]byte]Block),
	}

	buf := make([]byte, blockSize)
	for {
		n, err = src.Read(buf)
		if n != 0 {
			block = Block{
				Start:      start,
				End:        start + uint32(n),
				Checksum32: adler32.Checksum(buf[0:n]),
				Sha256hash: sha256.Sum256(buf[0:n]),
			}

			if ok := fingerprint.Blocks[block.Checksum32]; ok == nil {
				fingerprint.Blocks[block.Checksum32] = make(map[[sha256.Size]byte]Block)
			}
			fingerprint.Blocks[block.Checksum32][block.Sha256hash] = block

			start = block.End
		}
		if err != nil {
			if err == io.EOF {
				return fingerprint, nil
			} else {
				return fingerprint, err
			}
		}
	}
}
