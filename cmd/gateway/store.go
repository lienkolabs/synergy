package main

import (
	"log"
	"os"
	"sync"

	"github.com/lienkolabs/swell/util"
)

type block struct {
	data [][]byte
}

type blockchain struct {
	mu      sync.Mutex
	io      os.File
	blocks  []*block
	current *block
}

const (
	blocksignal  byte = 0
	actionsignal byte = 1
)

func newBlockBytes(epoch uint64) []byte {
	data := []byte{blocksignal}
	util.PutUint64(epoch, &data)
	return data
}

func (b *blockchain) NewBlock() {
	newBlock := &block{data: make([][]byte, 0)}
	b.current = newBlock
	b.blocks = append(b.blocks, b.current)
	epoch := uint64(len(b.blocks) - 1)
	data := newBlockBytes(epoch)
	if n, err := b.io.Write(data); n != len(data) || err != nil {
		log.Fatalf("could not write block: %v", err)
	}
}

func (b *blockchain) NewAction(action []byte) {
	b.mu.Lock()
	b.current.data = append(b.current.data, action)
	b.mu.Unlock()
	data := []byte{actionsignal}
	util.PutUint64(uint64(len(data)), &data)
	data = append(data, action...)
	if n, err := b.io.Write(data); n != len(data) || err != nil {
		log.Fatalf("could not write action: %v", err)
	}
}

func (b *blockchain) Close() {
	b.io.Close()
}

func OpenBlockchain(path string) *blockchain {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalf("could not open file: %v", err)
	}
	return &blockchain{
		mu: sync.Mutex{},
		io: *file,
	}
}
