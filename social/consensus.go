package social

import (
	"errors"

	"github.com/lienkolabs/swell/crypto"
)

// % + 1 Vote...or 100%
// Supermahority to change policy rule
// Majority for anything else
type Policy struct {
	Majority      int
	SuperMajority int
}

type VoteInstruction struct {
	Epoch     uint64
	Author    crypto.Token
	Reasons   string
	Hash      crypto.Hash
	Approve   bool
	Signature crypto.Signature
}

type Consensual interface {
	Consensus(hash crypto.Hash, votes []VoteInstruction) bool
	IsMember(token crypto.Token) bool
	IncludeMember(token crypto.Token)
	RemoveMember(token crypto.Token)
	ChangeMajority(majority int)
	ListOfMembers() map[crypto.Token]struct{}
}

func consensus(members map[crypto.Token]struct{}, votesRequired int, hash crypto.Hash, votes []VoteInstruction) bool {
	count := 0
	for _, vote := range votes {
		_, isMember := members[vote.Author]
		if isMember && vote.Approve && hash == vote.Hash && vote.Author.Verify(vote.Hash[:], vote.Signature) {
			count += 1
			if count >= votesRequired {
				return true
			}
		}
	}
	return false
}

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
