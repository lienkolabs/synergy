package api

import (
	"time"

	"github.com/lienkolabs/swell/crypto"
	"github.com/lienkolabs/synergy/social/state"
)

type EventsView struct {
	Hash        string
	Live        bool
	Description string
	StartAt     time.Time
	Collective  string
	Public      bool
}

type EventVoteAction struct {
	Kind   string // create, cancel, update
	Hash   string
	Update string
}

type EventsListView struct {
	Events []EventsView
}

type VoteUpdateEventView struct {
	Description     string
	OldDescription  string
	StartAt         string
	OldStartAt      string
	EstimatedEnd    string
	OldEstimatedEnd string
	Venue           string
	OldVenue        string
	Open            string
	OldOpen         string
	Public          string
	OldPublic       string
	Hash            string
	Reasons         string
	Collective      string
	Managing        bool
	VoteHash        string
}

func yesorno(b *bool) string {
	if b == nil {
		return ""
	}
	if *b {
		return "yes"
	}
	return "no"
}

func EventUpdateFromState(s *state.State, hash crypto.Hash, token crypto.Token) VoteUpdateEventView {
	update, ok := s.Proposals.UpdateEvent[hash]
	if !ok {
		return VoteUpdateEventView{}
	}
	old := update.Event
	vote := VoteUpdateEventView{
		OldDescription:  old.Description,
		OldStartAt:      old.StartAt.String(),
		OldEstimatedEnd: old.EstimatedEnd.String(),
		OldVenue:        old.Venue,
		OldOpen:         yesorno(&old.Open),
		OldPublic:       yesorno(&old.Public),
		Open:            yesorno(update.Open),
		Public:          yesorno(update.Public),
		Hash:            crypto.EncodeHash(old.Hash),
		Reasons:         update.Reasons,
		Collective:      old.Collective.Name,
		VoteHash:        crypto.EncodeHash(hash),
	}
	if update.Description != nil {
		vote.Description = *update.Description
	}
	if update.StartAt != nil {
		vote.StartAt = update.StartAt.String()
	}
	if update.EstimatedEnd != nil {
		vote.StartAt = update.StartAt.String()
	}
	if update.Venue != nil {
		vote.StartAt = *update.Venue
	}
	if old.Managers.IsMember(token) {
		vote.Managing = true
	}
	return vote
}

type EventDetailView struct {
	Live            bool
	Description     string
	StartAt         time.Time
	EstimatedEnd    time.Time
	Collective      string
	Venue           string
	Open            bool
	Public          bool
	ManagerMajority int
	Managers        []MemberDetailView
	Checkedin       []MemberDetailView
	Votes           []EventVoteAction
	Managing        bool
	Hash            string
}

func EventsFromState(state *state.State) EventsListView {
	view := EventsListView{
		Events: make([]EventsView, 0),
	}
	for _, event := range state.Events {
		itemView := EventsView{
			Hash: crypto.EncodeHash(event.Hash),
			Live: event.Live,

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
		Live:            event.Live,
		Description:     event.Description,
		StartAt:         event.StartAt,
		Collective:      event.Collective.Name,
		EstimatedEnd:    event.EstimatedEnd,
		Venue:           event.Venue,
		Open:            event.Open,
		Public:          event.Public,
		Checkedin:       make([]MemberDetailView, 0),
		ManagerMajority: event.Managers.Majority,
		Managers:        make([]MemberDetailView, 0),
		Votes:           make([]EventVoteAction, 0),
		Managing:        event.Managers.IsMember(token),
		Hash:            crypto.EncodeHash(hash),
		// Managers:        membersToHandles(event.Managers.ListOfMembers(), s),
	}
	for token, _ := range event.Managers.ListOfMembers() {
		handle, ok := s.Members[crypto.Hasher(token[:])]
		if ok {
			view.Managers = append(view.Managers, MemberDetailView{handle})
		}
	}
	for token := range event.Checkin {
		if handle, ok := s.Members[crypto.Hasher(token[:])]; ok {
			view.Checkedin = append(view.Checkedin, MemberDetailView{handle})
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
		StartAt:         event.StartAt,
		Live:            event.Live,
		Description:     event.Description,
		Collective:      event.Collective.Name,
		EstimatedEnd:    event.EstimatedEnd,
		Venue:           event.Venue,
		Open:            event.Open,
		Public:          event.Public,
		ManagerMajority: event.Managers.Majority,
		Managers:        make([]MemberDetailView, 0),
		Votes:           make([]EventVoteAction, 0),
		Managing:        event.Managers.IsMember(token),
		Hash:            crypto.EncodeHash(hash),
		// Managers:        membersToHandles(event.Managers.ListOfMembers(), s),
	}
	for token, _ := range event.Managers.ListOfMembers() {
		handle, ok := s.Members[crypto.Hasher(token[:])]
		if ok {
			view.Managers = append(view.Managers, MemberDetailView{handle})
		}
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

// Accept Checkin Event Vote

// type AcceptCheckinEventView struct {
// 	Description   string
// 	StartAt       time.Time
// 	Venue         string
// 	CheckingIn    string
// 	Reasons       string
// 	EventMajority int
// }

// func PendingCheckinEventFromState(state *state.State, hash crypto.Hash) *AcceptCheckinEventView {
// 	pending, ok := state.Proposals.CheckinEvent[hash]
// 	if !ok {
// 		return nil
// 	}
// 	checkin := pending.EventCheckin
// 	view := AcceptCheckinEventView{
// 		Description:   checkin.Description,
// 		StartAt:       checkin.StartAt,
// 		Venue:         checkin.Venue,
// 		Reasons:       checkin.Reasons,
// 		EventMajority: checkin.EventMajority,
// 	}
// 	view.CheckingIn = state.Members[crypto.Hasher(pending.Origin.Author[:])]
// 	return &view
// }
