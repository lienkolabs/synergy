package state

import (
	"errors"

	"github.com/lienkolabs/swell/crypto"
	"github.com/lienkolabs/swell/util"
	"github.com/lienkolabs/synergy/social/actions"
)

const (
	ReactionsCount   = 5
	ProposalDeadline = 30 * 24 * 60 * 60
)

type State struct {
	Epoch        uint64
	Members      map[crypto.Hash]string
	PendingMedia map[crypto.Hash]*PendingMedia // multi-part media file
	Media        map[crypto.Hash][]byte
	Drafts       map[crypto.Hash]*Draft
	Edits        map[crypto.Hash]*Edit
	Releases     map[crypto.Hash]*Release
	Events       map[crypto.Hash]*Event
	Collectives  map[crypto.Hash]*Collective
	Boards       map[crypto.Hash]*Board
	Proposals    map[crypto.Hash]Proposal // proposals pending vote actions
	Deadline     map[uint64][]crypto.Hash
	Reactions    [ReactionsCount]map[crypto.Hash]uint

	action Notifier
}

func (s *State) Collective(name string) (*Collective, bool) {
	hash := crypto.Hasher([]byte(name))
	col, ok := s.Collectives[hash]
	return col, ok
}

func (s *State) Board(name string) (*Board, bool) {
	hash := crypto.Hasher([]byte(name))
	col, ok := s.Boards[hash]
	return col, ok
}

func (s *State) IsMember(token crypto.Token) bool {
	hash := crypto.HashToken(token)
	_, ok := s.Members[hash]
	return ok
}

func (s *State) Notify(origin Action, objHash crypto.Hash) {
	s.action.Notify(origin, s.hashToObjectType(objHash), objHash)
}

func (s *State) NextBlock() {
	if deadline, ok := s.Deadline[s.Epoch]; ok {
		for _, hash := range deadline {
			delete(s.Proposals, hash)
			s.Notify(ExpireProposal, hash)
		}
	}
}

func (s *State) hashToObjectType(hash crypto.Hash) Object {
	if _, ok := s.Members[hash]; ok {
		return MemberObject
	}
	if _, ok := s.Drafts[hash]; ok {
		return DraftObject
	}
	if _, ok := s.Edits[hash]; ok {
		return EditObject
	}
	if _, ok := s.Media[hash]; ok {
		return MediaObject
	}
	return NoObject
}

func (s *State) setDeadline(epoch uint64, hash crypto.Hash) {
	if epoch <= s.Epoch {
		return
	}
	if deadlines, ok := s.Deadline[epoch]; ok {
		s.Deadline[epoch] = append(deadlines, hash)
	} else {
		s.Deadline[epoch] = []crypto.Hash{hash}
	}
}

func (s *State) AcceptCheckinEvent(accept *actions.AcceptCheckinEvent) error {
	// need to refactor action
	return nil
}

func (s *State) ImprintStamp(stamp *actions.ImprintStamp) error {
	if !s.IsMember(stamp.Author) {
		return errors.New("not a member")
	}
	release, ok := s.Releases[stamp.Hash]
	if !ok {
		return errors.New("release not found")
	}
	collective, ok := s.Collective(stamp.OnBehalfOf)
	if !ok {
		return errors.New("collective not found")
	}
	hash := crypto.Hasher(stamp.Serialize())
	vote := actions.Vote{
		Epoch:   stamp.Epoch,
		Author:  stamp.Author,
		Reasons: "commit",
		Hash:    hash,
		Approve: true,
	}
	newStamp := Stamp{
		Reputation: collective,
		Release:    release,
		Hash:       stamp.Hash,
		Votes:      []actions.Vote{vote},
	}
	if collective.Consensus(hash, newStamp.Votes) {
		newStamp.Imprinted = true
		release.Stamps = append(release.Stamps, &newStamp)
		return nil
	}
	s.Proposals[hash] = &newStamp
	return nil
}

func (s *State) ChekinEvent(checkin *actions.CheckinEvent) error {
	if !s.IsMember(checkin.Author) {
		return errors.New("not an author")
	}
	event, ok := s.Events[checkin.EventHash]
	if !ok {
		return errors.New("event not found")
	}
	if _, ok := event.Checkin[checkin.Author]; ok {
		return errors.New("already checkin")
	}
	event.Checkin[checkin.Author] = nil
	return nil
}

func (s *State) UpdateEvent(create *actions.UpdateEvent) error {
	event, ok := s.Events[create.EventHash]
	if !ok {
		return errors.New("event not found")
	}
	if !event.Managers.IsMember(create.Author) {
		return errors.New("not a manager of the event")
	}
	if create.Description != "" {
		event.Description = create.Description
	}
	if create.Venue != "" {
		event.Venue = create.Venue
	}
	// TODO what is about the other stuff???
	return nil
}

func (s *State) CancelEvent(cancel *actions.CancelEvent) error {
	event, ok := s.Events[cancel.Hash]
	if !ok {
		return errors.New("event not found")
	}
	if !event.Managers.IsMember(cancel.Author) {
		return errors.New("not a manager")
	}
	event.Live = false
	return nil
}

func (s *State) CreateEvent(create *actions.CreateEvent) error {
	if !s.IsMember(create.Author) {
		return errors.New("not a member")
	}
	collective, ok := s.Collective(create.OnBehalfOf)
	if !ok {
		return errors.New("collective not found")
	}
	hash := crypto.Hasher(create.Serialize())
	vote := actions.Vote{
		Epoch:   create.Epoch,
		Author:  create.Author,
		Reasons: "commit",
		Hash:    hash,
		Approve: true,
	}
	event := Event{
		Collective:   collective,
		StartAt:      create.StartAt,
		EstimatedEnd: create.EstimatedEnd,
		Description:  create.Description,
		Venue:        create.Venue,
		Open:         create.Open,
		Public:       create.Public,
		Hash:         hash,
		Votes:        []actions.Vote{vote},
		Checkin:      make(map[crypto.Token]*actions.AcceptCheckinEvent),
		Live:         false,
	}
	if len(create.Managers) > 0 {
		managers := make(map[crypto.Token]struct{})
		for _, manager := range create.Managers {
			managers[manager] = struct{}{}
		}
		event.Managers = &UnamedCollective{
			Members:  managers,
			Majority: 0,
		}
	}
	if collective.Consensus(hash, []actions.Vote{vote}) {
		if _, ok := s.Events[hash]; ok {
			return errors.New("event already booked")
		}
		s.Events[hash] = &event
		return nil
	}
	if _, ok := s.Proposals[hash]; ok {
		return errors.New("event already booked")
	}
	s.Proposals[hash] = &event
	return nil
}

func (s *State) MultipartMedia(media *actions.MultipartMedia) error {
	pending, ok := s.PendingMedia[media.Hash]
	if !ok {
		return errors.New("referred media not found")
	}
	total, err := pending.Append(media)
	if err != nil {
		return err
	}
	if total != nil {
		delete(s.PendingMedia, media.Hash)
		s.Media[media.Hash] = total
		s.Notify(MediaUpload, media.Hash)
	}
	return nil
}

func (s *State) Release(release *actions.ReleaseDraft) error {
	draft, ok := s.Drafts[release.ContentHash]
	if !ok {
		return errors.New("draft not found")
	}
	if !draft.Authors.IsMember(release.Author) {
		return errors.New("not an author")
	}
	hash := crypto.Hasher(release.Serialize())
	vote := actions.Vote{
		Epoch:   release.Epoch,
		Author:  release.Author,
		Reasons: "commit",
		Hash:    hash,
		Approve: true,
	}
	newRelease := Release{
		Epoch:  release.Epoch,
		Draft:  draft,
		Hash:   release.ContentHash,
		Votes:  []actions.Vote{vote},
		Stamps: make([]*Stamp, 0),
	}
	if draft.Authors.Consensus(hash, []actions.Vote{vote}) {
		if _, ok := s.Releases[release.ContentHash]; ok {
			return errors.New("already released")
		}
		newRelease.Released = true
		s.Releases[release.ContentHash] = &newRelease
		return nil
	}
	s.Proposals[hash] = &newRelease
	return nil
}

func (s *State) UpdateBoard(update *actions.UpdateBoard) error {
	if !s.IsMember(update.Author) {
		return errors.New("not a member")
	}
	board, ok := s.Board(update.Board)
	if !ok {
		return errors.New("board not found")
	}
	hash := crypto.Hasher([]byte(update.Serialize()))
	vote := actions.Vote{
		Epoch:   update.Epoch,
		Author:  update.Author,
		Reasons: "commit",
		Hash:    hash,
		Approve: true,
	}
	if board.Collective.Consensus(hash, []actions.Vote{vote}) {
		board.Keyword = update.Keywords
		board.Description = update.Description
		board.Editors.ChangeMajority(int(update.PinMajority))
		return nil
	}
	pending := PendingUpdateBoard{
		Keywords:    update.Keywords,
		Description: update.Description,
		PinMajority: int(update.PinMajority),
		Board:       board,
		Hash:        hash,
		Votes:       []actions.Vote{vote},
	}
	s.Proposals[hash] = &pending
	return nil
	// TODO notify
}

func (s *State) CreateBoard(board *actions.CreateBoard) error {
	if !s.IsMember(board.Author) {
		return errors.New("not a member")
	}
	if _, ok := s.Board(board.Name); ok {
		return errors.New("board already exists")
	}
	collective, ok := s.Collective(board.OnBehalfOf)
	if !ok {
		return errors.New("collective unkown")
	}
	newBoard := Board{
		Name:        board.Name,
		Keyword:     board.Keywords,
		Description: board.Description,
		Collective:  collective,
		Editors: &UnamedCollective{
			Members:  make(map[crypto.Token]struct{}),
			Majority: int(board.PinMajority),
		},
		Pinned: make([]*Draft, 0),
	}
	hash := crypto.Hasher([]byte(board.Serialize()))
	vote := actions.Vote{
		Epoch:   board.Epoch,
		Author:  board.Author,
		Reasons: "commit",
		Hash:    hash,
		Approve: true,
	}
	if collective.Consensus(hash, []actions.Vote{vote}) {
		s.Boards[hash] = &newBoard
		return nil
	}
	s.Proposals[hash] = &PendingBoard{
		Board: &newBoard,
		Hash:  crypto.Hasher(board.Serialize()),
		Votes: []actions.Vote{vote},
	}
	// TODO: notify
	return nil
}

func (s *State) SignIn(signin *actions.Signin) error {
	hash := crypto.HashToken(signin.Author)
	if _, ok := s.Members[hash]; ok {
		return errors.New("already a member of synergy")
	}
	s.Members[hash] = signin.Handle
	s.Notify(SigninAction, hash)
	return nil
}

func (s *State) CreateCollective(create *actions.CreateCollective) error {
	if s.IsMember(create.Author) {
		return errors.New("not a member of synergy")
	}
	if _, ok := s.Collective(create.Name); ok {
		return errors.New("collective already exists")
	}
	if create.Policy.Majority < 0 || create.Policy.Majority > 100 || create.Policy.SuperMajority < 0 || create.Policy.SuperMajority > 100 {
		return errors.New("invalid policy")
	}
	hash := crypto.Hasher([]byte(create.Name))
	s.Collectives[hash] = &Collective{
		Name:        create.Name,
		Members:     map[crypto.Token]struct{}{},
		Description: create.Description,
		Policy: actions.Policy{
			Majority:      create.Policy.Majority,
			SuperMajority: create.Policy.SuperMajority,
		},
	}
	return nil
}

func (s *State) UpdateCollective(update *actions.UpdateCollective) error {
	collective, ok := s.Collective(update.OnBehalfOf)
	if !ok {
		return errors.New("unkown collective")
	}
	if !collective.IsMember(update.Author) {
		return errors.New("not a member of collective")
	}
	hash := crypto.Hasher(update.Serialize()) // proposal hash = hash of instruction
	vote := actions.Vote{
		Epoch:   update.Epoch,
		Author:  update.Author,
		Reasons: "commit",
		Hash:    hash,
		Approve: true,
	}

	if update.Policy != nil {
		if update.Policy.Majority < 0 || update.Policy.Majority > 100 || update.Policy.SuperMajority < 0 || update.Policy.SuperMajority > 100 {
			return errors.New("invalid policy")
		}
		if collective.SuperConsensus(hash, []actions.Vote{vote}) {
			if update.Description != "" {
				collective.Description = update.Description
			}
			collective.Policy = actions.Policy{
				Majority:      update.Policy.Majority,
				SuperMajority: update.Policy.SuperMajority,
			}
			return nil
		}
	} else {
		if collective.Consensus(hash, []actions.Vote{vote}) {
			if update.Description != "" {
				collective.Description = update.Description
			}
			return nil
		}
	}

	pending := PendingUpdate{
		Update: update,
		// consensus is based on the collective composition at the moment
		// of incorporation of instruction
		Collective: collective.Photo(),
		Hash:       hash,
		Votes:      []actions.Vote{vote},
	}
	if update.Policy != nil {
		pending.ChangePolicy = true
	}
	s.Proposals[hash] = &pending
	s.setDeadline(update.Epoch+ProposalDeadline, hash)
	return nil

}

func (s *State) RequestMembership(request *actions.RequestMembership) error {
	if !s.IsMember(request.Author) {
		return errors.New("not a member of synergy")
	}
	collective, ok := s.Collective(request.Collective)
	if !ok {
		return errors.New("collective not found")
	}
	if collective.IsMember(request.Author) {
		return errors.New("already a member")
	}
	hash := crypto.Hasher(request.Serialize())
	pending := PendingRequestMembership{
		Request:    request,
		Collective: collective.Photo(),
		Hash:       hash,
		Votes:      make([]actions.Vote, 0),
	}
	s.Proposals[hash] = &pending
	s.setDeadline(request.Epoch+ProposalDeadline, hash)
	return nil
}

func (s *State) RemoveMember(remove *actions.RemoveMember) error {
	collective, ok := s.Collective(remove.OnBehalfOf)
	if !ok {
		return errors.New("collective not found")
	}
	if !collective.IsMember(remove.Author) {
		return errors.New("author not a member of collective")
	}
	if !collective.IsMember(remove.Member) {
		return errors.New("member to be removed not a member of collective")
	}
	if remove.Author.Equal(remove.Member) {
		delete(collective.Members, remove.Author)
		return nil
	}
	hash := crypto.Hasher(remove.Serialize())
	vote := actions.Vote{
		Epoch:   remove.Epoch,
		Author:  remove.Author,
		Reasons: "commit",
		Hash:    hash,
		Approve: true,
	}
	if collective.Consensus(hash, []actions.Vote{vote}) {
		delete(collective.Members, remove.Author)
		return nil
	}
	pending := PendingRemoveMember{
		Remove:     remove,
		Collective: collective.Photo(),
		Hash:       hash,
		Votes:      []actions.Vote{vote},
	}
	s.Proposals[hash] = &pending
	s.setDeadline(remove.Epoch+ProposalDeadline, hash)
	return nil
}

func (s *State) Reaction(reaction *actions.React) error {
	if reaction.Reaction >= ReactionsCount {
		return errors.New("invalid reaction")
	}
	// TODO: should check if hash is known?
	if count, ok := s.Reactions[reaction.Reaction][reaction.Hash]; ok {
		s.Reactions[reaction.Reaction][reaction.Hash] = count + 1
	} else {
		s.Reactions[reaction.Reaction][reaction.Hash] = 1
	}
	if obj := s.hashToObjectType(reaction.Hash); obj != NoObject {
		s.action.Notify(ReactAction, obj, reaction.Hash)
	}
	return nil
}

func (s *State) Edit(edit *actions.Edit) error {
	if _, ok := s.Media[edit.ContentHash]; ok {
		return errors.New("hash already claimed")
	}
	if _, ok := s.PendingMedia[edit.ContentHash]; ok {
		return errors.New("hash already claimed")
	}
	if edit.NumberOfParts > 1 {
		pending := PendingMedia{
			Hash:          edit.ContentHash,
			NumberOfParts: edit.NumberOfParts,
			Parts:         make([]*actions.MultipartMedia, int(edit.NumberOfParts)),
		}
		s.PendingMedia[edit.ContentHash] = &pending
		pending.Parts[0] = &actions.MultipartMedia{
			Hash: edit.ContentHash,
			Part: 0,
			Of:   edit.NumberOfParts,
			Data: edit.Content,
		}
	} else {
		s.Media[edit.ContentHash] = edit.Content
	}
	newEdit := Edit{
		Reasons:  edit.Reasons,
		Draft:    edit.EditedDraft,
		EditType: edit.ContentType,
		Votes: []actions.Vote{
			{
				Epoch:   edit.Epoch,
				Author:  edit.Author,
				Reasons: "submission",
				Hash:    edit.ContentHash,
				Approve: true,
			}},
	}

	if edit.OnBehalfOf != "" {
		collective, ok := s.Collective(edit.OnBehalfOf)
		if !ok {
			return errors.New("collective unkown")
		}
		if !collective.IsMember(edit.Author) {
			return errors.New("not a member of collective")
		}
		newEdit.Authors = collective
		if collective.Consensus(edit.ContentHash, newEdit.Votes) {
			s.Edits[edit.ContentHash] = &newEdit
		} else {
			s.Proposals[edit.ContentHash] = &newEdit
		}
	} else if edit.CoAuthors != nil && len(edit.CoAuthors) > 0 {
		newEdit.Authors = Authors(1+len(edit.CoAuthors), append(edit.CoAuthors, edit.Author)...)
		s.Proposals[edit.ContentHash] = &newEdit
	} else {
		newEdit.Authors = Authors(1, edit.Author)
		s.Edits[edit.ContentHash] = &newEdit
	}
	s.action.Notify(EditAction, DraftObject, edit.EditedDraft)
	s.action.Notify(EditAction, AuthorObject, crypto.HashToken(edit.Author))
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
func (s *State) Draft(draft *actions.Draft) error {
	if _, ok := s.Media[draft.ContentHash]; !ok {
		return errors.New("media file not available")
	}
	var previous *Draft
	if draft.PreviousDraft != crypto.ZeroHash {
		if previous, ok := s.Drafts[draft.PreviousDraft]; !ok {
			return errors.New("invalid previous version")
		} else {
			isPreviousAuthor := previous.Authors.IsMember(draft.Author)
			if !isPreviousAuthor {
				return errors.New("unauthorized version")
			}
		}
	}
	newDraft := &Draft{
		Title:       draft.Title,
		Description: draft.Description,
		//Authors:            draft.CoAuthors,
		DraftType:       draft.ContentType,
		DraftHash:       draft.ContentHash,
		PreviousVersion: previous,
		References:      draft.References,
		Votes:           make([]actions.Vote, 0),
	}
	if draft.CoAuthors == nil {
		if draft.OnBehalfOf == "" {
			// create single author collective
			newDraft.Authors = Authors(1, draft.Author)
		}
		behalf, ok := s.Collective(draft.OnBehalfOf)
		if !ok {
			return errors.New("named collective not recognizedx")
		}
		newDraft.Authors = behalf
	}
	selfVote := actions.Vote{
		Epoch:   draft.Epoch,
		Author:  draft.Author,
		Reasons: "submission",
		Hash:    draft.ContentHash,
		Approve: true,
	}
	s.Proposals[newDraft.DraftHash] = newDraft
	newDraft.IncorporateVote(selfVote, s)
	if newDraft.PreviousVersion != nil {
		s.action.Notify(DraftAction, DraftObject, draft.PreviousDraft)
	}
	return nil
}

func (s *State) Vote(vote actions.Vote) error {
	if proposed, ok := s.Proposals[vote.Hash]; ok {
		return proposed.IncorporateVote(vote, s)
	}
	if draft, ok := s.Drafts[vote.Hash]; ok {
		return draft.IncorporateVote(vote, s)
	}
	return errors.New("not open vote found")
}

func (s *State) Pin(pin actions.Pin) error {
	board, ok := s.Board(pin.Board)
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
	action := Pin{
		Hash:  hash,
		Epoch: pin.Epoch,
		Board: board,
		Draft: draft,
		Pin:   pin.Pin,
		Votes: make([]actions.Vote, 0),
	}
	selfVote := actions.Vote{
		Epoch:   pin.Epoch,
		Author:  pin.Author,
		Reasons: "submission",
		Hash:    hash,
		Approve: true,
	}
	s.Proposals[hash] = &action
	return action.IncorporateVote(selfVote, s)
}

func (s *State) BoardEditor(action actions.BoardEditor) error {
	board, ok := s.Board(action.Board)
	if !ok {
		return errors.New("invalid board")
	}
	if s.IsMember(action.Editor); !ok {
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
	proposal := BoardEditor{
		Hash:   hash,
		Epoch:  action.Epoch,
		Board:  board,
		Editor: action.Editor,
		Insert: action.Insert,
		Votes:  make([]actions.Vote, 0),
	}
	selfVote := actions.Vote{
		Epoch:   action.Epoch,
		Author:  action.Author,
		Reasons: "submission",
		Hash:    hash,
		Approve: true,
	}
	s.Proposals[hash] = &proposal
	return proposal.IncorporateVote(selfVote, s)
}
