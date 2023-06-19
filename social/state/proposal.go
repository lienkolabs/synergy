package state

import (
	"errors"

	"github.com/lienkolabs/swell/crypto"
	"github.com/lienkolabs/synergy/social/actions"
)

const (
	UpdateCollectiveProposal byte = iota
	RequestMembershipProposal
	RemoveMemberProposal
	DraftProposal
	EditProposal
	CreateBoardProposal
	UpdateBoardProposal
	PinProposal
	BoardEditorProposal
	ReleaseDraftProposal
	ImprintStampProposal
	ReactProposal
	CreateEventProposal
	CancelEventProposal
	UpdateEventProposal
)

var ErrProposalNotFound = errors.New("proposal not found")

type Proposal interface {
	IncorporateVote(vote actions.Vote, state *State) error
}

func NewProposals() *Proposals {
	return &Proposals{
		all:               make(map[crypto.Hash]byte),
		updateCollective:  make(map[crypto.Hash]*PendingUpdate),
		requestMembership: make(map[crypto.Hash]*PendingRequestMembership),
		removeMember:      make(map[crypto.Hash]*PendingRemoveMember),
		draf:              make(map[crypto.Hash]*Draft),
		edit:              make(map[crypto.Hash]*Edit),
		createBoard:       make(map[crypto.Hash]*PendingBoard),
		updateBoard:       make(map[crypto.Hash]*PendingUpdateBoard),
		pin:               make(map[crypto.Hash]*Pin),
		boardEditor:       make(map[crypto.Hash]*BoardEditor),
		releaseDraft:      make(map[crypto.Hash]*Release),
		imprintStamp:      make(map[crypto.Hash]*Stamp),
		//react map[crypto.Hash]*
		createEvent: make(map[crypto.Hash]*Event),
		cancelEvent: make(map[crypto.Hash]*CancelEvent),
		updateEvent: make(map[crypto.Hash]*EventUpdate),
	}
}

type Proposals struct {
	all               map[crypto.Hash]byte
	index             map[crypto.Token]map[crypto.Hash]struct{}
	updateCollective  map[crypto.Hash]*PendingUpdate
	requestMembership map[crypto.Hash]*PendingRequestMembership
	removeMember      map[crypto.Hash]*PendingRemoveMember
	draf              map[crypto.Hash]*Draft
	edit              map[crypto.Hash]*Edit
	createBoard       map[crypto.Hash]*PendingBoard
	updateBoard       map[crypto.Hash]*PendingUpdateBoard
	pin               map[crypto.Hash]*Pin
	boardEditor       map[crypto.Hash]*BoardEditor
	releaseDraft      map[crypto.Hash]*Release
	imprintStamp      map[crypto.Hash]*Stamp
	//react map[crypto.Hash]*
	createEvent map[crypto.Hash]*Event
	cancelEvent map[crypto.Hash]*CancelEvent
	updateEvent map[crypto.Hash]*EventUpdate
}

func (p *Proposals) Delete(hash crypto.Hash) {
	delete(p.all, hash)
	delete(p.updateCollective, hash)
	delete(p.requestMembership, hash)
	delete(p.removeMember, hash)
	delete(p.draf, hash)
	delete(p.edit, hash)
	delete(p.createBoard, hash)
	delete(p.updateBoard, hash)
	delete(p.pin, hash)
	delete(p.boardEditor, hash)
	delete(p.releaseDraft, hash)
	delete(p.imprintStamp, hash)
	//react map[crypto.Hash]*
	delete(p.createEvent, hash)
	delete(p.cancelEvent, hash)
	delete(p.updateEvent, hash)
	for _, hashes := range p.index {
		delete(hashes, hash)
	}
}

func (p *Proposals) indexHash(c Consensual, hash crypto.Hash) {
	for token := range c.ListOfMembers() {
		if hashes, ok := p.index[token]; ok {
			hashes[hash] = struct{}{}
		} else {
			p.index[token] = map[crypto.Hash]struct{}{hash: {}}
		}
	}
}

func (p *Proposals) GetVotes(token crypto.Token) map[crypto.Hash]struct{} {
	hashes := p.index[token]
	return hashes
}

func (p *Proposals) AddUpdateCollective(update *PendingUpdate) {
	p.all[update.Hash] = UpdateCollectiveProposal
	p.updateCollective[update.Hash] = update
}

func (p *Proposals) AddRequestMembership(update *PendingRequestMembership) {
	p.all[update.Hash] = RequestMembershipProposal
	p.requestMembership[update.Hash] = update
}

func (p *Proposals) AddPendingRemoveMember(update *PendingRemoveMember) {
	p.all[update.Hash] = RemoveMemberProposal
	p.removeMember[update.Hash] = update
}

func (p *Proposals) AddDraft(update *Draft) {
	p.all[update.DraftHash] = DraftProposal
	p.draf[update.DraftHash] = update
}

func (p *Proposals) AddEdit(update *Edit) {
	p.all[update.Edit] = EditProposal
	p.edit[update.Edit] = update
}

func (p *Proposals) AddPendingBoard(update *PendingBoard) {
	p.all[update.Hash] = CreateBoardProposal
	p.createBoard[update.Hash] = update
}

func (p *Proposals) AddPendingUpdateBoard(update *PendingUpdateBoard) {
	p.all[update.Hash] = UpdateBoardProposal
	p.updateBoard[update.Hash] = update
}

func (p *Proposals) AddPin(update *Pin) {
	p.all[update.Hash] = PinProposal
	p.pin[update.Hash] = update
}

func (p *Proposals) AddBoardEditor(update *BoardEditor) {
	p.all[update.Hash] = BoardEditorProposal
	p.boardEditor[update.Hash] = update
}

func (p *Proposals) AddRelease(update *Release) {
	p.all[update.Hash] = ReleaseDraftProposal
	p.releaseDraft[update.Hash] = update
}

func (p *Proposals) AddStamp(update *Stamp) {
	p.all[update.Hash] = ImprintStampProposal
	p.imprintStamp[update.Hash] = update
}

func (p *Proposals) AddEvent(update *Event) {
	p.all[update.Hash] = CreateEventProposal
	p.createEvent[update.Hash] = update
}

func (p *Proposals) AddCancelEvent(update *CancelEvent) {
	p.all[update.Hash] = CancelEventProposal
	p.cancelEvent[update.Hash] = update
}

func (p *Proposals) AddEventUpdate(update *EventUpdate) {
	p.all[update.Hash] = UpdateEventProposal
	p.updateEvent[update.Hash] = update
}

func (p *Proposals) Has(hash crypto.Hash) bool {
	_, ok := p.all[hash]
	return ok
}

func (p *Proposals) IncorporateVote(vote actions.Vote, state *State) error {
	hash := vote.Hash
	var proposal Proposal
	kind, ok := p.all[hash]
	if !ok {
		return ErrProposalNotFound
	}
	switch kind {
	case UpdateCollectiveProposal:
		proposal = p.updateCollective[hash]
	case RemoveMemberProposal:
		proposal = p.removeMember[hash]
	case DraftProposal:
		proposal = p.draf[hash]
	case EditProposal:
		proposal = p.edit[hash]
	case CreateBoardProposal:
		proposal = p.createBoard[hash]
	case UpdateBoardProposal:
		proposal = p.updateBoard[hash]
	case PinProposal:
		proposal = p.pin[hash]
	case BoardEditorProposal:
		proposal = p.boardEditor[hash]
	case ReleaseDraftProposal:
		proposal = p.releaseDraft[hash]
	case ImprintStampProposal:
		proposal = p.imprintStamp[hash]
	case ReactProposal:
		//
	case CreateEventProposal:
		proposal = p.createEvent[hash]
	case CancelEventProposal:
		proposal = p.cancelEvent[hash]
	case UpdateEventProposal:
		proposal = p.updateEvent[hash]
	}
	if proposal == nil {
		return ErrProposalNotFound
	}
	return proposal.IncorporateVote(vote, state)
}
