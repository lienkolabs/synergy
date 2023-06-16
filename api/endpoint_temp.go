package api

import (
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/lienkolabs/synergy/social/state"
)

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

// Events template struct

type EventsView struct {
	ID          string
	Description string
	Owner       string
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
			Owner:       event.Collective.Name,
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

// State

func NewStateView(state *state.State) *StateView {
	templates := make(map[string]*template.Template)
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
