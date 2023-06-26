package api

import (
	"time"

	"github.com/lienkolabs/swell/crypto"
	"github.com/lienkolabs/synergy/social/state"
)

type EventsView struct {
	Hash        string
	Description string
	StartAt     time.Time
	Collective  string
	Public      bool
}

type EventVoteAction struct {
	Kind string // create, cancel, update
	// OnBehalfOf string // collective, managers
	Hash string
}

type EventsListView struct {
	Events []EventsView
}

type EventDetailView struct {
	Description  string
	StartAt      time.Time
	EstimatedEnd time.Time
	Collective   string
	Venue        string
	Open         bool
	Public       bool
	Managers     []string
	Votes        []EventVoteAction
	Managing     bool
}

func EventsFromState(state *state.State) EventsListView {
	view := EventsListView{
		Events: make([]EventsView, 0),
	}
	for _, event := range state.Events {
		hash, _ := event.Hash.MarshalText()
		itemView := EventsView{
			Hash:        string(hash),
			Description: event.Description,
			StartAt:     event.StartAt,
			Collective:  event.Collective.Name,
			Public:      event.Public,
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
		Description:  event.Description,
		StartAt:      event.StartAt,
		Collective:   event.Collective.Name,
		EstimatedEnd: event.EstimatedEnd,
		Venue:        event.Venue,
		Open:         event.Open,
		Public:       event.Public,
		Managers:     membersToHandles(event.Managers.ListOfMembers(), s),
		Votes:        make([]EventVoteAction, 0),
		Managing:     event.Managers.IsMember(token),
	}
	pending := s.Proposals.GetVotes(token)
	if len(pending) > 0 {
		for pendingHash := range pending {
			hash, _ := pendingHash.MarshalText()
			vote := EventVoteAction{
				Hash: string(hash),
			}
			switch s.Proposals.Kind(pendingHash) {
			case state.CreateEventProposal:
				// collective vote
				vote.Kind = "Create"
				view.Votes = append(view.Votes, vote)
			case state.UpdateEventProposal:
				// managers vote
				vote.Kind = "Update"
				view.Votes = append(view.Votes, vote)
			case state.CancelEventProposal:
				// managers vote
				vote.Kind = "Cancel"
				view.Votes = append(view.Votes, vote)
			}
		}
	}
	return &view
}

func EventUpdateDetailFromState(s *state.State, hash crypto.Hash, token crypto.Token) *EventDetailView {
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
		Managing:     event.Managers.IsMember(token),
	}
	pending := s.Proposals.GetVotes(token)
	if len(pending) > 0 {
		for pendingHash := range pending {
			hash, _ := pendingHash.MarshalText()
			vote := EventVoteAction{
				Hash: string(hash),
			}
			if vote.Kind != "Update" {
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
