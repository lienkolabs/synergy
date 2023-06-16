package api

import (
	"crypto"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/lienkolabs/synergy/social/state"
)

// Boards template struct

type BoardsView struct {
	Name        string
	Description string
	Keywords    []string
	PinMajority int
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
}

func BoardsFromState(state *state.State) BoardsListView {
	view := BoardsListView{
		Boards: make([]BoardsView, 0),
	}
	for _, board := range state.Boards {
		itemView := BoardsView{
			Name:        board.Name,
			Description: board.Description,
			Keywords:    board.Keyword,
			PinMajority: board.Editors.Majority,
		}
		view.Boards = append(view.Boards, itemView)
	}
	return view
}

func BoardsDetailFromState(state *state.State, name string) *BoardDetailView {
	board, ok := state.Board(name)
	if !ok {
		return nil
	}
	view := BoardDetailView{
		Name:        board.Name,
		Description: board.Description,
		Collective:  board.Collective.Name,
		Keywords:    board.Keyword,
		PinMajority: board.Editors.Majority,
		Editors:     board.Editors.Members.Handle,
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
	}
	return &view
}

// Drafts template struct

type DraftsView struct {
	Title       string
	Author      string
	CoAuthors   []string
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
	Author       string
	CoAuthors    []string
	References   []string
	PreviousHash crypto.Hash
}

func DraftsFromState(state *state.State) DraftsListView {
	view := DraftsListView{
		Drafts: make([]DraftsView, 0),
	}
	for _, draft := range state.Drafts {
		itemView := DraftsView{
			Title:        draft.Title,
			Description:  draft.Description,
			Keywords:     draft.Keywords,
			Content:      draft.Content,
			Author:       draft.Author,
			CoAuthors:    draft.Authors,
			References:   draft.References,
			PreviousHash: draft.PreviousVersion,
		}
		view.Drafts = append(view.Drafts, itemView)
	}
	return view
}

func DraftDetailFromState(state *state.State, title string) *DraftDetailView {
	draft, ok := state.Drafts(title)
	if !ok {
		return nil
	}
	view := DraftDetailView{
		Title:       draft.Title,
		Author:      draft.Author,
		CoAuthors:   draft.CoAuthors,
		Description: draft.Description,
		Keywords:    draft.Keywords,
	}
	return &view
}

// Edits template struct

type EditsView struct {
	Title     string
	Author    string
	CoAuthors []string
	Reasons   string
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
			Title:     edit.Title,
			Author:    edit.Author,
			CoAuthors: edit.Authors,
			Reasons:   edit.Reasons,
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
			ID:          event.ID,
			Description: event.Description,
			Collective:  event.Collective.Name,
			Public:      event.Public,
			StartAt:     event.StartAt,
		}
		view.Events = append(view.Events, itemView)
	}
	return view
}

func EventDetailFromState(state *state.State, ID string) *EventDetailView {
	event, ok := state.Events(ID)
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
		Managers:     event.Managers,
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
	About  string
	Hash   string
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

func MemberDetailFromState(state *state.State, hash string) *MemberDetailView {
	members := state.Members
	member := ""
	for mhash, currentmember := range members {
		if hash == string(mhash[0]) {
			member = currentmember
		}
		if member != "" {
			break
		}
	}
	if member == "" {
		return nil
	}
	view := MemberDetailView{
		Handle: member.Handle,
		About:  member.About,
		Token:  member.Token,
	}
	return &view
}

// State

func NewStateView(state *state.State) *StateView {
	templates := make(map[string]*template.Template)

	// Board
	if t, err := template.ParseFiles("./api/templates/boards.html"); err != nil {
		log.Fatal(err)
	} else {
		templates["boards"] = t
	}
	if eventTemplate, err := template.ParseFiles("./api/templates/board.html"); err != nil {
		log.Fatal(err)
	} else {
		templates["board"] = eventTemplate
	}

	// Collectives
	if collectivesTemplate, err := template.ParseFiles("./api/templates/collectives.html"); err != nil {
		log.Fatal(err)
	} else {
		templates["collectives"] = collectivesTemplate
	}
	if collectiveTemplate, err := template.ParseFiles("./api/templates/collective.html"); err != nil {
		log.Fatal(err)
	} else {
		templates["collective"] = collectiveTemplate
	}

	// Drafts
	if draftsTemplate, err := template.ParseFiles("./api/templates/drafts.html"); err != nil {
		log.Fatal(err)
	} else {
		templates["drafts"] = draftsTemplate
	}
	if draftTemplate, err := template.ParseFiles("./api/templates/draft.html"); err != nil {
		log.Fatal(err)
	} else {
		templates["draft"] = draftTemplate
	}

	// Events
	if t, err := template.ParseFiles("./api/templates/events.html"); err != nil {
		log.Fatal(err)
	} else {
		templates["events"] = t
	}
	if eventTemplate, err := template.ParseFiles("./api/templates/event.html"); err != nil {
		log.Fatal(err)
	} else {
		templates["members"] = eventTemplate
	}

	// Members
	if t, err := template.ParseFiles("./api/templates/members.html"); err != nil {
		log.Fatal(err)
	} else {
		templates["members"] = t
	}
	if memberTemplate, err := template.ParseFiles("./api/templates/member.html"); err != nil {
		log.Fatal(err)
	} else {
		templates["members"] = memberTemplate
	}

	return &StateView{
		State:     state,
		Templates: templates,
	}
}

type StateView struct {
	State     *state.State
	Templates map[string]*template.Template
}

// endpoint "/collectives" vai ser respondido por esta função
func (s *StateView) CollectivesHandler(w http.ResponseWriter, r *http.Request) {
	t := s.Templates["colletives"]
	view := ColletivesFromState(s.State)
	t.Execute(w, view)
}

func (s *StateView) MembersHandler(w http.ResponseWriter, r *http.Request) {
	t := s.Templates["members"]
	view := MembersFromState(s.State)
	t.Execute(w, view)
}

/*
	se a é do tipo int entâo &a é do tipo *int
	se b é do tipo *int então *b é do tipo int
*/
