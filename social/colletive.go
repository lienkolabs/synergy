package social

import (
	"errors"

	"github.com/lienkolabs/swell/crypto"
)

func isValidVote(hash crypto.Hash, vote VoteInstruction, signatures []VoteInstruction) error {
	if vote.Hash != hash {
		return errors.New("invalid hash")
	}
	for _, cast := range signatures {
		if cast.Author == vote.Author {
			return errors.New("vote already cast")
		}
	}
	return nil
}

type Consensual interface {
	Consensus(hash crypto.Hash, votes []VoteInstruction) bool
	IsMember(token crypto.Token) bool
	IncludeMember(token crypto.Token)
	RemoveMember(token crypto.Token)
	ChangeMajority(majority int)
}

func consensus(members map[crypto.Token]struct{}, majority int, hash crypto.Hash, votes []VoteInstruction) bool {
	count := 0
	for _, vote := range votes {
		_, isMember := members[vote.Author]
		if isMember && vote.Approve && hash == vote.Hash && vote.Author.Verify(vote.Hash[:], vote.Signature) {
			count += 1
			if count >= majority {
				return true
			}
		}
	}
	return false
}

type Collective struct {
	Name          string // optional
	Members       map[crypto.Token]struct{}
	Majority      int // number of signatures to approve
	SuperMajority int
}

func (c *Collective) IncludeMember(token crypto.Token) {
	c.Members[token] = struct{}{}
}

func (c *Collective) RemoveMember(token crypto.Token) {
	delete(c.Members, token)
}

func (c *Collective) ChangeMajority(majority int) {
	c.Majority = majority
}

func (c *Collective) Consensus(hash crypto.Hash, votes []VoteInstruction) bool {
	return consensus(c.Members, c.Majority, hash, votes)
}

func (c *Collective) SuperConsensus(hash crypto.Hash, votes []VoteInstruction) bool {
	return consensus(c.Members, c.SuperMajority, hash, votes)
}

func (c *Collective) IsMember(token crypto.Token) bool {
	_, ok := c.Members[token]
	return ok
}

type UnamedCollective struct {
	Members  map[crypto.Token]struct{}
	Majority int
}

func (c *UnamedCollective) Consensus(hash crypto.Hash, votes []VoteInstruction) bool {
	return consensus(c.Members, c.Majority, hash, votes)
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

type VoteInstruction struct {
	Epoch     uint64
	Author    crypto.Token
	Reasons   string
	Hash      crypto.Hash
	Approve   bool
	Signature crypto.Signature
}

func Authors(majority int, tokens ...crypto.Token) Consensual {
	collective := UnamedCollective{
		Members:  make(map[crypto.Token]struct{}),
		Majority: majority,
	}
	for _, token := range tokens {
		collective.Members[token] = struct{}{}
	}
	return &collective
}
