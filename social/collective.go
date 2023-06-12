package social

import (
	"errors"

	"github.com/lienkolabs/swell/crypto"
)

type CreateCollectiveInstruction struct {
	Epoch       uint64
	Author      crypto.Token
	Reasons     string
	Name        string
	Description string
	Policy      Policy
}

func (c *CreateCollectiveInstruction) Serialize() []byte {
	return nil
}

type UpdateCollectiveInstruction struct {
	Epoch       uint64
	Author      crypto.Token
	OnBehalfOf  string
	Reasons     string
	Description string
	Policy      *Policy
}

func (u *UpdateCollectiveInstruction) Serialize() []byte {
	return nil
}

type RequestMembershipInstruction struct {
	Epoch      uint64
	Author     crypto.Token
	OnBehalfOf string
	Reasons    string
	Include    bool // true = author as new member, false = self exclude
}

func (r *RequestMembershipInstruction) Serialize() []byte {
	return nil
}

type RemoveMemberInstruction struct {
	Epoch      uint64
	Author     crypto.Token
	OnBehalfOf string
	Reasons    string
	Member     crypto.Token
}

func (r *RemoveMemberInstruction) Serialize() []byte {
	return nil
}

type Collective struct {
	Name        string
	Members     map[crypto.Token]struct{}
	Description string
	Policy      Policy
}

func (c *Collective) ListOfMembers() map[crypto.Token]struct{} {
	return nil
}

func (c *Collective) Photo() *Collective {
	cloned := Collective{
		Name:    c.Name,
		Members: make(map[crypto.Token]struct{}),
		Policy: Policy{
			Majority:      c.Policy.Majority,
			SuperMajority: c.Policy.SuperMajority,
		},
	}
	for member, _ := range c.Members {
		cloned.Members[member] = struct{}{}
	}
	return &cloned
}

func (c *Collective) IncludeMember(token crypto.Token) {
	c.Members[token] = struct{}{}
}

func (c *Collective) RemoveMember(token crypto.Token) {
	delete(c.Members, token)
}

func (c *Collective) ChangeMajority(majority int) {
	c.Policy.Majority = majority
}

func (c *Collective) Consensus(hash crypto.Hash, votes []VoteInstruction) bool {
	required := len(c.Members)*c.Policy.Majority/100 + 1
	if required > len(c.Members) {
		required = len(c.Members)
	}
	return consensus(c.Members, required, hash, votes)
}

func (c *Collective) SuperConsensus(hash crypto.Hash, votes []VoteInstruction) bool {
	required := len(c.Members)*c.Policy.SuperMajority/100 + 1
	if required > len(c.Members) {
		required = len(c.Members)
	}
	return consensus(c.Members, required, hash, votes)
}

func (c *Collective) IsMember(token crypto.Token) bool {
	_, ok := c.Members[token]
	return ok
}

type UnamedCollective struct {
	Members  map[crypto.Token]struct{}
	Majority int
}

func (c *UnamedCollective) ListOfMembers() map[crypto.Token]struct{} {
	return c.Members
}

func (c *UnamedCollective) Consensus(hash crypto.Hash, votes []VoteInstruction) bool {
	required := len(c.Members)*c.Majority/100 + 1
	if required > len(c.Members) {
		required = len(c.Members)
	}
	return consensus(c.Members, required, hash, votes)
}

func (c *UnamedCollective) IsMember(token crypto.Token) bool {
	_, ok := c.Members[token]
	return ok
}

func (c *UnamedCollective) IncludeMember(token crypto.Token) {
	c.Members[token] = struct{}{}
}

func (c *UnamedCollective) RemoveMember(token crypto.Token) {
	delete(c.Members, token)
}

func (c *UnamedCollective) ChangeMajority(majority int) {
	c.Majority = majority
}

type PendingUpdate struct {
	Update       *UpdateCollectiveInstruction
	Collective   *Collective
	Hash         crypto.Hash
	ChangePolicy bool
	Votes        []VoteInstruction
}

func (p *PendingUpdate) IncorporateVote(vote VoteInstruction, state *State) error {
	if err := isValidVote(p.Hash, vote, p.Votes); err != nil {
		return err
	}
	p.Votes = append(p.Votes, vote)
	if p.ChangePolicy {
		if !p.Collective.SuperConsensus(p.Hash, p.Votes) {
			return nil
		}
	} else {
		if !p.Collective.Consensus(p.Hash, p.Votes) {
			return nil
		}
	}
	// exclude pending update from live proposals because of consensus
	delete(state.Proposals, p.Hash)
	// update collective
	if p.Update.Description != "" {
		p.Collective.Description = p.Update.Description
	}
	if p.ChangePolicy {
		p.Collective.Policy = Policy{
			Majority:      p.Update.Policy.Majority,
			SuperMajority: p.Update.Policy.SuperMajority,
		}
	}
	return nil
}

type PendingRequestMembership struct {
	Request    *RequestMembershipInstruction
	Collective *Collective
	Hash       crypto.Hash
	Votes      []VoteInstruction
}

func (p *PendingRequestMembership) IncorporateVote(vote VoteInstruction, state *State) error {
	if err := isValidVote(p.Hash, vote, p.Votes); err != nil {
		return err
	}
	p.Votes = append(p.Votes, vote)
	if !p.Collective.Consensus(p.Hash, p.Votes) {
		return nil
	}
	delete(state.Proposals, p.Hash)
	collective, ok := state.Collectives[p.Collective.Name]
	if !ok {
		return errors.New("collective not found")
	}
	collective.Members[p.Request.Author] = struct{}{}
	return nil
}

type PendingRemoveMember struct {
	Remove     *RemoveMemberInstruction
	Collective *Collective
	Hash       crypto.Hash
	Votes      []VoteInstruction
}

func (p *PendingRemoveMember) IncorporateVote(vote VoteInstruction, state *State) error {
	if err := isValidVote(p.Hash, vote, p.Votes); err != nil {
		return err
	}
	p.Votes = append(p.Votes, vote)
	if !p.Collective.Consensus(p.Hash, p.Votes) {
		return nil
	}
	delete(state.Proposals, p.Hash)
	collective, ok := state.Collectives[p.Collective.Name]
	if !ok {
		return errors.New("collective not found")
	}
	delete(collective.Members, p.Remove.Member)
	return nil
}
