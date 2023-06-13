package api

import (
	"io"

	"github.com/lienkolabs/swell/util"
	"github.com/lienkolabs/synergy/social"
	"github.com/lienkolabs/synergy/social/state"
)

// API -> receives instructions from UI
// API 	-> receive instructions from gateway
// 		-> translate to UI updates

type Listener interface {
	Listen() []byte
}

type Gateway struct {
	state state.State
}

type Blockchain struct {
	listener Listener
	file     io.WriteCloser
	epoch    uint64
}

func (b *Blockchain) Write(data []byte) error {
	_, err := b.file.Write(data)
	return err
}

func (b *Blockchain) Listen() error {
	data := b.listener.Listen()
	kind := data[2]
	if kind == social.ZeroType {
		b.epoch, _ = util.ParseUint64(data, 3)
		return b.Write(data)
	} else {

	}

}
