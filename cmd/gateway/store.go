package main

import (
	"log"
	"os"
	"sync"

	"github.com/lienkolabs/breeze/util"
)

type block struct {
	data [][]byte
}

type blockchain struct {
	mu      sync.Mutex
	io      *os.File
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

func OpenBlockchain() (*blockchain, bool) {
	exists := true
	if stat, err := os.Stat("chain.dat"); err != nil || stat.Size() == 0 {
		exists = false
	}
	file, err := os.OpenFile("chain.dat", os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalf("could not access chain file: %v\n", err)
	}
	b := &blockchain{
		mu: sync.Mutex{},
		io: file,
	}
	signal := make([]byte, 9)
	for {
		if n, err := b.io.Read(signal); n != 9 || err != nil {
			break
		}
		number, _ := util.ParseUint64(signal, 1)
		if signal[0] == blocksignal {
			if number == 0 {
				if len(b.blocks) != 0 {
					log.Fatal("blockchain file corrupted")
				}
			} else if number != uint64(len(b.blocks)) {
				log.Fatal("blockchain file corrupted")
			} else {
				b.NewBlock()
			}
		} else if signal[0] == actionsignal {
			data := make([]byte, int(number))
			if n, _ := b.io.Read(signal); n != len(data) {
				log.Fatal("blockchain file corrupted")
			}
			b.NewAction(data)
		} else {
			log.Fatal("blockchain file corrupted")
		}
	}
	return b, exists
}