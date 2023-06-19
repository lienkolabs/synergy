package state

import (
	"errors"

	"github.com/lienkolabs/swell/crypto"
	"github.com/lienkolabs/synergy/social/actions"
)

type PendingUpdateBoard struct {
	Keywords    []string
	PinMajority int
	Description string
	Board       *Board
	Hash        crypto.Hash
	Votes       []actions.Vote
}

func (b *PendingUpdateBoard) IncorporateVote(vote actions.Vote, state *State) error {
	IsNewValidVote(vote, b.Votes, b.Hash)
	b.Votes = append(b.Votes, vote)
	if b.Board.Collective.Consensus(vote.Hash, b.Votes) {
		state.Proposals.Delete(b.Hash)
		b.Board.Editors.ChangeMajority(b.PinMajority)
		b.Board.Keyword = b.Keywords
		b.Board.Description = b.Description
	}
	return nil
}

type PendingBoard struct {
	Board *Board
	Hash  crypto.Hash
	Votes []actions.Vote
}

func (b *PendingBoard) IncorporateVote(vote actions.Vote, state *State) error {
	IsNewValidVote(vote, b.Votes, b.Hash)
	b.Votes = append(b.Votes, vote)
	if b.Board.Collective.Consensus(vote.Hash, b.Votes) {
		state.Proposals.Delete(b.Hash)
		hash := crypto.Hasher([]byte(b.Board.Name))
		if _, ok := state.Boards[hash]; ok {
			return errors.New("board already exists")
		}
		state.Boards[hash] = b.Board
	}
	return nil
}

type Board struct {
	Name        string
	Keyword     []string
	Description string
	Collective  *Collective
	Editors     *UnamedCollective
	Pinned      []*Draft
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
	Votes []actions.Vote
}

func (p *Pin) IncorporateVote(vote actions.Vote, state *State) error {
	IsNewValidVote(vote, p.Votes, p.Hash)
	p.Votes = append(p.Votes, vote)
	if p.Board.Editors.Consensus(vote.Hash, p.Votes) {
		state.Proposals.Delete(p.Hash)
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
	Votes  []actions.Vote
}

func (e *BoardEditor) IncorporateVote(vote actions.Vote, state *State) error {
	IsNewValidVote(vote, e.Votes, e.Hash)
	e.Votes = append(e.Votes, vote)
	if e.Board.Collective.Consensus(vote.Hash, e.Votes) {
		state.Proposals.Delete(vote.Hash)
		if e.Insert {
			e.Board.Editors.IncludeMember(e.Editor)
		} else {
			e.Board.Editors.RemoveMember(e.Editor)
		}
		return nil
	}
	return nil
}
