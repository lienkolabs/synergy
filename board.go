package synergy

import (
	"errors"

	"github.com/lienkolabs/swell/crypto"
)

type PinInstruction struct {
	Epoch         uint64
	Author        crypto.Token
	Board         string
	Draft         crypto.Hash
	Pin           bool
	HashSignature crypto.Signature
}

type BoardEditorInstruction struct {
	Epoch         uint64
	Author        crypto.Token
	Board         string
	Editor        crypto.Token
	Insert        bool
	HashSignature crypto.Signature
}

type Board struct {
	Name           string
	Keyword        []string
	Administrators Consensual
	Editors        *UnamedCollective
	Pinned         []*Draft
}

func (b *Board) Pin(d *Draft) error {
	for _, pinned := range b.Pinned {
		if pinned == d {
			return errors.New("already pinned")
		}
	}
	b.Pinned = append(b.Pinned, d)
	return nil
}

func (b *Board) Remove(d *Draft) error {
	for n, pinned := range b.Pinned {
		if pinned == d {
			b.Pinned = append(b.Pinned[0:n], b.Pinned[n+1:]...)
			return nil
		}
	}
	return errors.New("not pinned")
}

func (b *Board) Last(n int) []*Draft {
	if len(b.Pinned) <= n {
		return b.Pinned
	}
	return b.Pinned[len(b.Pinned)-n:]
}

func (b *Board) First(n int) []*Draft {
	if len(b.Pinned) <= n {
		return b.Pinned
	}
	return b.Pinned[0:n]
}

type BoardPinAction struct {
	Hash  crypto.Hash // hash of combined draft hash + board name + epoch + pin/remove?
	Epoch uint64
	Board *Board
	Draft *Draft
	Pin   bool
	Votes []VoteInstruction
}

func (p *BoardPinAction) IncorporateVote(vote VoteInstruction, state *State) error {
	if vote.Hash != p.Hash {
		return errors.New("invalid hash")
	}
	for _, cast := range p.Votes {
		if cast.Author == vote.Author {
			return errors.New("vote already cast")
		}
	}
	p.Votes = append(p.Votes, vote)
	if p.Board.Editors.Consensus(vote.Hash, p.Votes) {
		delete(state.Proposals, p.Hash)
		if p.Pin {
			return p.Board.Pin(p.Draft)
		}
		return p.Board.Remove(p.Draft)
	}
	return nil
}

type BoardEditorAction struct {
	Hash   crypto.Hash // hash of combined draft hash + board name + epoch + pin/remove?
	Epoch  uint64
	Board  *Board
	Editor crypto.Token
	Insert bool
	Votes  []VoteInstruction
}

func (e *BoardEditorAction) IncorporateVote(vote VoteInstruction, state *State) error {
	if vote.Hash != e.Hash {
		return errors.New("invalid hash")
	}
	for _, cast := range e.Votes {
		if cast.Author == vote.Author {
			return errors.New("vote already cast")
		}
	}
	e.Votes = append(e.Votes, vote)
	if e.Board.Administrators.Consensus(vote.Hash, e.Votes) {
		delete(state.Proposals, vote.Hash)
		if e.Insert {
			e.Board.Editors.IncludeMember(e.Editor)
		} else {
			e.Board.Editors.RemoveMember(e.Editor)
		}
		return nil
	}
	return nil
}
