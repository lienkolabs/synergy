package api

import (
	"log"
	"net/http"
	"text/template"

	"github.com/lienkolabs/synergy/social/state"
)

type MembersView struct {
	Hash   string
	Handle string
}

type MembersListView struct {
	Members []MembersView
}

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

func NewStateView(state *state.State) *StateView {
	templates := make(map[string]*template.Template)
	if t, err := template.ParseFiles("./api/templates/members.html"); err != nil {
		log.Fatal(err)
	} else {
		templates["members"] = t
	}

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
