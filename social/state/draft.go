package state

import (
	"github.com/lienkolabs/swell/crypto"
	"github.com/lienkolabs/synergy/social/actions"
)

type Draft struct {
	Title           string
	Description     string
	Authors         Consensual
	DraftType       string
	DraftHash       crypto.Hash // this must be a valid Media in the state
	PreviousVersion *Draft
	References      []crypto.Hash
	Votes           []actions.VoteAction
	Aproved         bool
}

// IncorpoateVote checks if vote scope is correct (hash) if vote was not alrerady
// cast. If new valid vote returns if the new vote is sufficient for consensus
func (d *Draft) IncorporateVote(vote actions.VoteAction, state *State) error {
	if err := IsNewValidVote(vote, d.Votes, d.DraftHash); err != nil {
		return err
	}
	d.Votes = append(d.Votes, vote)
	if d.Aproved {
		return nil
	}
	if d.PreviousVersion != nil && !d.PreviousVersion.Authors.Consensus(d.DraftHash, d.Votes) {
		return nil
	}
	newMembers := d.Authors.ListOfMembers()
	if newMembers == nil {
		if d.Authors.Consensus(d.DraftHash, d.Votes) {
			d.Aproved = true
			state.Drafts[d.DraftHash] = d
			return nil
		}
	} else {
		if d.PreviousVersion != nil {
			previous := d.PreviousVersion.Authors.ListOfMembers()
			if previous != nil { // previous is not collective
				for member, _ := range newMembers {
					if _, ok := previous[member]; ok {
						delete(newMembers, member)
					}
				}
			}
		}
	}
	if newMembers != nil {
		for _, vote := range d.Votes {
			delete(newMembers, vote.Author)
		}
	}
	if newMembers == nil || len(newMembers) == 0 {
		d.Aproved = true
		state.Drafts[d.DraftHash] = d
	}
	return nil
}
