package api

import (
	"time"

	"github.com/lienkolabs/swell/crypto"
	"github.com/lienkolabs/synergy/social/state"
)

type EventsView struct {
	Hash            string
	Live            bool
	PendingApproval bool
	Description     string
	StartAt         time.Time
	Collective      string
	Public          bool
}

type EventVoteAction struct {
	Kind   string // create, cancel, update
	Hash   string
	Update string
}

type EventsListView struct {
	Events []EventsView
}

type EventDetailView struct {
	Live            bool
	PendingApproval bool
	Description     string
	StartAt         time.Time
	EstimatedEnd    time.Time
	Collective      string
	Venue           string
	Open            bool
	Public          bool
	Managers        []string
	Checkedin       []string
	Votes           []EventVoteAction
	Managing        bool
	Hash            string
}

func EventsFromState(state *state.State) EventsListView {
	view := EventsListView{
		Events: make([]EventsView, 0),
	}
	for _, event := range state.Events {
		pending := false
		// event_prop := state.Proposals.GetEvent(event.Hash)
		// if event_prop != nil {
		// 	vote := EventVoteAction{
		// 		Hash: crypto.EncodeHash(event.Hash),
		// 	}
		// 	if vote.Kind == "Create" {
		// 		pending = true
		// 	}
		// }
		itemView := EventsView{
			Hash:            crypto.EncodeHash(event.Hash),
			Live:            event.Live,
			PendingApproval: pending,
			Description:     event.Description,
			StartAt:         event.StartAt,
			Collective:      event.Collective.Name,
			Public:          event.Public,
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
		Live:            event.Live,
		Description:     event.Description,
		StartAt:         event.StartAt,
		Collective:      event.Collective.Name,
		EstimatedEnd:    event.EstimatedEnd,
		Venue:           event.Venue,
		Open:            event.Open,
		Public:          event.Public,
		Checkedin:       make([]string, 0),
		Managers:        membersToHandles(event.Managers.ListOfMembers(), s),
		Votes:           make([]EventVoteAction, 0),
		Managing:        event.Managers.IsMember(token),
		Hash:            crypto.EncodeHash(hash),
		PendingApproval: false,
	}
	for token := range event.Checkin {
		if handle, ok := s.Members[crypto.Hasher(token[:])]; ok {
			view.Checkedin = append(view.Checkedin, handle)
		}
	}
	pending := s.Proposals.GetVotes(token)
	if len(pending) > 0 {
		for pendingHash := range pending {
			vote := EventVoteAction{
				Hash: crypto.EncodeHash(pendingHash),
			}
			switch s.Proposals.Kind(pendingHash) {
			case state.CreateEventProposal:
				// collective vote
				vote.Kind = "Create"
				view.Votes = append(view.Votes, vote)
				view.PendingApproval = true
			case state.UpdateEventProposal:
				// managers vote
				vote.Kind = "Update"
				vote.Update = "Update"
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
		Live:         event.Live,
		Description:  event.Description,
		Collective:   event.Collective.Name,
		EstimatedEnd: event.EstimatedEnd,
		Venue:        event.Venue,
		Open:         event.Open,
		Public:       event.Public,
		Managers:     membersToHandles(event.Managers.ListOfMembers(), s),
		Votes:        make([]EventVoteAction, 0),
		Managing:     event.Managers.IsMember(token),
		Hash:         crypto.EncodeHash(hash),
	}
	pending := s.Proposals.GetVotes(token)
	if len(pending) > 0 {
		for pendingHash := range pending {
			vote := EventVoteAction{
				Hash: crypto.EncodeHash(pendingHash),
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
