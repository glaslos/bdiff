// Copyright (c) 2023 Lukas Rist
//
// Copyright 2015 Monmohan Singh. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package bdiff

import (
	"crypto/sha256"
	"io"
)

type processingResult struct {
	blockMatch   bool
	matchedBlock Block
	windowState  *State
	readPtr      int64
	eof          bool
}

func processBlock(r io.Reader, rptr uint32, filesz uint32, s Fingerprint, delta *[]Block) processingResult {
	brem := s.BlockSize
	if (rptr + s.BlockSize) > filesz {
		brem = filesz - rptr
	}
	if brem == 0 {
		return processingResult{false, Block{}, nil, rptr, true}
	}

	buf := make([]byte, brem)
	n, err := io.ReadFull(r, buf)
	if err != nil || n != brem {

	}

	checksum, state := Checksum(buf)
	matchblock, matched := matchBlock(checksum, sha256.Sum256(buf), s)
	return processingResult{matched, matchblock, state, rptr, false}

}

func processRolling(r io.Reader, st *State, rptr int64, filesz int64, s Fingerprint, delta *[]Block) processingResult {
	diff := *delta
	db := &diff[len(diff)-1]
	brem := filesz - (rptr + int64(len(st.window)))
	if brem == 0 {
		db.RawBytes = append(db.RawBytes, st.window...)
		*delta = diff
		return processingResult{false, Block{}, nil, rptr, true}
	}
	fb := st.window[0]
	db.RawBytes = append(db.RawBytes, fb)
	b := make([]byte, 1)
	_, e := io.ReadFull(r, b)
	if e != nil {

	}
	rptr++
	checksum := st.UpdateWindow(b[0])
	matchblock, matched := matchBlock(checksum, sha256.Sum256(st.window), s)
	return processingResult{matched, matchblock, st, rptr, false}
}

func Diff(r io.Reader, filesz int64, s Fingerprint) []Block {
	var (
		delta     []Block
		state     *State
		rptr      int64
		result    processingResult
		blockMode bool
	)
	blockMode = true
	for {
		if blockMode {
			result = processBlock(r, rptr, filesz, s, &delta)
			rptr = result.readPtr
			state = result.windowState
			if result.eof {
				return delta
			}
			if result.blockMatch {
				delta = append(delta, result.matchedBlock)
				rptr += int64(len(state.window))
				continue
			}
			delta = append(delta, Block{HasData: true, Start: rptr})
			blockMode = false
		}
		result = processRolling(r, state, rptr, filesz, s, &delta)
		rptr = result.readPtr
		state = result.windowState

		if result.eof {
			return delta
		}
		if result.blockMatch {
			delta = append(delta, result.matchedBlock)
			rptr += int64(len(state.window))
			blockMode = true
			continue
		}

	}

}
