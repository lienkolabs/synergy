package state

import (
	"errors"

	"github.com/lienkolabs/swell/crypto"
	"github.com/lienkolabs/synergy/social/actions"
)

type Consensual interface {
	Consensus(hash crypto.Hash, votes []actions.Vote) bool
	IsMember(token crypto.Token) bool
	IncludeMember(token crypto.Token)
	RemoveMember(token crypto.Token)
	ChangeMajority(majority int)
	ListOfMembers() map[crypto.Token]struct{}
	ListOfTokens() map[crypto.Token]struct{}
	CollectiveName() string
	GetPolicy() (majority int, supermajority int)
	Unanimous(hash crypto.Hash, votes []actions.Vote) bool
}

func consensus(members map[crypto.Token]struct{}, votesRequired int, hash crypto.Hash, votes []actions.Vote) bool {
	count := 0
	for _, vote := range votes {
		_, isMember := members[vote.Author]
		if isMember && vote.Approve && hash == vote.Hash {
			count += 1
			if count >= votesRequired {
				return true
			}
		}
	}
	return false
}

func isValidVote(hash crypto.Hash, vote actions.Vote, signatures []actions.Vote) error {
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
