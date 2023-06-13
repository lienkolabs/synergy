package api

import (
	"github.com/lienkolabs/swell/crypto"
	"github.com/lienkolabs/swell/util"
	"github.com/lienkolabs/synergy/social/actions"
)

type Attorney struct {
	author  crypto.Token
	pk      crypto.PrivateKey
	wallet  crypto.PrivateKey
	pending map[crypto.Hash]actions.Action
	epoch   uint64
}

func NewAttorney(pk, wallet crypto.PrivateKey, token crypto.Token) *Attorney {
	return &Attorney{
		pk:     pk,
		wallet: wallet,
	}
}

func (a *Attorney) DressAction(action actions.Action) []byte {
	bytes := action.Serialize()
	dress := []byte{0}
	dress = append(dress, bytes[0:8+crypto.TokenSize]...)
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
