package social

import (
	"errors"

	"github.com/lienkolabs/swell/crypto"
)

type DraftInstruction struct {
	Epoch         uint64
	Author        crypto.Token
	OnBehalfOf    string
	Reasons       string
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

func (d *DraftInstruction) AddDeliberation(vote VoteInstruction, state *State) error {
	return nil
}

type Draft struct {
	Title           string
	Description     string
	Authors         Consensual
	DraftType       string
	DraftHash       crypto.Hash // this must be a valid Media in the state
	PreviousVersion *Draft
	References      []crypto.Hash
	Votes           []VoteInstruction
	Aproved         bool
}

// IncorpoateVote checks if vote scope is correct (hash) if vote was not alrerady
// cast. If new valid vote returns if the new vote is sufficient for consensus
func (d *Draft) IncorporateVote(vote VoteInstruction, state *State) error {
	if vote.Hash != d.DraftHash {
		return errors.New("invalid hash")
	}
	for _, cast := range d.Votes {
		if cast.Author == vote.Author {
			return errors.New("vote already cast")
		}
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
