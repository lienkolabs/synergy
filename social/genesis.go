package social

import (
	"fmt"
	"sync"
	"time"

	"github.com/lienkolabs/swell/crypto"
	"github.com/lienkolabs/synergy/social/state"
)

func TestGenesisState(users map[crypto.Token]string) *state.State {
	genesis := state.GenesisState()
	for user, handle := range users {
		genesis.Members[crypto.HashToken(user)] = handle
	}
	return genesis
}

type Gateway struct {
	mu       *sync.Mutex
	incoming chan []byte
	newBlock []chan uint64
	stop     chan chan struct{}
	State    *state.State
}

func (g *Gateway) Stop() {
	resp := make(chan struct{})
	g.stop <- resp
	<-resp
}

func (g *Gateway) Action(data []byte) {
	g.incoming <- data
}

func (g *Gateway) Register() chan uint64 {
	g.mu.Lock()
	defer g.mu.Unlock()
	blockEvent := make(chan uint64)
	g.newBlock = append(g.newBlock, blockEvent)
	return blockEvent
}

func SelfGateway(engine *state.State) *Gateway {
	gateway := &Gateway{
		mu:       &sync.Mutex{},
		incoming: make(chan []byte),
		newBlock: make([]chan uint64, 0),
		stop:     make(chan chan struct{}),
		State:    engine,
	}

	ticker := time.NewTicker(time.Second)

	go func() {
		defer gateway.mu.Unlock()
		for {
			select {
			case <-ticker.C:
				gateway.mu.Lock()
				engine.Epoch += 1
				for _, emit := range gateway.newBlock {
					emit <- engine.Epoch
				}
				gateway.mu.Unlock()
			case action := <-gateway.incoming:
				if err := engine.Action(action); err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("Action performed")
				}
			case resp := <-gateway.stop:
				gateway.mu.Lock()
				for _, event := range gateway.newBlock {
					close(event)
				}
				resp <- struct{}{}
				return
			}
		}
	}()
	return gateway
}

func Undress(data []byte) []byte {
	// ignore first byte (breeze version)
	head := data[1 : 8+crypto.TokenSize+1]
	// ignore protocol id and tail
	// tail = wallet signature + fee + wallet + attorney signature + attorney
	tailSize := 2*crypto.SignatureSize + 2*crypto.TokenSize + 8
	return append(head, data[8+crypto.TokenSize+1+4:len(data)-tailSize]...)
}
