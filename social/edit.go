package social

import "github.com/lienkolabs/swell/crypto"

type EditInstruction struct {
	Epoch         uint64
	Author        crypto.Token
	CoAuthors     []crypto.Token
	OnBehalfOf    string
	Reasons       string
	EditedDraft   crypto.Hash
	EditType      string
	EditHash      crypto.Hash
	EditSignature crypto.Signature
	Parts         byte
	Data          []byte
}

type Edit struct {
	Authors    Consensual
	Reasons    string
	Draft      crypto.Hash
	EditType   string
	Edit       crypto.Hash
	Signatures []VoteInstruction
}

func (e *Edit) IncorporateVote(vote VoteInstruction, state *State) error {
	if err := isValidVote(e.Edit, vote, e.Signatures); err != nil {
		return err
	}
	e.Signatures = append(e.Signatures, vote)
	if !e.Authors.Consensus(e.Edit, e.Signatures) {
		return nil
	}
	// new consensus
	delete(state.Proposals, e.Edit)
	// to do where to put edits?
	return nil
}
