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

// Collectives template struct

type CollectivesView struct {
	Name        string
	Description string
}

type CollectivesListView struct {
	Collectives []CollectivesView
}

type CollectiveDetailView struct {
	Name          string
	Description   string
	Majority      int
	SuperMajority int
	Members       []MemberDetailView
}

func ColletivesFromState(state *state.State) CollectivesListView {
	view := CollectivesListView{
		Collectives: make([]CollectivesView, 0),
	}
	for _, collective := range state.Collectives {
		itemView := CollectivesView{
			Name:        collective.Name,
			Description: collective.Description,
		}
		view.Collectives = append(view.Collectives, itemView)
	}
	return view
}

func CollectiveDetailFromState(state *state.State, name string) *CollectiveDetailView {
	collective, ok := state.Collective(name)
	if !ok {
		return nil
	}
	view := CollectiveDetailView{
		Name:          collective.Name,
		Description:   collective.Description,
		Majority:      collective.Policy.Majority,
		SuperMajority: collective.Policy.SuperMajority,
		Members:       make([]MemberDetailView, 0),
	}
	for token, _ := range collective.Members {
		handle, ok := state.Members[crypto.Hasher(token[:])]
		if ok {
			view.Members = append(view.Members, MemberDetailView{handle})
		}
	}
	return &view
}

// Drafts template struct

type DraftsView struct {
	Title       string
	Authors     []string
	Description string
	Keywords    []string
}

type DraftsListView struct {
	Drafts []DraftsView
}

type DraftDetailView struct {
	Title        string
	Description  string
	Keywords     []string
	Content      string
	Authors      []string
	References   []string
	PreviousHash string
}

func membersToHandles(members map[crypto.Token]struct{}, state *state.State) []string {
	handles := make([]string, 0)
	for member, _ := range members {
		handle, ok := state.Members[crypto.Hasher(member[:])]
		if ok {
			handles = append(handles, handle)
		}
	}
	return handles
}

func hashesToString(hashes []crypto.Hash) []string {
	output := make([]string, 0)
	for _, hash := range hashes {
		text, err := hash.MarshalText()
		if err != nil {
			output = append(output, string(text))
		}
	}
	return output
}

func DraftsFromState(state *state.State) DraftsListView {
	view := DraftsListView{
		Drafts: make([]DraftsView, 0),
	}
	for _, draft := range state.Drafts {
		itemView := DraftsView{
			Title:       draft.Title,
			Authors:     membersToHandles(draft.Authors.ListOfMembers(), state),
			Description: draft.Description,
			Keywords:    draft.Keywords,
		}
		view.Drafts = append(view.Drafts, itemView)
	}
	return view
}

func DraftDetailFromState(state *state.State, hash crypto.Hash) *DraftDetailView {
	draft, ok := state.Drafts[hash]
	if !ok {
		return nil
	}
	view := DraftDetailView{
		Title:       draft.Title,
		Description: draft.Description,
		Keywords:    draft.Keywords,
		Authors:     membersToHandles(draft.Authors.ListOfMembers(), state),
		References:  hashesToString(draft.References),
	}
	text, _ := draft.PreviousVersion.DraftHash.MarshalText()
	view.PreviousHash = string(text)
	return &view
}

// Edits template struct

type EditsView struct {
	Title   string
	Authors []string
	Reasons string
}

type EditsListView struct {
	Edits []EditsView
}

func EditsFromState(state *state.State) EditsListView {
	view := EditsListView{
		Edits: make([]EditsView, 0),
	}
	for _, edit := range state.Edits {
		itemView := EditsView{
			Title:   edit.Draft.Title,
			Authors: membersToHandles(edit.Authors.ListOfMembers(), state),
			Reasons: edit.Reasons,
		}
		view.Edits = append(view.Edits, itemView)
	}
	return view
}

// Events template struct

type EventsView struct {
	ID          string
	Description string
	Collective  string
	Public      bool
	StartAt     time.Time
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
}

func EventsFromState(state *state.State) EventsListView {
	view := EventsListView{
		Events: make([]EventsView, 0),
	}
	for _, event := range state.Events {
		itemView := EventsView{
			Description: event.Description,
			Collective:  event.Collective.Name,
			Public:      event.Public,
			StartAt:     event.StartAt,
		}
		view.Events = append(view.Events, itemView)
	}
	return view
}

func EventDetailFromState(state *state.State, hash crypto.Hash) *EventDetailView {
	event, ok := state.Events[hash]
	if !ok {
		return nil
	}
	view := EventDetailView{
		StartAt:      event.StartAt,
		Description:  event.Description,
		Collective:   event.Collective.Name,
		EstimatedEnd: event.EstimatedEnd,
		Venue:        event.Venue,
		Open:         event.Open,
		Public:       event.Public,
		Managers:     membersToHandles(event.Managers.ListOfMembers(), state),
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
