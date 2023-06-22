package state

import (
	"errors"
	"sync"

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
		mu:                &sync.Mutex{},
		all:               make(map[crypto.Hash]byte),
		index:             make(map[crypto.Token]*SetOfHashes),
		UpdateCollective:  make(map[crypto.Hash]*PendingUpdate),
		RequestMembership: make(map[crypto.Hash]*PendingRequestMembership),
		RemoveMember:      make(map[crypto.Hash]*PendingRemoveMember),
		Draft:             make(map[crypto.Hash]*Draft),
		Edit:              make(map[crypto.Hash]*Edit),
		CreateBoard:       make(map[crypto.Hash]*PendingBoard),
		UpdateBoard:       make(map[crypto.Hash]*PendingUpdateBoard),
		Pin:               make(map[crypto.Hash]*Pin),
		BoardEditor:       make(map[crypto.Hash]*BoardEditor),
		ReleaseDraft:      make(map[crypto.Hash]*Release),
		ImprintStamp:      make(map[crypto.Hash]*Stamp),
		//react map[crypto.Hash]*
		CreateEvent: make(map[crypto.Hash]*Event),
		CancelEvent: make(map[crypto.Hash]*CancelEvent),
		UpdateEvent: make(map[crypto.Hash]*EventUpdate),
	}
}

type SetOfHashes struct {
	set map[crypto.Hash]struct{}
}

func NewSetOfHashes() *SetOfHashes {
	return &SetOfHashes{
		set: make(map[crypto.Hash]struct{}),
	}
}

func (s *SetOfHashes) Add(hash crypto.Hash) {
	s.set[hash] = struct{}{}
}

func (s *SetOfHashes) Remove(hash crypto.Hash) {
	delete(s.set, hash)
}

func (s *SetOfHashes) All() map[crypto.Hash]struct{} {
	all := make(map[crypto.Hash]struct{})
	if s.set == nil || len(s.set) == 0 {
		return all
	}
	for hash := range s.set {
		all[hash] = struct{}{}
	}
	return all
}

type Proposals struct {
	mu                *sync.Mutex
	all               map[crypto.Hash]byte
	index             map[crypto.Token]*SetOfHashes
	UpdateCollective  map[crypto.Hash]*PendingUpdate
	RequestMembership map[crypto.Hash]*PendingRequestMembership
	RemoveMember      map[crypto.Hash]*PendingRemoveMember
	Draft             map[crypto.Hash]*Draft
	Edit              map[crypto.Hash]*Edit
	CreateBoard       map[crypto.Hash]*PendingBoard
	UpdateBoard       map[crypto.Hash]*PendingUpdateBoard
	Pin               map[crypto.Hash]*Pin
	BoardEditor       map[crypto.Hash]*BoardEditor
	ReleaseDraft      map[crypto.Hash]*Release
	ImprintStamp      map[crypto.Hash]*Stamp
	//react map[crypto.Hash]*
	CreateEvent map[crypto.Hash]*Event
	CancelEvent map[crypto.Hash]*CancelEvent
	UpdateEvent map[crypto.Hash]*EventUpdate
}

func (p *Proposals) GetEvent(hash crypto.Hash) *Event {
	e, _ := p.CreateEvent[hash]
	return e
}

func (p *Proposals) Delete(hash crypto.Hash) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.all, hash)
	delete(p.UpdateCollective, hash)
	delete(p.RequestMembership, hash)
	delete(p.RemoveMember, hash)
	delete(p.Draft, hash)
	delete(p.Edit, hash)
	delete(p.CreateBoard, hash)
	delete(p.UpdateBoard, hash)
	delete(p.Pin, hash)
	delete(p.BoardEditor, hash)
	delete(p.ReleaseDraft, hash)
	delete(p.ImprintStamp, hash)
	//react map[crypto.Hash]*
	delete(p.CreateEvent, hash)
	delete(p.CancelEvent, hash)
	delete(p.UpdateEvent, hash)
	for _, hashes := range p.index {
		hashes.Remove(hash)
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
	p.mu.Lock()
	defer p.mu.Unlock()
	members := c.ListOfTokens()
	for token := range members {
		if _, ok := p.index[token]; !ok {
			p.index[token] = NewSetOfHashes()
		}
		p.index[token].Add(hash)
	}
}

func (p *Proposals) GetVotes(token crypto.Token) map[crypto.Hash]struct{} {
	hashes := p.index[token]
	if hashes == nil {
		return nil
	}
	return hashes.All()
}

func (p *Proposals) AddUpdateCollective(update *PendingUpdate) {
	p.indexHash(update.Collective, update.Hash)
	p.all[update.Hash] = UpdateCollectiveProposal
	p.UpdateCollective[update.Hash] = update
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
	p.Draft[update.DraftHash] = update
}

func (p *Proposals) AddEdit(update *Edit) {
	p.indexHash(update.Draft.Authors, update.Edit)
	p.indexHash(update.Authors, update.Edit)
	p.all[update.Edit] = EditProposal
	p.Edit[update.Edit] = update
}

func (p *Proposals) AddPendingBoard(update *PendingBoard) {
	p.indexHash(update.Board.Collective, update.Hash)
	p.all[update.Hash] = CreateBoardProposal
	p.CreateBoard[update.Hash] = update
}

func (p *Proposals) AddPendingUpdateBoard(update *PendingUpdateBoard) {
	p.indexHash(update.Board.Editors, update.Hash)
	p.all[update.Hash] = UpdateBoardProposal
	p.UpdateBoard[update.Hash] = update
}

func (p *Proposals) AddPin(update *Pin) {
	p.indexHash(update.Board.Editors, update.Hash)
	p.all[update.Hash] = PinProposal
	p.Pin[update.Hash] = update
}

func (p *Proposals) AddBoardEditor(update *BoardEditor) {
	p.indexHash(update.Board.Collective, update.Hash)
	p.all[update.Hash] = BoardEditorProposal
	p.BoardEditor[update.Hash] = update
}

func (p *Proposals) AddRelease(update *Release) {
	p.indexHash(update.Draft.Authors, update.Hash)
	p.all[update.Hash] = ReleaseDraftProposal
	p.ReleaseDraft[update.Hash] = update
}

func (p *Proposals) AddStamp(update *Stamp) {
	p.indexHash(update.Reputation, update.Hash) // reputation aqui é = um membro ou coletivo que vai dar o stamp ??
	p.all[update.Hash] = ImprintStampProposal
	p.ImprintStamp[update.Hash] = update
}

func (p *Proposals) AddEvent(update *Event) {
	p.indexHash(update.Collective, update.Hash)
	p.all[update.Hash] = CreateEventProposal
	p.CreateEvent[update.Hash] = update
}

func (p *Proposals) AddCancelEvent(update *CancelEvent) {
	p.indexHash(update.Event.Collective, update.Hash)
	p.all[update.Hash] = CancelEventProposal
	p.CancelEvent[update.Hash] = update
}

func (p *Proposals) AddEventUpdate(update *EventUpdate) {
	p.indexHash(update.Event.Managers, update.Hash)
	p.all[update.Hash] = UpdateEventProposal
	p.UpdateEvent[update.Hash] = update
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
		proposal = p.UpdateCollective[hash]
	case RequestMembershipProposal:
		proposal = p.RequestMembership[hash]
	case RemoveMemberProposal:
		proposal = p.RemoveMember[hash]
	case DraftProposal:
		proposal = p.Draft[hash]
	case EditProposal:
		proposal = p.Edit[hash]
	case CreateBoardProposal:
		proposal = p.CreateBoard[hash]
	case UpdateBoardProposal:
		proposal = p.UpdateBoard[hash]
	case PinProposal:
		proposal = p.Pin[hash]
	case BoardEditorProposal:
		proposal = p.BoardEditor[hash]
	case ReleaseDraftProposal:
		proposal = p.ReleaseDraft[hash]
	case ImprintStampProposal:
		proposal = p.ImprintStamp[hash]
	case ReactProposal:
		//
	case CreateEventProposal:
		proposal = p.CreateEvent[hash]
	case CancelEventProposal:
		proposal = p.CancelEvent[hash]
	case UpdateEventProposal:
		proposal = p.UpdateEvent[hash]
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
		proposal := p.UpdateCollective[hash]
		return proposal.Collective.Name
	case RemoveMemberProposal:
		proposal := p.RemoveMember[hash]
		return proposal.Collective.Name
	case DraftProposal:
		proposal := p.Draft[hash]
		return proposal.Authors.CollectiveName()
	case EditProposal:
		proposal := p.Edit[hash]
		return proposal.Authors.CollectiveName()
	case CreateBoardProposal:
		proposal := p.CreateBoard[hash]
		return proposal.Board.Collective.Name
	case UpdateBoardProposal:
		proposal := p.UpdateBoard[hash]
		return proposal.Board.Collective.Name
	case PinProposal:
		proposal := p.Pin[hash]
		return proposal.Board.Name
	case BoardEditorProposal:
		proposal := p.BoardEditor[hash]
		return proposal.Board.Collective.Name
	case ReleaseDraftProposal:
		proposal := p.ReleaseDraft[hash]
		return proposal.Draft.Authors.CollectiveName()
	case ImprintStampProposal:
		proposal := p.ImprintStamp[hash]
		return proposal.Reputation.Name
	case ReactProposal:
		//
	case CreateEventProposal:
		proposal := p.CreateEvent[hash]
		return proposal.Collective.Name
	case CancelEventProposal:
		proposal := p.CancelEvent[hash]
		return proposal.Event.Collective.Name
	case UpdateEventProposal:
		proposal := p.UpdateEvent[hash]
		return proposal.Event.Collective.Name
	}
	return ""
}
