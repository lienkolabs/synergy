package api

import (
	"time"

	"github.com/lienkolabs/swell/crypto"
	"github.com/lienkolabs/synergy/social/state"
)

// Boards template struct

type BoardsView struct {
	Name       string
	Hash       string
	Collective string
	Keywords   []string
}

type BoardsListView struct {
	Boards []BoardsView
}

type BoardDetailView struct {
	Name        string
	Description string
	Collective  string
	Keywords    []string
	PinMajority int
	Editors     []string
	Drafts      []DraftsView
}

func BoardsFromState(state *state.State) BoardsListView {
	view := BoardsListView{
		Boards: make([]BoardsView, 0),
	}
	for _, board := range state.Boards {
		hash, _ := board.Hash.MarshalText()
		itemView := BoardsView{
			Name:       board.Name,
			Hash:       string(hash),
			Collective: board.Collective.Name,
			Keywords:   board.Keyword,
		}
		view.Boards = append(view.Boards, itemView)
	}
	return view
}

func BoardDetailFromState(state *state.State, hash crypto.Hash) *BoardDetailView {
	board, ok := state.Boards[hash]
	if !ok {
		return nil
	}
	view := BoardDetailView{
		Name:        board.Name,
		Description: board.Description,
		Collective:  board.Collective.Name,
		Keywords:    board.Keyword,
		PinMajority: board.Editors.Majority,
		Editors:     make([]string, 0),
		Drafts:      make([]DraftsView, 0),
	}
	for token, _ := range board.Editors.Members {
		handle, ok := state.Members[crypto.Hasher(token[:])]
		if ok {
			view.Editors = append(view.Editors, handle)
		}
	}
	for _, d := range board.Pinned {
		draftView := DraftsView{
			Title:       d.Title,
			Authors:     make([]string, 0),
			Description: d.Description,
			Keywords:    d.Keywords,
		}
		view.Drafts = append(view.Drafts, draftView)
	}

	return &view
}

// Events template struct

type EventsView struct {
	ID          string
	Hash        string
	Description string
	Collective  string
	Public      bool
	StartAt     time.Time
}

type EventVoteAction struct {
	Kind       string // create, cancel, update
	OnBehalfOf string // collective, managers
	Hash       string
}

type EventsListView struct {
	Events []EventsView
}

type EventDetailView struct {
	StartAt      time.Time
	Description  string
	Collective   string
	EstimatedEnd time.Time
	Venue        string
	Open         bool
	Public       bool
	Managers     []string
	Votes        []EventVoteAction
}

func EventsFromState(state *state.State) EventsListView {
	view := EventsListView{
		Events: make([]EventsView, 0),
	}
	for _, event := range state.Events {
		hash, _ := event.Hash.MarshalText()
		itemView := EventsView{
			Description: event.Description,
			Hash:        string(hash),
			Collective:  event.Collective.Name,
			Public:      event.Public,
			StartAt:     event.StartAt,
		}
		view.Events = append(view.Events, itemView)
	}
	return view
}

func EventDetailFromState(s *state.State, hash crypto.Hash, token crypto.Token) *EventDetailView {
	event, ok := s.Events[hash]
	if !ok {
		event = s.Proposals.GetEvent(hash)
		if event == nil {
			return nil
		}
	}
	view := EventDetailView{
		StartAt:      event.StartAt,
		Description:  event.Description,
		Collective:   event.Collective.Name,
		EstimatedEnd: event.EstimatedEnd,
		Venue:        event.Venue,
		Open:         event.Open,
		Public:       event.Public,
		Managers:     membersToHandles(event.Managers.ListOfMembers(), s),
		Votes:        make([]EventVoteAction, 0),
	}
	pending := s.Proposals.GetVotes(token)
	if len(pending) > 0 {
		for pendingHash := range pending {
			hash, _ := pendingHash.MarshalText()
			vote := EventVoteAction{
				OnBehalfOf: s.Proposals.OnBehalfOf(pendingHash),
				Hash:       string(hash),
			}
			switch s.Proposals.Kind(pendingHash) {
			case state.CreateEventProposal:
				vote.Kind = "Create"
			case state.UpdateEventProposal:
				vote.Kind = "Update"
			case state.CancelEventProposal:
				vote.Kind = "Cancel"
			}
			if vote.Kind != "" {
				view.Votes = append(view.Votes, vote)
			}
		}
	}
	return &view
}

// Members template struct

type MembersView struct {
	Hash   string
	Handle string
}

type MembersListView struct {
	Members []MembersView
}

type MemberDetailView struct {
	Handle string
}

func MembersFromState(state *state.State) MembersListView {
	view := MembersListView{
		Members: make([]MembersView, 0),
	}
	for hash, member := range state.Members {
		hashText, _ := hash.MarshalText()
		itemView := MembersView{
			Hash:   string(hashText),
			Handle: member,
		}
		view.Members = append(view.Members, itemView)
	}
	return view
}

func MemberDetailFromState(state *state.State, handle string) *MemberDetailView {
	_, ok := state.MembersIndex[handle]
	if !ok {
		return nil
	}
	view := MemberDetailView{
		Handle: handle,
	}
	return &view
}
