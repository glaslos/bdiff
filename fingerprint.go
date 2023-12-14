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

type Block struct {
	Start, End uint32
	Checksum32 uint32
	Sha256hash [sha256.Size]byte
	HasData    bool
	RawBytes   []byte
}

type Fingerprint struct {
	BlockSize uint32
	Blocks    map[uint32]map[[sha256.Size]byte]Block
}

func NewFingerprint(src io.Reader, blockSize uint32) (Fingerprint, error) {
	buf := make([]byte, blockSize)

	n, start := 0, uint32(0)

	var (
		err   error
		block Block
	)

	fingerprint := Fingerprint{
		BlockSize: blockSize,
		Blocks:    make(map[uint32]map[[sha256.Size]byte]Block),
	}

	for {
		n, err = src.Read(buf)
		block = Block{
			Start:      start,
			End:        start + uint32(n),
			Checksum32: adler32.Checksum(buf[0:n]),
			Sha256hash: sha256.Sum256(buf[0:n]),
		}

		if sha2blk := fingerprint.Blocks[block.Checksum32]; sha2blk == nil {
			fingerprint.Blocks[block.Checksum32] = make(map[[sha256.Size]byte]Block)
		}
		fingerprint.Blocks[block.Checksum32][block.Sha256hash] = block

		start = block.End
		if err != nil {
			if err == io.EOF {
				return fingerprint, nil
			} else {
				return fingerprint, err
			}
		}
	}
}
