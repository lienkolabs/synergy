package state

import (
	"github.com/lienkolabs/swell/crypto"
	"github.com/lienkolabs/synergy/social/actions"
)

type Edit struct {
	Authors  Consensual
	Reasons  string
	Draft    *Draft
	EditType string
	Edit     crypto.Hash
	Votes    []actions.Vote
}

func (e *Edit) IncorporateVote(vote actions.Vote, state *State) error {
	if err := IsNewValidVote(vote, e.Votes, e.Edit); err != nil {
		return err
	}
	e.Votes = append(e.Votes, vote)
	if !e.Authors.Consensus(e.Edit, e.Votes) {
		return nil
	}
	// new consensus
	delete(state.Proposals, e.Edit)
	// to do where to put edits?
	return nil
}
