package api

import (
	"github.com/lienkolabs/swell/crypto"
	"github.com/lienkolabs/swell/util"
	"github.com/lienkolabs/synergy/social"
	"github.com/lienkolabs/synergy/social/actions"
)

type Attorney struct {
	author  crypto.Token
	pk      crypto.PrivateKey
	wallet  crypto.PrivateKey
	pending map[crypto.Hash]actions.Action
	epoch   uint64
	gateway chan []byte
}

func NewAttorneyServer(pk crypto.PrivateKey, token crypto.Token, port int, gateway *social.Gateway) *Attorney {
	attorney := Attorney{
		author:  token,
		pk:      pk,
		wallet:  pk,
		pending: make(map[crypto.Hash]actions.Action),
		epoch:   0,
	}
	blockEvent := gateway.Register()
	send := make(chan actions.Action)
	go func() {
		for {
			select {
			case epoch := <-blockEvent:
				attorney.epoch = epoch
			case action := <-send:
				gateway.Action(attorney.DressAction(action))
			}
		}
	}()

	return nil
}

func NewAttorney(pk, wallet crypto.PrivateKey, token crypto.Token) *Attorney {
	return &Attorney{
		pk:     pk,
		wallet: wallet,
	}
}

func (a *Attorney) Send(all []actions.Action) {
	for _, action := range all {
		dressed := a.DressAction(action)
		a.gateway <- dressed
	}
}

// Dress a giving action with current epoch, attorney´s author
// attorneys signature, attorneys wallet and wallet signature
func (a *Attorney) DressAction(action actions.Action) []byte {
	bytes := action.Serialize()
	dress := []byte{0}
	util.PutUint64(a.epoch, &dress)
	util.PutToken(a.author, &dress)
	dress = append(dress, 0, 1, 0, 0) // axe void synergy
	dress = append(dress, bytes[8+crypto.TokenSize:]...)
	util.PutToken(a.pk.PublicKey(), &dress)
	signature := a.pk.Sign(dress)
	util.PutSignature(signature, &dress)
	util.PutToken(a.wallet.PublicKey(), &dress)
	util.PutUint64(0, &dress) // zero fee
	signature = a.wallet.Sign(dress)
	util.PutSignature(signature, &dress)
	return dress
}

func (a *Attorney) Confirmed(hash crypto.Hash) {
	delete(a.pending, hash)
}
