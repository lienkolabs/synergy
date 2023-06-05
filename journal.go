package synergy

import (
	"errors"

	"github.com/lienkolabs/swell/crypto"
)

type JournalEditorInstruction struct {
	Epoch         uint64
	Author        crypto.Token
	Journal       string
	Editor        crypto.Token
	Insert        bool
	HashSignature crypto.Signature
}

type Journal struct {
	Name           string
	Keyword        []string
	Administrators Consensual
	Editors        *UnamedCollective
	Published      []*Post
}

func (b *Journal) Last(n int) []*Post {
	if len(b.Published) <= n {
		return b.Published
	}
	return b.Published[len(b.Published)-n:]
}

func (b *Journal) First(n int) []*Post {
	if len(b.Published) <= n {
		return b.Published
	}
	return b.Published[0:n]
}

type JournalEditorAction struct {
	Hash   crypto.Hash // hash of combined draft hash + board name + epoch + pin/remove?
	Epoch  uint64
	Board  *Board
	Editor crypto.Token
	Insert bool
	Votes  []VoteInstruction
}

func (e *JournalEditorAction) IncorporateVote(vote VoteInstruction, state *State) error {
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
