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
	String     string
}

func (b *Block) Print() {
	fmt.Printf("\t%d:%d: %v: %x\t%s\n", b.Start, b.End, b.HasData, b.Checksum32, b.String)
}

type Delta []Block

func (d Delta) Print() {
	for _, block := range d {
		block.Print()
	}
}

type Fingerprint struct {
	BlockSize uint32
	Blocks    map[uint32]map[[sha256.Size]byte]Block
}

func (fp *Fingerprint) Print() {
	for checksum, blocks := range fp.Blocks {
		fmt.Printf("%d: ", checksum)
		for _, block := range blocks {
			block.Print()
		}
	}
}
