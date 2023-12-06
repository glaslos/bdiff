// Copyright (c) 2023 Lukas Rist
//
// Copyright 2015 Monmohan Singh. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package bdiff

import (
	"bytes"
	"crypto/sha256"
	"hash"
	"io"

	"github.com/chmduquesne/rollinghash/adler32"
)

type processingResult struct {
	blockMatch   bool
	matchedBlock Block
	index        uint32
	consumed     int
	eof          bool
}

func matchesBlock(checksum uint32, sha256 [sha256.Size]byte, s Fingerprint) (Block, bool) {
	if sha2blk, ok := s.Blocks[checksum]; ok {
		if block, m := sha2blk[sha256]; m {
			return block, true
		}
	}
	return Block{}, false

}

func processRolling(r io.Reader, index uint32, fileSize uint32, f Fingerprint, delta *[]Block, h *adler32.Adler32) (processingResult, error) {
	diff := *delta
	db := &diff[len(diff)-1]

	var b []byte
	bw := bytes.NewBuffer(b)
	_, err := h.WriteWindow(bw)
	if err != nil {
		return processingResult{}, err
	}

	bufferSize := fileSize - (index + uint32(len(b)))
	if bufferSize == 0 {
		db.RawBytes = append(db.RawBytes, b...)
		*delta = diff
		return processingResult{false, Block{}, index, 0, true}, nil
	}
	fb := b[0]
	db.RawBytes = append(db.RawBytes, fb)

	buf := make([]byte, 1)
	n, err := io.ReadFull(r, buf)
	if err != nil {
		return processingResult{}, err
	}
	index++
	h.Roll(buf[0])

	bw = bytes.NewBuffer(b)
	_, err = h.WriteWindow(bw)
	if err != nil {
		return processingResult{}, err
	}
	matchblock, matched := matchesBlock(h.Sum32(), sha256.Sum256(b), f)
	return processingResult{matched, matchblock, index, n, false}, nil
}

func processBlock(r io.Reader, index uint32, fileSize uint32, f Fingerprint, delta *[]Block, h hash.Hash32) (processingResult, error) {
	buffSize := f.BlockSize
	if (index + f.BlockSize) > fileSize {
		buffSize = fileSize - index
	}
	if buffSize == 0 {
		return processingResult{false, Block{}, index, 0, true}, nil
	}

	buf := make([]byte, buffSize)
	if _, err := io.ReadFull(r, buf); err != nil {
		return processingResult{}, err
	}

	n, err := h.Write(buf)
	if err != nil {
		return processingResult{}, err
	}

	block, matched := matchesBlock(h.Sum32(), sha256.Sum256(buf), f)
	return processingResult{matched, block, index, n, false}, nil

}

func Diff(r io.Reader, fileSize uint32, f Fingerprint) ([]Block, error) {
	var (
		delta     []Block
		index     uint32
		result    processingResult
		blockMode bool
		err       error
	)
	h := adler32.New()
	blockMode = true
	for {
		if blockMode {
			result, err = processBlock(r, index, fileSize, f, &delta, h)
			if err != nil {
				return delta, err
			}
			index = result.index
			if result.eof {
				return delta, nil
			}
			if result.blockMatch {
				delta = append(delta, result.matchedBlock)
				index += uint32(result.consumed)
				continue
			}
			delta = append(delta, Block{HasData: true, Start: index})
			blockMode = false
		}
		result, err = processRolling(r, index, fileSize, f, &delta, h)
		if err != nil {
			return nil, err
		}
		index = result.index

		if result.eof {
			return delta, nil
		}
		if result.blockMatch {
			delta = append(delta, result.matchedBlock)
			index += uint32(result.consumed)
			blockMode = true
			continue
		}

	}

}
