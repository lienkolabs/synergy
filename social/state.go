package social

import (
	"errors"

	"github.com/lienkolabs/swell/crypto"
	"github.com/lienkolabs/swell/util"
)

const ReactionsCount = 5

type State struct {
	Members          map[crypto.Token]struct{}
	PendingMedia     map[crypto.Hash]*PendingMedia // multi-part media file
	Media            map[crypto.Hash][]byte
	Drafts           map[crypto.Hash]*Draft
	Edits            map[crypto.Hash]*Edit
	NamedCollectives map[string]*Collective
	Boards           map[string]*Board
	Proposals        map[crypto.Hash]Proposal // proposals pending vote actions
	Reactions        [ReactionsCount]map[crypto.Hash]uint
	Events           chan crypto.Hash
}

func (s *State) IncorporateReaction(reaction *ReactionInstruction) error {
	if reaction.Reaction >= ReactionsCount {
		return errors.New("invalid reaction")
	}
	// TODO: should check if hash is known?
	if count, ok := s.Reactions[reaction.Reaction][reaction.Hash]; ok {
		s.Reactions[reaction.Reaction][reaction.Hash] = count + 1
	} else {
		s.Reactions[reaction.Reaction][reaction.Hash] = 1
	}
	return nil
}

func (s *State) IncorporateEditInstruction(edit *EditInstruction) error {
	if _, ok := s.Media[edit.EditHash]; ok {
		return errors.New("hash already claimed")
	}
	if _, ok := s.PendingMedia[edit.EditHash]; ok {
		return errors.New("hash already claimed")
	}
	if edit.Parts > 1 {
		pending := PendingMedia{
			Hash:          edit.EditHash,
			NumberOfParts: edit.Parts,
			Parts:         make([]*MultipartMedia, edit.Parts),
		}
		s.PendingMedia[edit.EditHash] = &pending
		pending.Parts[0] = &MultipartMedia{
			Hash: edit.EditHash,
			Part: 0,
			Of:   edit.Parts,
			Data: edit.Data,
		}
	} else {
		s.Media[edit.EditHash] = edit.Data
	}
	newEdit := Edit{
		Reasons:  edit.Reasons,
		Draft:    edit.EditedDraft,
		EditType: edit.EditType,
		Signatures: []VoteInstruction{
			{
				Epoch:     edit.Epoch,
				Author:    edit.Author,
				Reasons:   "submission",
				Hash:      edit.EditHash,
				Approve:   true,
				Signature: edit.EditSignature,
			}},
	}

	if edit.OnBehalfOf != "" {
		collective, ok := s.NamedCollectives[edit.OnBehalfOf]
		if !ok {
			return errors.New("collective unkown")
		}
		if !collective.IsMember(edit.Author) {
			return errors.New("not a member of collective")
		}
		newEdit.Authors = collective
		if collective.Consensus(edit.EditHash, newEdit.Signatures) {
			s.Edits[edit.EditHash] = &newEdit
		} else {
			s.Proposals[edit.EditHash] = &newEdit
		}
	} else if edit.CoAuthors != nil && len(edit.CoAuthors) > 0 {
		newEdit.Authors = Authors(1+len(edit.CoAuthors), append(edit.CoAuthors, edit.Author)...)
		s.Proposals[edit.EditHash] = &newEdit
	} else {
		newEdit.Authors = Authors(1, edit.Author)
		s.Edits[edit.EditHash] = &newEdit
	}
	return nil
}

// IncorporateDraftInstruction checks if proposed draft is valid and if so
// incorporate it as ProposedDraft if further consent is necessary of as
// Draft if the instruction author has alone authority to submit the draft.
//
// Checks:
// a) It must refer to a known media file hash (pending media will not be
//
//	accepted)
//
// b) If it has a designated previous version, the instruction auhtor must be
//
//	an accredited author of the previous version or a member of the collective
//
// c) If draft is submitted on behalf of a named collective, this name must
//
//	be recognized by the state
//
// d)
func (s *State) IncorporateDraftInstruction(draft *DraftInstruction) error {
	if _, ok := s.Media[draft.DraftHash]; !ok {
		return errors.New("media file not available")
	}
	var previous *Draft
	if draft.PreviousVersion != crypto.ZeroHash {
		if previous, ok := s.Drafts[draft.PreviousVersion]; !ok {
			return errors.New("invalid previous version")
		} else {
			isPreviousAuthor := previous.Authors.IsMember(draft.InstructionAuthor)
			if !isPreviousAuthor {
				return errors.New("unauthorized version")
			}
		}
	}
	newDraft := &Draft{
		Title:              draft.DraftTitle,
		Abstract:           draft.DraftAbstract,
		Authors:            draft.DraftAuthors,
		DraftType:          draft.DraftType,
		DraftHash:          draft.DraftHash,
		PreviousVersion:    previous,
		InternalReferences: draft.InternalReferences,
		Votes:              make([]VoteInstruction, 0),
	}
	if draft.DraftAuthors == nil {
		if draft.OnBehalfOf == "" {
			// create single author collective
			newDraft.Authors = Authors(1, draft.InstructionAuthor)
		}
		behalf, ok := s.NamedCollectives[draft.OnBehalfOf]
		if !ok {
			return errors.New("named collective not recognizedx")
		}
		newDraft.Authors = behalf
	}
	selfVote := VoteInstruction{
		Epoch:     draft.Epoch,
		Author:    draft.InstructionAuthor,
		Reasons:   "submission",
		Hash:      draft.DraftHash,
		Approve:   true,
		Signature: draft.HashSignature,
	}
	s.Proposals[newDraft.DraftHash] = newDraft
	newDraft.IncorporateVote(selfVote, s)
	return nil
}

func (s *State) IncorporateGenericVoteInstruction(vote VoteInstruction) error {
	if proposed, ok := s.Proposals[vote.Hash]; ok {
		return proposed.IncorporateVote(vote, s)
	}
	if draft, ok := s.Drafts[vote.Hash]; ok {
		return draft.IncorporateVote(vote, s)
	}
	return errors.New("invalid hash")
}

func (s *State) IncorporatePinInstruction(pin PinInstruction) error {
	board, ok := s.Boards[pin.Board]
	if !ok {
		return errors.New("invalid board")
	}
	draft, ok := s.Drafts[pin.Draft]
	if !ok {
		return errors.New("invalid draft")
	}
	bytes := make([]byte, 0)
	util.PutUint64(pin.Epoch, &bytes)
	util.PutHash(pin.Draft, &bytes)
	util.PutString(pin.Board, &bytes)
	if pin.Pin {
		util.PutByte(1, &bytes)
	} else {
		util.PutByte(0, &bytes)
	}
	hash := crypto.Hasher(bytes)
	action := BoardPinAction{
		Hash:  hash,
		Epoch: pin.Epoch,
		Board: board,
		Draft: draft,
		Pin:   pin.Pin,
		Votes: make([]VoteInstruction, 0),
	}
	selfVote := VoteInstruction{
		Epoch:     pin.Epoch,
		Author:    pin.Author,
		Reasons:   "submission",
		Hash:      hash,
		Approve:   true,
		Signature: pin.HashSignature,
	}
	s.Proposals[hash] = &action
	return action.IncorporateVote(selfVote, s)
}

func (s *State) IncorporateBoardEditorInstruction(action BoardEditorInstruction) error {
	board, ok := s.Boards[action.Board]
	if !ok {
		return errors.New("invalid board")
	}
	if _, ok := s.Members[action.Editor]; !ok {
		return errors.New("invalid editor")
	}

	bytes := make([]byte, 0)
	util.PutUint64(action.Epoch, &bytes)
	util.PutToken(action.Editor, &bytes)
	util.PutString(action.Board, &bytes)
	if action.Insert {
		util.PutByte(1, &bytes)
	} else {
		util.PutByte(0, &bytes)
	}
	hash := crypto.Hasher(bytes)
	proposal := BoardEditorAction{
		Hash:   hash,
		Epoch:  action.Epoch,
		Board:  board,
		Editor: action.Editor,
		Insert: action.Insert,
		Votes:  make([]VoteInstruction, 0),
	}
	selfVote := VoteInstruction{
		Epoch:     action.Epoch,
		Author:    action.Author,
		Reasons:   "submission",
		Hash:      hash,
		Approve:   true,
		Signature: action.HashSignature,
	}
	s.Proposals[hash] = &proposal
	return proposal.IncorporateVote(selfVote, s)
}
