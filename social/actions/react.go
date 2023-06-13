package actions

import (
	"github.com/lienkolabs/swell/crypto"
	"github.com/lienkolabs/swell/util"
)

type ReactAction struct {
	Epoch      uint64
	Author     crypto.Token
	Reasons    string
	OnBehalfOf string
	Hash       crypto.Hash
	Reaction   byte
}

func (c *ReactAction) Serialize() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(c.Epoch, &bytes)
	util.PutToken(c.Author, &bytes)
	util.PutByte(IReactAction, &bytes)
	util.PutString(c.Reasons, &bytes)
	util.PutHash(c.Hash, &bytes)
	util.PutByte(c.Reaction, &bytes)
	return bytes
}

func ParseReactAction(create []byte) *ReactAction {
	action := ReactAction{}
	position := 0
	action.Epoch, position = util.ParseUint64(create, position)
	action.Author, position = util.ParseToken(create, position)
	if create[position] != IReactAction {
		return nil
	}
	position += 1
	action.Reasons, position = util.ParseString(create, position)
	action.Hash, position = util.ParseHash(create, position)
	action.Reaction, position = util.ParseByte(create, position)
	if position != len(create) {
		return nil
	}
	return &action
}
