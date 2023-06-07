package synergy

import (
	"errors"

	"github.com/lienkolabs/swell/crypto"
)

// Instruction to require consensus over an inclusion/exlusion of editor from a
// Zine.
type ZinelEditorInstruction struct {
	Epoch         uint64
	Author        crypto.Token
	Zine          string
	Editor        crypto.Token
	Insert        bool
	HashSignature crypto.Signature
}

type LaunchZineInstruction struct {
	Epoch        uint64
	Author       crypto.Token
	Reasons      string
	Name         string
	BoardPolicy  Policy // Majority = Elect Editors, Supermajority change Board or policy
	EditorPolicy Policy // Majority = Supermajority = Publish draft
}

type Zine struct {
	Name           string
	Keyword        []string
	Administrators Consensual
	Editors        *UnamedCollective
	Published      []*Post
}

func (b *Zine) Last(n int) []*Post {
	if len(b.Published) <= n {
		return b.Published
	}
	return b.Published[len(b.Published)-n:]
}

func (b *Zine) First(n int) []*Post {
	if len(b.Published) <= n {
		return b.Published
	}
	return b.Published[0:n]
}

type ZineEditorAction struct {
	Hash   crypto.Hash // hash of combined draft hash + board name + epoch + pin/remove?
	Epoch  uint64
	Zine   *Zine
	Editor crypto.Token
	Insert bool
	Votes  []VoteInstruction
}

func (e *ZineEditorAction) IncorporateVote(vote VoteInstruction, state *State) error {
	if vote.Hash != e.Hash {
		return errors.New("invalid hash")
	}
	for _, cast := range e.Votes {
		if cast.Author == vote.Author {
			return errors.New("vote already cast")
		}
	}
	e.Votes = append(e.Votes, vote)
	if e.Zine.Administrators.Consensus(vote.Hash, e.Votes) {
		delete(state.Proposals, vote.Hash)
		if e.Insert {
			e.Zine.Editors.IncludeMember(e.Editor)
		} else {
			e.Zine.Editors.RemoveMember(e.Editor)
		}
		return nil
	}
	return nil
}
