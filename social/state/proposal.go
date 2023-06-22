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
	UnkownProposal
)

var proposalNames = []string{
	"Update Collective",
	"Request Membership",
	"Remove Member",
	"Draft",
	"Edit",
	"Create Board",
	"Update Board",
	"Pin",
	"Board Editor",
	"Release Draft",
	"Imprint Stamp",
	"React",
	"Create Event",
	"Cancel Event",
	"Update Event",
	"Unkown",
}

var ErrProposalNotFound = errors.New("proposal not found")

type Proposal interface {
	IncorporateVote(vote actions.Vote, state *State) error
}

func NewProposals() *Proposals {
	return &Proposals{
		all:               make(map[crypto.Hash]byte),
		updateCollective:  make(map[crypto.Hash]*PendingUpdate),
		RequestMembership: make(map[crypto.Hash]*PendingRequestMembership),
		RemoveMember:      make(map[crypto.Hash]*PendingRemoveMember),
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
	RequestMembership map[crypto.Hash]*PendingRequestMembership
	RemoveMember      map[crypto.Hash]*PendingRemoveMember
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

func (p *Proposals) GetEvent(hash crypto.Hash) *Event {
	e, _ := p.createEvent[hash]
	return e
}

func (p *Proposals) Delete(hash crypto.Hash) {
	delete(p.all, hash)
	delete(p.updateCollective, hash)
	delete(p.RequestMembership, hash)
	delete(p.RemoveMember, hash)
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

func (p *Proposals) Kind(hash crypto.Hash) byte {
	kind, ok := p.all[hash]
	if !ok {
		return UnkownProposal
	}
	return kind
}

func (p *Proposals) KindText(hash crypto.Hash) string {
	return proposalNames[p.Kind(hash)]
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
	p.indexHash(update.Collective, update.Hash)
	p.all[update.Hash] = UpdateCollectiveProposal
	p.updateCollective[update.Hash] = update
}

func (p *Proposals) AddRequestMembership(update *PendingRequestMembership) {
	p.indexHash(update.Collective, update.Hash)
	p.all[update.Hash] = RequestMembershipProposal
	p.RequestMembership[update.Hash] = update
}

func (p *Proposals) AddPendingRemoveMember(update *PendingRemoveMember) {
	p.indexHash(update.Collective, update.Hash)
	p.all[update.Hash] = RemoveMemberProposal
	p.RemoveMember[update.Hash] = update
}

func (p *Proposals) AddDraft(update *Draft) {
	p.indexHash(update.Authors, update.DraftHash)
	if update.PreviousVersion != nil {
		p.indexHash(update.PreviousVersion.Authors, update.DraftHash)
	}
	p.all[update.DraftHash] = DraftProposal
	p.draf[update.DraftHash] = update
}

func (p *Proposals) AddEdit(update *Edit) {
	p.indexHash(update.Draft.Authors, update.Edit)
	p.indexHash(update.Authors, update.Edit)
	p.all[update.Edit] = EditProposal
	p.edit[update.Edit] = update
}

func (p *Proposals) AddPendingBoard(update *PendingBoard) {
	p.indexHash(update.Board.Collective, update.Hash)
	p.all[update.Hash] = CreateBoardProposal
	p.createBoard[update.Hash] = update
}

func (p *Proposals) AddPendingUpdateBoard(update *PendingUpdateBoard) {
	p.indexHash(update.Board.Editors, update.Hash)
	p.all[update.Hash] = UpdateBoardProposal
	p.updateBoard[update.Hash] = update
}

func (p *Proposals) AddPin(update *Pin) {
	p.indexHash(update.Board.Editors, update.Hash)
	p.all[update.Hash] = PinProposal
	p.pin[update.Hash] = update
}

func (p *Proposals) AddBoardEditor(update *BoardEditor) {
	p.indexHash(update.Board.Collective, update.Hash)
	p.all[update.Hash] = BoardEditorProposal
	p.boardEditor[update.Hash] = update
}

func (p *Proposals) AddRelease(update *Release) {
	p.indexHash(update.Draft.Authors, update.Hash)
	p.all[update.Hash] = ReleaseDraftProposal
	p.releaseDraft[update.Hash] = update
}

func (p *Proposals) AddStamp(update *Stamp) {
	p.indexHash(update.Reputation, update.Hash) // reputation aqui é = um membro ou coletivo que vai dar o stamp ??
	p.all[update.Hash] = ImprintStampProposal
	p.imprintStamp[update.Hash] = update
}

func (p *Proposals) AddEvent(update *Event) {
	p.indexHash(update.Collective, update.Hash)
	p.all[update.Hash] = CreateEventProposal
	p.createEvent[update.Hash] = update
}

func (p *Proposals) AddCancelEvent(update *CancelEvent) {
	p.indexHash(update.Event.Collective, update.Hash)
	p.all[update.Hash] = CancelEventProposal
	p.cancelEvent[update.Hash] = update
}

func (p *Proposals) AddEventUpdate(update *EventUpdate) {
	p.indexHash(update.Event.Managers, update.Hash)
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
		proposal = p.RemoveMember[hash]
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

func (p *Proposals) OnBehalfOf(hash crypto.Hash) string {
	kind, ok := p.all[hash]
	if !ok {
		return ""
	}
	switch kind {
	case UpdateCollectiveProposal:
		proposal := p.updateCollective[hash]
		return proposal.Collective.Name
	case RemoveMemberProposal:
		proposal := p.RemoveMember[hash]
		return proposal.Collective.Name
	case DraftProposal:
		proposal := p.draf[hash]
		return proposal.Authors.CollectiveName()
	case EditProposal:
		proposal := p.edit[hash]
		return proposal.Authors.CollectiveName()
	case CreateBoardProposal:
		proposal := p.createBoard[hash]
		return proposal.Board.Collective.Name
	case UpdateBoardProposal:
		proposal := p.updateBoard[hash]
		return proposal.Board.Collective.Name
	case PinProposal:
		proposal := p.pin[hash]
		return proposal.Board.Name
	case BoardEditorProposal:
		proposal := p.boardEditor[hash]
		return proposal.Board.Collective.Name
	case ReleaseDraftProposal:
		proposal := p.releaseDraft[hash]
		return proposal.Draft.Authors.CollectiveName()
	case ImprintStampProposal:
		proposal := p.imprintStamp[hash]
		return proposal.Reputation.Name
	case ReactProposal:
		//
	case CreateEventProposal:
		proposal := p.createEvent[hash]
		return proposal.Collective.Name
	case CancelEventProposal:
		proposal := p.cancelEvent[hash]
		return proposal.Event.Collective.Name
	case UpdateEventProposal:
		proposal := p.updateEvent[hash]
		return proposal.Event.Collective.Name
	}
	return ""
}
