package actions

import (
	"github.com/lienkolabs/swell/crypto"
	"github.com/lienkolabs/swell/util"
)

type CreateBoardAction struct {
	Epoch       uint64
	Author      crypto.Token
	Reasons     string
	OnBehalfOf  string
	Name        string
	Description string
	Keywords    []string
	PinMajority byte
}

func (c *CreateBoardAction) Serialize() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(c.Epoch, &bytes)
	util.PutToken(c.Author, &bytes)
	util.PutByte(ICreateBoardAction, &bytes)
	util.PutString(c.Reasons, &bytes)
	util.PutString(c.OnBehalfOf, &bytes)
	util.PutString(c.Name, &bytes)
	util.PutString(c.Description, &bytes)
	PutKeywords(c.Keywords, &bytes)
	util.PutByte(c.PinMajority, &bytes)
	return bytes
}

func ParseCreateAction(create []byte) *CreateBoardAction {
	action := CreateBoardAction{}
	position := 0
	action.Epoch, position = util.ParseUint64(create, position)
	action.Author, position = util.ParseToken(create, position)
	if create[position] != ICreateBoardAction {
		return nil
	}
	position += 1
	action.Reasons, position = util.ParseString(create, position)
	action.OnBehalfOf, position = util.ParseString(create, position)
	action.Name, position = util.ParseString(create, position)
	action.Description, position = util.ParseString(create, position)
	action.Keywords, position = ParseKeywords(create, position)
	action.PinMajority, position = util.ParseByte(create, position)
	if position != len(create) {
		return nil
	}
	return &action
}

type UpdateBoardAction struct {
	Epoch       uint64
	Author      crypto.Token
	Reasons     string
	Board       string
	Description string
	Keywords    []string
	PinMajority byte
}

func (c *UpdateBoardAction) Serialize() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(c.Epoch, &bytes)
	util.PutToken(c.Author, &bytes)
	util.PutByte(IUpdateBoardAction, &bytes)
	util.PutString(c.Reasons, &bytes)
	util.PutString(c.Board, &bytes)
	util.PutString(c.Description, &bytes)
	PutKeywords(c.Keywords, &bytes)
	util.PutByte(c.PinMajority, &bytes)
	return bytes
}

func ParseUpdateBoardAction(create []byte) *UpdateBoardAction {
	action := UpdateBoardAction{}
	position := 0
	action.Epoch, position = util.ParseUint64(create, position)
	action.Author, position = util.ParseToken(create, position)
	if create[position] != IUpdateBoardAction {
		return nil
	}
	position += 1
	action.Reasons, position = util.ParseString(create, position)
	action.Board, position = util.ParseString(create, position)
	action.Description, position = util.ParseString(create, position)
	action.Keywords, position = ParseKeywords(create, position)
	action.PinMajority, position = util.ParseByte(create, position)
	if position != len(create) {
		return nil
	}
	return &action
}

type PinAction struct {
	Epoch   uint64
	Author  crypto.Token
	Reasons string
	Board   string
	Draft   crypto.Hash
	Pin     bool
}

func (c *PinAction) Serialize() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(c.Epoch, &bytes)
	util.PutToken(c.Author, &bytes)
	util.PutByte(IPinAction, &bytes)
	util.PutString(c.Reasons, &bytes)
	util.PutString(c.Board, &bytes)
	util.PutHash(c.Draft, &bytes)
	util.PutBool(c.Pin, &bytes)
	return bytes
}

func ParsePinAction(create []byte) *PinAction {
	action := PinAction{}
	position := 0
	action.Epoch, position = util.ParseUint64(create, position)
	action.Author, position = util.ParseToken(create, position)
	if create[position] != IPinAction {
		return nil
	}
	position += 1
	action.Reasons, position = util.ParseString(create, position)
	action.Board, position = util.ParseString(create, position)
	action.Draft, position = util.ParseHash(create, position)
	action.Pin, position = util.ParseBool(create, position)
	if position != len(create) {
		return nil
	}
	return &action
}

type BoardEditorAction struct {
	Epoch   uint64
	Author  crypto.Token
	Reasons string
	Board   string
	Editor  crypto.Token
	Insert  bool
}

func (c *BoardEditorAction) Serialize() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(c.Epoch, &bytes)
	util.PutToken(c.Author, &bytes)
	util.PutByte(IBoardEditorAction, &bytes)
	util.PutString(c.Reasons, &bytes)
	util.PutString(c.Board, &bytes)
	util.PutToken(c.Editor, &bytes)
	util.PutBool(c.Insert, &bytes)
	return bytes
}

func ParseBoardEditorAction(create []byte) *BoardEditorAction {
	action := BoardEditorAction{}
	position := 0
	action.Epoch, position = util.ParseUint64(create, position)
	action.Author, position = util.ParseToken(create, position)
	if create[position] != IBoardEditorAction {
		return nil
	}
	position += 1
	action.Reasons, position = util.ParseString(create, position)
	action.Board, position = util.ParseString(create, position)
	action.Editor, position = util.ParseToken(create, position)
	action.Insert, position = util.ParseBool(create, position)
	if position != len(create) {
		return nil
	}
	return &action
}
