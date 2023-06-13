package state

import (
	"errors"

	"github.com/lienkolabs/swell/crypto"
	"github.com/lienkolabs/synergy/social/actions"
)

type Board struct {
	Name       string
	Keyword    []string
	Collective *Collective
	Editors    *UnamedCollective
	Pinned     []*Draft
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

type Pin struct {
	Hash  crypto.Hash // hash of original instruction
	Epoch uint64
	Board *Board
	Draft *Draft
	Pin   bool
	Votes []actions.VoteAction
}

func (p *Pin) IncorporateVote(vote actions.VoteAction, state *State) error {
	IsNewValidVote(vote, p.Votes, p.Hash)
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

type BoardEditor struct {
	Hash   crypto.Hash // hash of combined draft hash + board name + epoch + pin/remove?
	Epoch  uint64
	Board  *Board
	Editor crypto.Token
	Insert bool
	Votes  []actions.VoteAction
}

func (e *BoardEditor) IncorporateVote(vote actions.VoteAction, state *State) error {
	IsNewValidVote(vote, e.Votes, e.Hash)
	e.Votes = append(e.Votes, vote)
	if e.Board.Collective.Consensus(vote.Hash, e.Votes) {
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
