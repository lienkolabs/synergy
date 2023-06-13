package actions

import (
	"github.com/lienkolabs/swell/crypto"
	"github.com/lienkolabs/swell/util"
)

type CreateCollective struct {
	Epoch       uint64
	Author      crypto.Token
	Reasons     string
	Name        string
	Description string
	Policy      Policy
}

func (c *CreateCollective) Serialize() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(c.Epoch, &bytes)
	util.PutToken(c.Author, &bytes)
	util.PutByte(ICreateCollectiveAction, &bytes)
	util.PutString(c.Reasons, &bytes)
	util.PutString(c.Name, &bytes)
	util.PutString(c.Description, &bytes)
	PutPolicy(c.Policy, &bytes)
	return bytes
}

func ParseCreateCollective(create []byte) *CreateCollective {
	action := CreateCollective{}
	position := 0
	action.Epoch, position = util.ParseUint64(create, position)
	action.Author, position = util.ParseToken(create, position)
	if create[position] != ICreateCollectiveAction {
		return nil
	}
	position += 1
	action.Reasons, position = util.ParseString(create, position)
	action.Name, position = util.ParseString(create, position)
	action.Description, position = util.ParseString(create, position)
	action.Policy, position = ParsePolicy(create, position)
	if position != len(create) {
		return nil
	}
	return &action
}

type UpdateCollective struct {
	Epoch       uint64
	Author      crypto.Token
	Reasons     string
	OnBehalfOf  string
	Description string
	Policy      *Policy
}

func (c *UpdateCollective) Serialize() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(c.Epoch, &bytes)
	util.PutToken(c.Author, &bytes)
	util.PutByte(IUpdateCollectiveAction, &bytes)
	util.PutString(c.Reasons, &bytes)
	util.PutString(c.OnBehalfOf, &bytes)
	util.PutString(c.Description, &bytes)
	if c.Policy != nil {
		util.PutByte(1, &bytes) // there is a policy
		PutPolicy(*c.Policy, &bytes)
	} else {
		util.PutByte(0, &bytes) // there is no policy
	}
	return bytes
}

func ParseUpdateCollective(update []byte) *UpdateCollective {
	action := UpdateCollective{}
	position := 0
	action.Epoch, position = util.ParseUint64(update, position)
	action.Author, position = util.ParseToken(update, position)
	if update[position] != IUpdateCollectiveAction {
		return nil
	}
	position += 1
	action.Reasons, position = util.ParseString(update, position)
	action.OnBehalfOf, position = util.ParseString(update, position)
	action.Description, position = util.ParseString(update, position)
	if update[position] == 1 {
		var policy Policy
		policy, position = ParsePolicy(update, position)
		action.Policy = &policy
	} else if update[position] != 0 {
		return nil
	}
	position += 1
	if position != len(update) {
		return nil
	}
	return &action
}

type RequestMembership struct {
	Epoch      uint64
	Author     crypto.Token
	Reasons    string
	Collective string
	Include    bool
}

func (c *RequestMembership) Serialize() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(c.Epoch, &bytes)
	util.PutToken(c.Author, &bytes)
	util.PutByte(IRequestMembershipAction, &bytes)
	util.PutString(c.Reasons, &bytes)
	util.PutString(c.Collective, &bytes)
	util.PutBool(c.Include, &bytes)
	return bytes
}

func ParseRequestMembership(update []byte) *RequestMembership {
	action := RequestMembership{}
	position := 0
	action.Epoch, position = util.ParseUint64(update, position)
	action.Author, position = util.ParseToken(update, position)
	if update[position] != IRequestMembershipAction {
		return nil
	}
	position += 1
	action.Reasons, position = util.ParseString(update, position)
	action.Collective, position = util.ParseString(update, position)
	action.Include, position = util.ParseBool(update, position)
	if position != len(update) {
		return nil
	}
	return &action
}

type RemoveMember struct {
	Epoch      uint64
	Author     crypto.Token
	OnBehalfOf string
	Reasons    string
	Member     crypto.Token
}

func (c *RemoveMember) Serialize() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(c.Epoch, &bytes)
	util.PutToken(c.Author, &bytes)
	util.PutByte(IRemoveMemberAction, &bytes)
	util.PutString(c.Reasons, &bytes)
	util.PutString(c.OnBehalfOf, &bytes)
	util.PutToken(c.Member, &bytes)
	return bytes
}

func ParseRemoveMember(update []byte) *RemoveMember {
	action := RemoveMember{}
	position := 0
	action.Epoch, position = util.ParseUint64(update, position)
	action.Author, position = util.ParseToken(update, position)
	if update[position] != IRemoveMemberAction {
		return nil
	}
	position += 1
	action.Reasons, position = util.ParseString(update, position)
	action.OnBehalfOf, position = util.ParseString(update, position)
	action.Member, position = util.ParseToken(update, position)
	if position != len(update) {
		return nil
	}
	return &action
}
