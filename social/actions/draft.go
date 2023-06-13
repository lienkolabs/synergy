package actions

import (
	"github.com/lienkolabs/swell/crypto"
	"github.com/lienkolabs/swell/util"
)

type DraftAction struct {
	Epoch         uint64
	Author        crypto.Token
	Reasons       string
	OnBehalfOf    string
	CoAuthors     []crypto.Token
	Policy        *Policy
	Title         string
	Keywords      []string
	Description   string
	ContentType   string
	ContentHash   crypto.Hash // hash of the entire content, not of the part
	NumberOfParts byte
	Content       []byte // entire content of the first part
	PreviousDraft crypto.Hash
	References    []crypto.Hash
}

func (c *DraftAction) Serialize() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(c.Epoch, &bytes)
	util.PutToken(c.Author, &bytes)
	util.PutByte(IDraftAction, &bytes)
	util.PutString(c.Reasons, &bytes)
	util.PutString(c.OnBehalfOf, &bytes)
	PutTokenArray(c.CoAuthors, &bytes)
	PutOptionalPolicy(c.Policy, &bytes)
	util.PutString(c.Title, &bytes)
	PutKeywords(c.Keywords, &bytes)
	util.PutString(c.Description, &bytes)
	util.PutString(c.ContentType, &bytes)
	util.PutHash(c.ContentHash, &bytes)
	util.PutByte(c.NumberOfParts, &bytes)
	util.PutByteArray(c.Content, &bytes)
	util.PutHash(c.PreviousDraft, &bytes)
	PutHashArray(c.References, &bytes)
	return bytes
}

func ParseDraftAction(create []byte) *DraftAction {
	action := DraftAction{}
	position := 0
	action.Epoch, position = util.ParseUint64(create, position)
	action.Author, position = util.ParseToken(create, position)
	if create[position] != IDraftAction {
		return nil
	}
	position += 1
	action.Reasons, position = util.ParseString(create, position)
	action.OnBehalfOf, position = util.ParseString(create, position)
	action.CoAuthors, position = ParseTokenArray(create, position)
	action.Title, position = util.ParseString(create, position)
	action.Keywords, position = ParseKeywords(create, position)
	action.Description, position = util.ParseString(create, position)
	action.ContentType, position = util.ParseString(create, position)
	action.ContentHash, position = util.ParseHash(create, position)
	action.NumberOfParts, position = util.ParseByte(create, position)
	action.Content, position = util.ParseByteArray(create, position)
	action.PreviousDraft, position = util.ParseHash(create, position)
	action.References, position = ParseHashArray(create, position)
	if position != len(create) {
		return nil
	}
	return &action
}

type ReleaseDraftAction struct {
	Epoch       uint64
	Author      crypto.Token
	Reasons     string
	ContentHash crypto.Hash
}

func (c *ReleaseDraftAction) Serialize() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(c.Epoch, &bytes)
	util.PutToken(c.Author, &bytes)
	util.PutByte(IReleaseDraftAction, &bytes)
	util.PutString(c.Reasons, &bytes)
	util.PutHash(c.ContentHash, &bytes)
	return bytes
}

func ParseReleaseDraftAction(create []byte) *ReleaseDraftAction {
	action := ReleaseDraftAction{}
	position := 0
	action.Epoch, position = util.ParseUint64(create, position)
	action.Author, position = util.ParseToken(create, position)
	if create[position] != IReleaseDraftAction {
		return nil
	}
	position += 1
	action.Reasons, position = util.ParseString(create, position)
	action.ContentHash, position = util.ParseHash(create, position)
	if position != len(create) {
		return nil
	}
	return &action
}
