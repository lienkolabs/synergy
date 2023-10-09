package state

import (
	"github.com/lienkolabs/breeze/crypto"
	"github.com/lienkolabs/synergy/social/actions"
)

type Draft struct {
	Title           string
	Date            uint64
	Description     string
	Authors         Consensual
	DraftType       string
	DraftHash       crypto.Hash // this must be a valid Media in the state
	PreviousVersion *Draft
	Keywords        []string
	References      []crypto.Hash
	Votes           []actions.Vote
	Pinned          []*Board
	Edits           []*Edit
	Aproved         bool
}

// RULES:
// ======
// If there is no previous version:
// If on behalf of collective, collective must consent
// If with co-authors every co-author must consent
//
// If there is a previous version
// If new version is with co-authors, every new co-author must consent
// Previous authors must collectively consent according to policy
// Current authors must collectively consent according to policy
func (d *Draft) Consensus() bool {
	if d.Aproved {
		return true
	}

	if !d.Authors.Consensus(d.DraftHash, d.Votes) {
		return false
	}
	previous := d.PreviousVersion
	if previous == nil {
		if d.Authors.CollectiveName() == "" {
			// every co-author must vote
			return d.Authors.Unanimous(d.DraftHash, d.Votes)
		}
		// it is a collective with consensus formed
		return true
	}
	if !previous.Authors.Consensus(d.DraftHash, d.Votes) {
		return false
	}
	collective := d.Authors.CollectiveName()
	if collective != "" {
		// current version is collective and previous version consent, its ok
		return true
	}
	previousCollective := previous.Authors.CollectiveName()
	if previousCollective != "" {
		// if previous is a collective and current is not... every co-author must
		// sign
		return d.Authors.Unanimous(d.DraftHash, d.Votes)
	}
	// here we know that neither current nor previous version is a collective
	currentMembers := d.Authors.ListOfMembers()
	previousMembers := previous.Authors.ListOfMembers()
	newMembers := make(map[crypto.Token]struct{})
	for token, _ := range currentMembers {
		if _, ok := previousMembers[token]; !ok {
			newMembers[token] = struct{}{}
		}
	}
	for _, vote := range d.Votes {
		if _, ok := newMembers[vote.Author]; ok && vote.Hash == d.DraftHash {
			delete(newMembers, vote.Author)
		}
	}
	return len(newMembers) == 0
}

// IncorpoateVote checks if vote scope is correct (hash) if vote was not alrerady
// cast. If new valid vote returns if the new vote is sufficient for consensus
func (d *Draft) IncorporateVote(vote actions.Vote, state *State) error {
	if err := IsNewValidVote(vote, d.Votes, d.DraftHash); err != nil {
		return err
	}
	d.Votes = append(d.Votes, vote)
	if d.Aproved {
		return nil
	}
	consensus := d.Consensus()
	if consensus {
		d.Aproved = true
		state.Proposals.Delete(d.DraftHash)
		state.Drafts[d.DraftHash] = d
		state.IndexConsensus(d.DraftHash, true)
	}
	return nil
}
