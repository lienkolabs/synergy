package actions

import (
	"time"

	"github.com/lienkolabs/swell/crypto"
	"github.com/lienkolabs/swell/util"
)

type CreateEvent struct {
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

func (c *CreateEvent) Serialize() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(c.Epoch, &bytes)
	util.PutToken(c.Author, &bytes)
	util.PutByte(ACreateEvent, &bytes)
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

func ParseCreateEvent(create []byte) *CreateEvent {
	action := CreateEvent{}
	position := 0
	action.Epoch, position = util.ParseUint64(create, position)
	action.Author, position = util.ParseToken(create, position)
	if create[position] != ACreateEvent {
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

type CancelEvent struct {
	Epoch   uint64
	Author  crypto.Token
	Reasons string
	Hash    crypto.Hash
}

func (c *CancelEvent) Serialize() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(c.Epoch, &bytes)
	util.PutToken(c.Author, &bytes)
	util.PutByte(ACancelEvent, &bytes)
	util.PutString(c.Reasons, &bytes)
	util.PutHash(c.Hash, &bytes)
	return bytes
}

func ParseCancelEvent(create []byte) *CancelEvent {
	action := CancelEvent{}
	position := 0
	action.Epoch, position = util.ParseUint64(create, position)
	action.Author, position = util.ParseToken(create, position)
	if create[position] != ACancelEvent {
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

type UpdateEvent struct {
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

func (c *UpdateEvent) Serialize() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(c.Epoch, &bytes)
	util.PutToken(c.Author, &bytes)
	util.PutByte(AUpdateEvent, &bytes)
	util.PutString(c.Reasons, &bytes)
	util.PutHash(c.EventHash, &bytes)
	util.PutString(c.Description, &bytes)
	util.PutString(c.Venue, &bytes)
	util.PutBool(c.Open, &bytes)
	util.PutBool(c.Public, &bytes)
	PutTokenArray(c.Managers, &bytes)
	return bytes
}

func ParseUpdateEvent(create []byte) *UpdateEvent {
	action := UpdateEvent{}
	position := 0
	action.Epoch, position = util.ParseUint64(create, position)
	action.Author, position = util.ParseToken(create, position)
	if create[position] != AUpdateEvent {
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

type CheckinEvent struct {
	Epoch     uint64
	Author    crypto.Token
	Reasons   string
	EventHash crypto.Hash
}

func (c *CheckinEvent) Serialize() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(c.Epoch, &bytes)
	util.PutToken(c.Author, &bytes)
	util.PutByte(ACheckinEvent, &bytes)
	util.PutString(c.Reasons, &bytes)
	util.PutHash(c.EventHash, &bytes)
	return bytes
}

func ParseCheckinEvent(create []byte) *CheckinEvent {
	action := CheckinEvent{}
	position := 0
	action.Epoch, position = util.ParseUint64(create, position)
	action.Author, position = util.ParseToken(create, position)
	if create[position] != ACheckinEvent {
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

type AcceptCheckinEvent struct {
	Epoch          uint64
	Author         crypto.Token
	Reasons        string
	EventHash      crypto.Hash
	CheckedIn      crypto.Token
	SecretKey      []byte // diffie-hellman
	ContentType    string
	PrivateContent []byte
}

func (c *AcceptCheckinEvent) Serialize() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(c.Epoch, &bytes)
	util.PutToken(c.Author, &bytes)
	util.PutByte(ACheckinEvent, &bytes)
	util.PutString(c.Reasons, &bytes)
	util.PutHash(c.EventHash, &bytes)
	util.PutToken(c.CheckedIn, &bytes)
	util.PutByteArray(c.SecretKey, &bytes)
	util.PutString(c.ContentType, &bytes)
	util.PutByteArray(c.PrivateContent, &bytes)
	return bytes
}

func ParseAcceptCheckinEvent(create []byte) *AcceptCheckinEvent {
	action := AcceptCheckinEvent{}
	position := 0
	action.Epoch, position = util.ParseUint64(create, position)
	action.Author, position = util.ParseToken(create, position)
	if create[position] != ACheckinEvent {
		return nil
	}
	position += 1
	action.Reasons, position = util.ParseString(create, position)
	action.EventHash, position = util.ParseHash(create, position)
	action.CheckedIn, position = util.ParseToken(create, position)
	action.SecretKey, position = util.ParseByteArray(create, position)
	action.ContentType, position = util.ParseString(create, position)
	action.PrivateContent, position = util.ParseByteArray(create, position)
	if position != len(create) {
		return nil
	}
	return &action
}
