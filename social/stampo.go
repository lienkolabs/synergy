package social

import (
	"errors"

	"github.com/lienkolabs/swell/crypto"
)

type ImprintStampInstruction struct {
	Epoch      uint64
	Author     crypto.Token
	OnBehalfOf string
	Reasons    string
	Hash       crypto.Hash
}

type Stamp struct {
	Reputation *Collective
	Release    *Release
	Hash       crypto.Hash
	Votes      []VoteInstruction
	Imprinted  bool
}

func (p *Stamp) IncorporateVote(vote VoteInstruction, state *State) error {
	if vote.Hash != p.Hash {
		return errors.New("invalid hash")
	}
	for _, cast := range p.Votes {
		if cast.Author == vote.Author {
			return errors.New("vote already cast")
		}
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
	p.Imprinted = true
	if p.Release.Stamps == nil {
		p.Release.Stamps = []*Stamp{p}
	} else {
		p.Release.Stamps = append(p.Release.Stamps, p)
	}
	delete(state.Proposals, p.Hash)
	return nil
}

type Release struct {
	Draft    *Draft
	Hash     crypto.Hash // (hash of the original instruction to release)
	Votes    []VoteInstruction
	Released bool
	Stamps   []*Stamp
}

type ReleaseInstruction struct {
	Epoch     uint64
	Author    crypto.Token
	DraftHash crypto.Hash
}

type Post struct {
	Hash    crypto.Hash
	Draft   *Draft
	Journal *Zine
}

func (p *Release) IncorporateVote(vote VoteInstruction, state *State) error {
	if vote.Hash != p.Hash {
		return errors.New("invalid hash")
	}
	for _, cast := range p.Votes {
		if cast.Author == vote.Author {
			return errors.New("vote already cast")
		}
	}
	p.Votes = append(p.Votes, vote)
	if p.Released {
		return nil
	}
	if !p.Draft.Authors.Consensus(p.Hash, p.Votes) {
		return nil
	}
	// new consensus
	p.Released = true
	delete(state.Proposals, p.Hash)
	if _, ok := state.Releases[p.Draft.DraftHash]; !ok {
		state.Releases[p.Draft.DraftHash] = p
		return nil
	}
	return errors.New("already released")
}
