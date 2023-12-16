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
	"fmt"
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

func matchesBlock(checksum uint32, sha256 [sha256.Size]byte, fp Fingerprint) (Block, bool) {
	if sha2blk, ok := fp.Blocks[checksum]; ok {
		if block, m := sha2blk[sha256]; m {
			return block, true
		}
	}
	return Block{}, false

}

func processRolling(r io.Reader, index uint32, fileSize int, f Fingerprint, delta *Delta, h *adler32.Adler32) (processingResult, error) {
	diff := *delta
	db := &diff[len(diff)-1]

	b := []byte{}
	bw := bytes.NewBuffer(b)
	_, err := h.WriteWindow(bw)
	if err != nil {
		return processingResult{}, err
	}

	bufferSize := uint32(fileSize) - (index + uint32(bw.Len()))
	if bufferSize == 0 {
		db.RawBytes = append(db.RawBytes, b...)
		*delta = diff
		return processingResult{false, Block{}, index, 0, true}, nil
	}

	fb := bw.Bytes()[0]
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
	block, matched := matchesBlock(h.Sum32(), sha256.Sum256(bw.Bytes()), f)
	return processingResult{matched, block, index, n, false}, nil
}

func processBlock(src io.Reader, index uint32, fileSize int, fp Fingerprint, delta *Delta) (processingResult, *adler32.Adler32, error) {
	h := adler32.New()
	buffSize := fp.BlockSize
	if (index + fp.BlockSize) > uint32(fileSize) {
		buffSize = uint32(fileSize) - index
	}

	if buffSize == 0 {
		return processingResult{false, Block{}, index, 0, true}, h, nil
	}

	buf := make([]byte, buffSize)
	if _, err := io.ReadFull(src, buf); err != nil {
		return processingResult{}, nil, fmt.Errorf("failed to read block: %w", err)
	}

	_, err := h.Write(buf)
	if err != nil {
		return processingResult{}, nil, err
	}

	block, matched := matchesBlock(h.Sum32(), sha256.Sum256(buf), fp)
	return processingResult{matched, block, index, int(buffSize), false}, h, nil

}

func Diff(src io.Reader, fileSize int, f Fingerprint) (Delta, error) {
	var (
		delta     Delta
		index     uint32
		result    processingResult
		blockMode bool
		err       error
		h         *adler32.Adler32
	)
	blockMode = true
	for {
		if blockMode {
			result, h, err = processBlock(src, index, fileSize, f, &delta)
			if err != nil {
				return delta, fmt.Errorf("failed to process block: %w", err)
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
		result, err = processRolling(src, index, fileSize, f, &delta, h)
		if err != nil {
			return nil, err
		}
		index = result.index + uint32(result.consumed)

		if result.eof {
			return delta, nil
		}
		if result.blockMatch {
			delta = append(delta, result.matchedBlock)
			blockMode = true
			continue
		}

	}

}
