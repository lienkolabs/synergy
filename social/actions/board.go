package actions

import (
	"github.com/lienkolabs/swell/crypto"
	"github.com/lienkolabs/swell/util"
)

type CreateBoard struct {
	Epoch       uint64
	Author      crypto.Token
	Reasons     string
	OnBehalfOf  string
	Name        string
	Description string
	Keywords    []string
	PinMajority byte
}

func (c *CreateBoard) Serialize() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(c.Epoch, &bytes)
	util.PutToken(c.Author, &bytes)
	util.PutByte(ACreateBoard, &bytes)
	util.PutString(c.Reasons, &bytes)
	util.PutString(c.OnBehalfOf, &bytes)
	util.PutString(c.Name, &bytes)
	util.PutString(c.Description, &bytes)
	PutKeywords(c.Keywords, &bytes)
	util.PutByte(c.PinMajority, &bytes)
	return bytes
}

func ParseCreateBoard(create []byte) *CreateBoard {
	action := CreateBoard{}
	position := 0
	action.Epoch, position = util.ParseUint64(create, position)
	action.Author, position = util.ParseToken(create, position)
	if create[position] != ACreateBoard {
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

type UpdateBoard struct {
	Epoch       uint64
	Author      crypto.Token
	Reasons     string
	Board       string
	Description *string
	Keywords    *[]string
	PinMajority *byte
}

func (c *UpdateBoard) Serialize() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(c.Epoch, &bytes)
	util.PutToken(c.Author, &bytes)
	util.PutByte(AUpdateBoard, &bytes)
	util.PutString(c.Reasons, &bytes)
	util.PutString(c.Board, &bytes)
	if c.Description != nil {
		util.PutByte(1, &bytes)
		util.PutString(*c.Description, &bytes)
	} else {
		util.PutByte(0, &bytes)
	}
	if c.Keywords != nil {
		util.PutByte(1, &bytes)
		PutKeywords(*c.Keywords, &bytes)
	} else {
		util.PutByte(0, &bytes)
	}
	if c.PinMajority != nil {
		util.PutByte(1, &bytes)
		util.PutByte(*c.PinMajority, &bytes)
	} else {
		util.PutByte(0, &bytes)
	}
	return bytes
}

func ParseUpdateBoard(create []byte) *UpdateBoard {
	action := UpdateBoard{}
	position := 0
	action.Epoch, position = util.ParseUint64(create, position)
	action.Author, position = util.ParseToken(create, position)
	if create[position] != AUpdateBoard {
		return nil
	}
	position += 1
	action.Reasons, position = util.ParseString(create, position)
	action.Board, position = util.ParseString(create, position)
	if create[position] == 0 {
		position += 1
	} else {
		var des string
		position += 1
		des, position = util.ParseString(create, position)
		action.Description = &des
	}
	if create[position] == 0 {
		position += 1
	} else {
		var key []string
		position += 1
		key, position = ParseKeywords(create, position)
		action.Keywords = &key
	}
	if create[position] == 0 {
		position += 1
	} else {
		var pin byte
		position += 1
		pin, position = util.ParseByte(create, position)
		action.PinMajority = &pin
	}
	if position != len(create) {
		return nil
	}
	return &action
}

type Pin struct {
	Epoch   uint64
	Author  crypto.Token
	Reasons string
	Board   string
	Draft   crypto.Hash
	Pin     bool
}

func (c *Pin) Serialize() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(c.Epoch, &bytes)
	util.PutToken(c.Author, &bytes)
	util.PutByte(APin, &bytes)
	util.PutString(c.Reasons, &bytes)
	util.PutString(c.Board, &bytes)
	util.PutHash(c.Draft, &bytes)
	util.PutBool(c.Pin, &bytes)
	return bytes
}

func ParsePin(create []byte) *Pin {
	action := Pin{}
	position := 0
	action.Epoch, position = util.ParseUint64(create, position)
	action.Author, position = util.ParseToken(create, position)
	if create[position] != APin {
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

type BoardEditor struct {
	Epoch   uint64
	Author  crypto.Token
	Reasons string
	Board   string
	Editor  crypto.Token
	Insert  bool
}

func (c *BoardEditor) Serialize() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(c.Epoch, &bytes)
	util.PutToken(c.Author, &bytes)
	util.PutByte(ABoardEditor, &bytes)
	util.PutString(c.Reasons, &bytes)
	util.PutString(c.Board, &bytes)
	util.PutToken(c.Editor, &bytes)
	util.PutBool(c.Insert, &bytes)
	return bytes
}

func ParseBoardEditor(create []byte) *BoardEditor {
	action := BoardEditor{}
	position := 0
	action.Epoch, position = util.ParseUint64(create, position)
	action.Author, position = util.ParseToken(create, position)
	if create[position] != ABoardEditor {
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
