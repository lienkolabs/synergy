package social

import (
	"fmt"
	"log"
	"sync"

	"github.com/lienkolabs/breeze/crypto"
	"github.com/lienkolabs/breeze/network/trusted"
	"github.com/lienkolabs/breeze/util"
	"github.com/lienkolabs/synergy/social/state"
)

type Proxy struct {
	mu      sync.Mutex
	state   *state.State
	conn    *trusted.SignedConnection
	viewers []chan uint64
	epoch   uint64
}

func (p *Proxy) Stop() {
	p.conn.Shutdown()
}

func (p *Proxy) State() *state.State {
	return p.state
}

func (p *Proxy) Epoch() uint64 {
	return p.epoch
}

func (p *Proxy) Action(data []byte) {
	undressed := Undress(data)
	if err := p.conn.Send(undressed); err != nil {
		log.Printf("error sending action: %v", err)
	}
}

func (p *Proxy) Register() chan uint64 {
	viewer := make(chan uint64)
	p.mu.Lock()
	defer p.mu.Unlock()
	p.viewers = append(p.viewers, viewer)
	return viewer
}

func SelfProxyState(host string, hostToken crypto.Token, credential crypto.PrivateKey, genesis *state.State) *Proxy {
	conn, err := trusted.Dial(host, credential, hostToken)
	if err != nil {
		log.Fatalf("could not connect to host: %v", err)
	}
	proxy := &Proxy{
		mu:      sync.Mutex{},
		state:   genesis,
		conn:    conn,
		viewers: make([]chan uint64, 0),
		epoch:   0,
	}

	go func() {
		for {
			data, err := conn.Read()
			if err != nil {
				log.Fatalf("error reading from host: %v", err)
			}
			if data[0] == 0 {
				if len(data) == 9 {
					proxy.epoch, _ = util.ParseUint64(data, 1)
					proxy.mu.Lock()
					for _, v := range proxy.viewers {
						v <- proxy.epoch
					}
					proxy.mu.Unlock()
				} else {
					log.Print("invalid epoch message")
				}
			} else if data[0] == 1 {
				if len(data) > 1 {
					action := data[1:]
					if err := proxy.state.Action(action); err != nil {
						log.Printf("invalid action: %v", err)
					} else {
						fmt.Println("action received")
					}
				}
			}
		}
	}()
	return proxy
}
