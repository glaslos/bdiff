package bdiff

import (
	"crypto/sha256"
	"fmt"
)

type Block struct {
	Start, End uint32
	Checksum32 uint32
	Sha256hash [sha256.Size]byte
	HasData    bool
	RawBytes   []byte
}

type Delta []Block

type Fingerprint struct {
	BlockSize uint32
	Blocks    map[uint32]map[[sha256.Size]byte]Block
}

func (fp *Fingerprint) String() {
	for checksum, blocks := range fp.Blocks {
		fmt.Printf("%d: ", checksum)
		for sha, block := range blocks {
			fmt.Printf("\t%d:%d %x: %v: %x\n", block.Start, block.End, sha, block.HasData, block.Sha256hash)
		}
	}
}
