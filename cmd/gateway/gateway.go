package main

import (
	"errors"
	"fmt"
	"net"
	"sync"

	"github.com/lienkolabs/breeze/crypto"
	"github.com/lienkolabs/breeze/network/trusted"
	"github.com/lienkolabs/synergy/social/state"
)

// Cached is used to buffer actions while the receiving node is syncing
type Cached struct {
	mu    sync.Mutex
	conn  *trusted.SignedConnection
	cache [][]byte
	ready bool
	alive bool
}

func NewCached(conn *trusted.SignedConnection) *Cached {
	return &Cached{
		mu:    sync.Mutex{},
		conn:  conn,
		cache: make([][]byte, 0),
		ready: false,
		alive: true,
	}
}

// SendDirect sends data directly to the connection without transiting the cache
func (c *Cached) SendDirect(data []byte) error {
	return c.conn.Send(data)
}

// Send wither sends data to the node or caches it to send it latet
func (c *Cached) Send(data []byte) error {
	if c.ready {
		return c.Send(data)
	}
	if !c.alive {
		return fmt.Errorf("connection is dead")
	}
	c.mu.Lock()
	c.cache = append(c.cache, data)
	c.mu.Unlock()
	return nil
}

// a gateway provides connectivity to submit actions to the breeze network
type Network struct {
	mu      sync.Mutex
	inbound map[crypto.Token]*Cached
	Actions chan []byte
	chain   *blockchain
	synergy *state.State
}

// sync a new connection
func (n *Network) Sync(conn *Cached) {
	n.mu.Lock()
	syncEpoch := uint64(len(n.chain.blocks)) + 1
	syncAction := len(n.chain.current.data)
	n.mu.Unlock()
	for epoch := uint64(0); epoch < syncEpoch-1; epoch++ {
		conn.SendDirect(newBlockBytes(uint64(epoch)))
		for _, action := range n.chain.blocks[int(epoch)].data {
			conn.Send(append([]byte{actionsignal}, action...))
		}
	}
	currentBlockCache := make([][]byte, syncAction)
	n.chain.mu.Lock()
	currentBlock := n.chain.blocks[(int(syncEpoch))-1].data
	for n := 0; n < syncAction; n++ {
		currentBlockCache[n] = currentBlock[n]
	}
	n.chain.mu.Unlock()
	for _, action := range currentBlockCache {
		conn.Send(append([]byte{actionsignal}, action...))
	}
	for {
		conn.mu.Lock()
		if len(conn.cache) == 0 {
			conn.ready = true
			conn.mu.Unlock()
			return
		}
		data := conn.cache[0]
		conn.cache = conn.cache[1:]
		conn.mu.Unlock()
		if err := conn.SendDirect(data); err != nil {
			conn.conn.Shutdown()
			conn.alive = false
			return
		}
	}
}

func NewActionsGateway(port int, credentials crypto.PrivateKey, chain *blockchain) (chan trusted.Message, error) {
	validate := trusted.AcceptAllConnections
	listeners, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		return nil, err
	}

	genesis := state.GenesisState(nil)
	if genesis == nil {
		return nil, errors.New("could not create genesis state")
	}
	gateway := Network{
		mu:      sync.Mutex{},
		inbound: make(map[crypto.Token]*Cached),
		Actions: make(chan []byte),
		chain:   chain,
		synergy: genesis,
	}

	for _, block := range chain.blocks {
		for _, action := range block.data {
			if err := genesis.Action(action); err != nil {
				return nil, fmt.Errorf("blockchain has invalid action: %v", err)
			}
		}
	}

	messages := make(chan trusted.Message)
	shutDown := make(chan crypto.Token) // receive connection shutdown

	go func() {
		for {
			if conn, err := listeners.Accept(); err == nil {
				trustedConn, err := trusted.PromoteConnection(conn, credentials, validate)
				if err != nil {
					conn.Close()
				} else {
					cached := NewCached(trustedConn)
					gateway.inbound[trustedConn.Token] = cached
					go gateway.Sync(cached)
					trustedConn.Listen(messages, shutDown)
				}
			} else {
				return
			}
		}
	}()

	go func() {
		for {
			select {
			case token := <-shutDown:
				gateway.mu.Lock()
				delete(gateway.inbound, token)
				gateway.mu.Unlock()
			case msg := <-messages:
				chain.NewAction(msg.Data)
			}
		}
	}()

	return messages, nil
}