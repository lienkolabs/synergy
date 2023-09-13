package state

import (
	"errors"

	"github.com/lienkolabs/breeze/crypto"
	"github.com/lienkolabs/synergy/social/actions"
)

type Stamp struct {
	Reputation *Collective
	Release    *Release
	Hash       crypto.Hash
	Votes      []actions.Vote
	Imprinted  bool
}

func (p *Stamp) IncorporateVote(vote actions.Vote, state *State) error {
	if err := IsNewValidVote(vote, p.Votes, p.Hash); err != nil {
		return err
	}
	if !p.Reputation.IsMember(vote.Author) {
		return errors.New("author is not a recognized member of the collective")
	}
	p.Votes = append(p.Votes, vote)
	if p.Imprinted {
		return nil
	}
	if !p.Reputation.Consensus(vote.Hash, p.Votes) {
		return nil
	}
	// new consensus
	state.IndexConsensus(vote.Hash, true)
	p.Imprinted = true
	//if state.index != nil {
	//	state.index.AddStampToCollective(p, p.Reputation)
	//}
	if p.Release.Stamps == nil {
		p.Release.Stamps = []*Stamp{p}
	} else {
		p.Release.Stamps = append(p.Release.Stamps, p)
	}
	state.Proposals.Delete(p.Hash)
	return nil
}

type Release struct {
	Epoch    uint64
	Draft    *Draft
	Hash     crypto.Hash // (hash of the original instruction to release)
	Votes    []actions.Vote
	Released bool
	Stamps   []*Stamp
}

func (p *Release) IncorporateVote(vote actions.Vote, state *State) error {
	if err := IsNewValidVote(vote, p.Votes, p.Hash); err != nil {
		return err
	}
	p.Votes = append(p.Votes, vote)
	if p.Released {
		return nil
	}
	if !p.Draft.Authors.Consensus(p.Hash, p.Votes) {
		return nil
	}
	// new consensus
	state.IndexConsensus(vote.Hash, true)
	p.Released = true
	state.Proposals.Delete(p.Hash)
	if _, ok := state.Releases[p.Draft.DraftHash]; !ok {
		state.Releases[p.Draft.DraftHash] = p
		return nil
	}
	return errors.New("already released")
}
