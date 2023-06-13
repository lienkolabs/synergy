package actions

import (
	"github.com/lienkolabs/swell/crypto"
	"github.com/lienkolabs/swell/util"
)

type ImprintStampAction struct {
	Epoch      uint64
	Author     crypto.Token
	Reasons    string
	OnBehalfOf string
	Hash       crypto.Hash
}

func (c *ImprintStampAction) Serialize() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(c.Epoch, &bytes)
	util.PutToken(c.Author, &bytes)
	util.PutByte(IImprintStampAction, &bytes)
	util.PutString(c.Reasons, &bytes)
	util.PutHash(c.Hash, &bytes)
	return bytes
}

func ParseImprintStampAction(create []byte) *ImprintStampAction {
	action := ImprintStampAction{}
	position := 0
	action.Epoch, position = util.ParseUint64(create, position)
	action.Author, position = util.ParseToken(create, position)
	if create[position] != IImprintStampAction {
		return nil
	}
	position += 1
	action.Reasons, position = util.ParseString(create, position)
	action.Hash, position = util.ParseHash(create, position)
	if position != len(create) {
		return nil
	}
	return &action
}
