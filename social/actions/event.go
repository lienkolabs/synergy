package actions

import (
	"time"

	"github.com/lienkolabs/swell/crypto"
	"github.com/lienkolabs/swell/util"
)

type CreateEventAction struct {
	Epoch        uint64
	Author       crypto.Token
	Reasons      string
	OnBehalfOf   string // non-optional
	StartAt      time.Time
	EstimatedEnd time.Time
	Description  string
	Venue        string
	Open         bool
	Public       bool
	Managers     []crypto.Token // default é qualquer um do coletivo
}

func (c *CreateEventAction) Serialize() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(c.Epoch, &bytes)
	util.PutToken(c.Author, &bytes)
	util.PutByte(ICreateEventAction, &bytes)
	util.PutString(c.Reasons, &bytes)
	util.PutString(c.OnBehalfOf, &bytes)
	util.PutTime(c.StartAt, &bytes)
	util.PutTime(c.EstimatedEnd, &bytes)
	util.PutString(c.Description, &bytes)
	util.PutString(c.Venue, &bytes)
	util.PutBool(c.Open, &bytes)
	util.PutBool(c.Public, &bytes)
	PutTokenArray(c.Managers, &bytes)
	return bytes
}

func ParseCreateEventAction(create []byte) *CreateEventAction {
	action := CreateEventAction{}
	position := 0
	action.Epoch, position = util.ParseUint64(create, position)
	action.Author, position = util.ParseToken(create, position)
	if create[position] != ICreateEventAction {
		return nil
	}
	position += 1
	action.Reasons, position = util.ParseString(create, position)
	action.OnBehalfOf, position = util.ParseString(create, position)
	action.StartAt, position = util.ParseTime(create, position)
	action.EstimatedEnd, position = util.ParseTime(create, position)
	action.Description, position = util.ParseString(create, position)
	action.Venue, position = util.ParseString(create, position)
	action.Open, position = util.ParseBool(create, position)
	action.Public, position = util.ParseBool(create, position)
	action.Managers, position = ParseTokenArray(create, position)

	if position != len(create) {
		return nil
	}
	return &action
}

type CancelEventAction struct {
	Epoch   uint64
	Author  crypto.Token
	Reasons string
	Hash    crypto.Hash
}

func (c *CancelEventAction) Serialize() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(c.Epoch, &bytes)
	util.PutToken(c.Author, &bytes)
	util.PutByte(ICancelEventAction, &bytes)
	util.PutString(c.Reasons, &bytes)
	util.PutHash(c.Hash, &bytes)
	return bytes
}

func ParseCancelEventAction(create []byte) *CancelEventAction {
	action := CancelEventAction{}
	position := 0
	action.Epoch, position = util.ParseUint64(create, position)
	action.Author, position = util.ParseToken(create, position)
	if create[position] != ICancelEventAction {
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

type UpdateEventAction struct {
	Epoch       uint64
	Author      crypto.Token
	Reasons     string
	EventHash   crypto.Hash
	Description string
	Venue       string
	Open        bool
	Public      bool
	Managers    []crypto.Token
}

func (c *UpdateEventAction) Serialize() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(c.Epoch, &bytes)
	util.PutToken(c.Author, &bytes)
	util.PutByte(IUpdateEventAction, &bytes)
	util.PutString(c.Reasons, &bytes)
	util.PutHash(c.EventHash, &bytes)
	util.PutString(c.Description, &bytes)
	util.PutString(c.Venue, &bytes)
	util.PutBool(c.Open, &bytes)
	util.PutBool(c.Public, &bytes)
	PutTokenArray(c.Managers, &bytes)
	return bytes
}

func ParseUpdateEventAction(create []byte) *UpdateEventAction {
	action := UpdateEventAction{}
	position := 0
	action.Epoch, position = util.ParseUint64(create, position)
	action.Author, position = util.ParseToken(create, position)
	if create[position] != IUpdateEventAction {
		return nil
	}
	position += 1
	action.Reasons, position = util.ParseString(create, position)
	action.EventHash, position = util.ParseHash(create, position)
	action.Description, position = util.ParseString(create, position)
	action.Venue, position = util.ParseString(create, position)
	action.Open, position = util.ParseBool(create, position)
	action.Public, position = util.ParseBool(create, position)
	action.Managers, position = ParseTokenArray(create, position)

	if position != len(create) {
		return nil
	}
	return &action
}

type CheckinEventAction struct {
	Epoch     uint64
	Author    crypto.Token
	Reasons   string
	EventHash crypto.Hash
}

func (c *CheckinEventAction) Serialize() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(c.Epoch, &bytes)
	util.PutToken(c.Author, &bytes)
	util.PutByte(ICheckinEventAction, &bytes)
	util.PutString(c.Reasons, &bytes)
	util.PutHash(c.EventHash, &bytes)
	return bytes
}

func ParseCheckinEventAction(create []byte) *CheckinEventAction {
	action := CheckinEventAction{}
	position := 0
	action.Epoch, position = util.ParseUint64(create, position)
	action.Author, position = util.ParseToken(create, position)
	if create[position] != ICheckinEventAction {
		return nil
	}
	position += 1
	action.Reasons, position = util.ParseString(create, position)
	action.EventHash, position = util.ParseHash(create, position)
	if position != len(create) {
		return nil
	}
	return &action
}

type AcceptCheckinEventAction struct {
	Epoch          uint64
	Author         crypto.Token
	Reasons        string
	Hash           crypto.Hash
	SecretKey      []byte // diffie-hellman
	ContentType    string
	PrivateContent []byte
}

func (c *AcceptCheckinEventAction) Serialize() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(c.Epoch, &bytes)
	util.PutToken(c.Author, &bytes)
	util.PutByte(ICheckinEventAction, &bytes)
	util.PutString(c.Reasons, &bytes)
	util.PutHash(c.Hash, &bytes)
	util.PutByteArray(c.SecretKey, &bytes)
	util.PutString(c.ContentType, &bytes)
	util.PutByteArray(c.PrivateContent, &bytes)
	return bytes
}

func ParseAcceptCheckinEventAction(create []byte) *AcceptCheckinEventAction {
	action := AcceptCheckinEventAction{}
	position := 0
	action.Epoch, position = util.ParseUint64(create, position)
	action.Author, position = util.ParseToken(create, position)
	if create[position] != ICheckinEventAction {
		return nil
	}
	position += 1
	action.Reasons, position = util.ParseString(create, position)
	action.Hash, position = util.ParseHash(create, position)
	action.SecretKey, position = util.ParseByteArray(create, position)
	action.ContentType, position = util.ParseString(create, position)
	action.PrivateContent, position = util.ParseByteArray(create, position)
	if position != len(create) {
		return nil
	}
	return &action
}
