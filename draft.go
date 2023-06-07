package synergy

import (
	"errors"

	"github.com/lienkolabs/swell/crypto"
)

// % + 1 Vote...or 100%
// Supermahority to change policy rule
// Majority for anything else
type Policy struct {
	Majority      int
	SuperMajority int
}

type AlternativeDraftInstruction struct {
	Epoch         uint64
	Author        crypto.Token
	OnBehalfOf    string // collective if any
	CoAuthors     []crypto.Token
	Policy        Policy
	Reasons       string
	Title         string
	Keywords      []string
	Description   string
	ContentHash   crypto.Hash
	ContentType   string
	NumberOfParts byte
	Content       []byte
	PreviousDraft crypto.Hash
	References    []crypto.Hash
}

func (d *AlternativeDraftInstruction) AddDeliberation(vote VoteInstruction, state *State) error {
	return nil
}

type DraftInstruction struct {
	Epoch              uint64
	InstructionAuthor  crypto.Token
	OnBehalfOf         string // collective if any
	Reasons            string
	DraftAuthors       Consensual // if nil then collective = instruction author
	DraftTitle         string
	DraftAbstract      string
	DraftType          string
	DraftHash          crypto.Hash // this must be a valid Media in the state
	PreviousVersion    crypto.Hash
	InternalReferences []crypto.Hash // optional list of other synergy content
	HashSignature      crypto.Signature
}

type Draft struct {
	Title              string
	Abstract           string
	Authors            Consensual
	DraftType          string
	DraftHash          crypto.Hash // this must be a valid Media in the state
	PreviousVersion    *Draft
	InternalReferences []crypto.Hash
	Votes              []VoteInstruction
	Aproved            bool
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
	if !d.Authors.Consensus(d.DraftHash, d.Votes) {
		return nil
	}
	if d.PreviousVersion != nil {
		if !d.PreviousVersion.Authors.Consensus(d.DraftHash, d.Votes) {
			return nil
		}
	}
	// new consensus
	d.Aproved = true
	delete(state.Proposals, d.DraftHash)
	state.Drafts[d.DraftHash] = d
	return nil
}

type PostInstruction struct {
	Epoch             uint64
	InstructionAuthor crypto.Token
	Draft             crypto.Hash
	Journal           string
	HashSignature     crypto.Signature // hash of epoch, draft hash + journal
}

type Post struct {
	Hash      crypto.Hash
	Draft     *Draft
	Journal   *Zine
	Votes     []VoteInstruction
	Published bool
}

func (p *Post) IncorporateVote(vote VoteInstruction, state *State) error {
	if vote.Hash != p.Hash {
		return errors.New("invalid hash")
	}
	for _, cast := range p.Votes {
		if cast.Author == vote.Author {
			return errors.New("vote already cast")
		}
	}
	p.Votes = append(p.Votes, vote)
	if p.Published {
		return nil
	}
	if !p.Draft.Authors.Consensus(p.Hash, p.Votes) {
		return nil
	}
	// new consensus
	p.Published = true
	delete(state.Proposals, p.Hash)
	p.Journal.Published = append(p.Journal.Published, p)
	return nil
}
