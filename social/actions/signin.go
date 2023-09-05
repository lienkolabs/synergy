package actions

import (
	"github.com/lienkolabs/swell/crypto"
	"github.com/lienkolabs/swell/util"
)

type Signin struct {
	Epoch   uint64
	Author  crypto.Token
	Reasons string
	Handle  string // provided by the protocol connection rules
}

func (c *Signin) Authored() crypto.Token {
	return c.Author
}

func (c *Signin) Serialize() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(c.Epoch, &bytes)
	util.PutToken(c.Author, &bytes)
	util.PutByte(ASignIn, &bytes)
	util.PutString(c.Reasons, &bytes)
	util.PutString(c.Handle, &bytes)
	return bytes
}

func ParseSignIn(create []byte) *Signin {
	action := Signin{}
	position := 0
	action.Epoch, position = util.ParseUint64(create, position)
	action.Author, position = util.ParseToken(create, position)
	if create[position] != ASignIn {
		return nil
	}
	position += 1
	action.Reasons, position = util.ParseString(create, position)
	action.Handle, position = util.ParseString(create, position)
	if position != len(create) {
		return nil
	}
	return &action
}
