package actions

import (
	"github.com/lienkolabs/swell/crypto"
	"github.com/lienkolabs/swell/util"
)

type Signin struct {
	Epoch   uint64
	Author  crypto.Token
	Reasons string
}

func (c *Signin) Serialize() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(c.Epoch, &bytes)
	util.PutToken(c.Author, &bytes)
	util.PutByte(ISignInAction, &bytes)
	return bytes
}

func ParseSignIn(create []byte) *Signin {
	action := Signin{}
	position := 0
	action.Epoch, position = util.ParseUint64(create, position)
	action.Author, position = util.ParseToken(create, position)
	if create[position] != ISignInAction {
		return nil
	}
	position += 1
	action.Reasons, position = util.ParseString(create, position)
	if position != len(create) {
		return nil
	}
	return &action
}
