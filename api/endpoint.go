package api

import (
	"fmt"
	"strings"
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
	Hash        string
	Description string
	Keywords    []string
}

type DraftsListView struct {
	Drafts []DraftsView
}

type AuthorDetail struct {
	Name       string
	Collective bool
}

type ReferenceDetail struct {
	Title  string
	Author string
	Date   string
}

// co-autor, stamp, pin, version, release
type DraftVoteAction struct {
	Kind       string
	OnBehalfOf string // collective or board editor
	Hash       string
}

type DraftDetailView struct {
	Title       string
	Date        string
	Description string
	Keywords    []string
	Hash        string
	//Content      string
	Authors      []AuthorDetail
	References   []ReferenceDetail
	PreviousHash string
	Pinned       []string
	Edited       bool
	Released     bool
	Stamps       []string
	Votes        []DraftVoteAction
	Policy       Policy
	Authorship   bool
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

func PinList(pin []*state.Board) []string {
	list := make([]string, 0)
	if len(pin) == 0 {
		return list
	}
	for _, p := range pin {
		list = append(list, p.Name)
	}
	return list
}

func StampList(stamps []*state.Collective) []string {
	list := make([]string, 0)
	if len(stamps) == 0 {
		return list
	}
	for _, p := range stamps {
		list = append(list, p.Name)
	}
	return list
}

func References(r []crypto.Hash, s *state.State) []ReferenceDetail {
	references := make([]ReferenceDetail, 0)
	for _, hash := range r {
		if draft, ok := s.Drafts[hash]; ok {
			reference := ReferenceDetail{
				Title:  draft.Title,
				Author: authorsEtAll(draft.Authors, s),
				Date:   fmt.Sprintf("%v", draft.Date.Year()),
			}
			references = append(references, reference)
		}
	}
	return references
}

func authorsEtAll(c state.Consensual, s *state.State) string {
	authors := AuthorList(c, s)
	if len(authors) == 0 {
		return ""
	}
	N := len(authors)
	tail := ""
	if len(authors) > 3 {
		N = 3
		tail = " et al."
	}
	authorlist := make([]string, N)
	for n := 0; n < N; n++ {
		authorlist[n] = authors[n].Name
	}

	return fmt.Sprintf("%v%v", strings.Join(authorlist, ","), tail)
}

func AuthorList(c state.Consensual, s *state.State) []AuthorDetail {
	if c == nil {
		return []AuthorDetail{}
	}
	name := c.CollectiveName()
	if name == "" {
		author := AuthorDetail{
			Name:       name,
			Collective: true,
		}
		return []AuthorDetail{author}
	}
	authors := make([]AuthorDetail, 0)
	for token, _ := range c.ListOfMembers() {
		if handle, ok := s.Members[crypto.Hasher(token[:])]; ok {
			authors = append(authors, AuthorDetail{Name: handle})
		}
	}
	return authors
}

func DraftsFromState(state *state.State) DraftsListView {
	view := DraftsListView{
		Drafts: make([]DraftsView, 0),
	}
	for _, draft := range state.Drafts {
		hash, _ := draft.DraftHash.MarshalText()
		itemView := DraftsView{
			Title:       draft.Title,
			Hash:        string(hash),
			Authors:     membersToHandles(draft.Authors.ListOfMembers(), state),
			Description: draft.Description,
			Keywords:    draft.Keywords,
		}
		view.Drafts = append(view.Drafts, itemView)
	}
	return view
}

func DraftDetailFromState(s *state.State, hash crypto.Hash, token crypto.Token) *DraftDetailView {
	draft, ok := s.Drafts[hash]
	if !ok {
		return nil
	}
	view := DraftDetailView{
		Title:       draft.Title,
		Description: draft.Description,
		Keywords:    draft.Keywords,
		Authors:     AuthorList(draft.Authors, s),
		References:  References(draft.References, s),
		Pinned:      PinList(draft.Pinned),
		Stamps:      StampList(draft.Stamps),
		Votes:       make([]DraftVoteAction, 0),
		Authorship:  draft.Authors.IsMember(token),
	}
	pending := s.Proposals.GetVotes(token)
	if len(pending) > 0 {
		for pendingHash := range pending {
			hash, _ := pendingHash.MarshalText()
			vote := DraftVoteAction{
				OnBehalfOf: s.Proposals.OnBehalfOf(pendingHash),
				Hash:       string(hash),
			}
			switch s.Proposals.Kind(pendingHash) {
			case state.DraftProposal:
				vote.Kind = "Authorship"
			case state.PinProposal:
				vote.Kind = "Pin"
			case state.ImprintStampProposal:
				vote.Kind = "Stamp"
			}
			if vote.Kind != "" {
				view.Votes = append(view.Votes, vote)
			}
		}
	}
	if _, ok := s.Releases[draft.DraftHash]; ok {
		view.Released = true
	}
	if len(draft.Edits) > 0 {
		view.Edited = true
	}
	if draft.PreviousVersion != nil {
		text, _ := draft.PreviousVersion.DraftHash.MarshalText()
		view.PreviousHash = string(text)
	}
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
	Hash        string
	Description string
	Collective  string
	Public      bool
	StartAt     time.Time
}

type EventVoteAction struct {
	Kind       string // create, update, cancel
	OnBehalfOf string // collective or managers
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

func EventDetailFromState(state *state.State, hash crypto.Hash, token crypto.Token) *EventDetailView {
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
		Votes:        make([]EventVoteAction, 0),
	}
	pending := state.Proposals.GetVotes(token)
	if len(pending) > 0 {
		for pendingHash := range pending {
			hash, _ := pendingHash.MarshalText()
			vote := EventVoteAction{
				OnBehalfOf: state.Proposals.OnBehalfOf(pendingHash),
				Hash:       string(hash),
			}
			switch state.Proposals.Kind(pendingHash) {
			case state.EventProposal:
				vote.Kind = "Create"
			case state.EventUpdate:
				vote.Kind = "Update"
			case state.EventCancelation:
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

// Votes template struct

type VotesView struct {
	Token crypto.Token
	Hash  string
}

type VotesListView struct {
	Votes []VotesView
}

type VoteDetailView struct {
	Hash string
}

// func VotesFromState(state *state.State) VotesListView {
// 	view := VotesListView{
// 		Votes: make([]VotesView, 0),
// 	}
// 	for hash, _ := range state.Proposals.GetVotes() {
// 		hashText, _ := hash.MarshalText()
// 		itemView := VotesView{
// 			Hash: string(hashText),
// 		}
// 		view.Votes = append(view.Votes, itemView)
// 	}
// 	return view
// }

// func VoteDetailFromState(state *state.State, hash string) *VoteDetailView {
// 	ok := state.Vote(hash)
// 	if !ok {
// 		return nil
// 	}
// 	view := VoteDetailView{
// 		Hash: hash,
// 	}
// 	return &view
// }
