package api

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/lienkolabs/swell/crypto"
	"github.com/lienkolabs/synergy/social/state"
)

// State

var templatesNames = []string{
	"boards", "board", "collectives", "collective", "draft", "drafts", "edits", "events",
	"event", "member", "members",
}

func launchTemplate(name string) *template.Template {
	filePath := fmt.Sprintf("./api/templates/%v.html", name)
	t, err := template.ParseFiles(filePath)
	if err != nil {
		log.Fatal(err)
	}
	return t
}

func NewStateView(state *state.State) *StateView {
	templates := make(map[string]*template.Template)
	for _, name := range templatesNames {
		t := launchTemplate(name)
		if t != nil {
			templates[name] = t
		}
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

func (a *Attorney) BoardsHandler(w http.ResponseWriter, r *http.Request) {
	t := a.templates["boards"]
	view := BoardsFromState(a.state)
	t.Execute(w, view)
}

func (a *Attorney) BoardHandler(w http.ResponseWriter, r *http.Request) {
	boardName := r.URL.Path
	boardName = strings.Replace(boardName, "/board/", "", 1)
	t := a.templates["board"]
	view := BoardDetailFromState(a.state, boardName)
	t.Execute(w, view)
}

func (a *Attorney) CollectivesHandler(w http.ResponseWriter, r *http.Request) {
	t := a.templates["colletives"]
	view := ColletivesFromState(a.state)
	t.Execute(w, view)
}

func (a *Attorney) CollectiveHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path
	name = strings.Replace(name, "/collective/", "", 1)
	t := a.templates["colletive"]
	view := CollectiveDetailFromState(a.state, name)
	t.Execute(w, view)
}

func (a *Attorney) DraftsHandler(w http.ResponseWriter, r *http.Request) {
	t := a.templates["drafts"]
	view := DraftsFromState(a.state)
	t.Execute(w, view)
}

func (a *Attorney) DraftHandler(w http.ResponseWriter, r *http.Request) {
	hashEncoded := r.URL.Path
	hashEncoded = strings.Replace(hashEncoded, "/draft/", "", 1)
	var hash crypto.Hash
	if hash.UnmarshalText([]byte(hashEncoded)) != nil {
		// REPONSE ERRROR
	}
	t := a.templates["draft"]
	view := DraftDetailFromState(a.state, hash)
	t.Execute(w, view)
}

func (a *Attorney) EditsHandler(w http.ResponseWriter, r *http.Request) {
	t := a.templates["edits"]
	view := EditsFromState(a.state)
	t.Execute(w, view)
}

func (a *Attorney) EventsHandler(w http.ResponseWriter, r *http.Request) {
	t := a.templates["events"]
	view := EventsFromState(a.state)
	t.Execute(w, view)
}

func (a *Attorney) EventHandler(w http.ResponseWriter, r *http.Request) {
	hashEncoded := r.URL.Path
	hashEncoded = strings.Replace(hashEncoded, "/event/", "", 1)
	var hash crypto.Hash
	if hash.UnmarshalText([]byte(hashEncoded)) != nil {
		// REPONSE ERRROR
	}
	t := a.templates["event"]
	view := EventDetailFromState(a.state, hash)
	t.Execute(w, view)
}

func (a *Attorney) MembersHandler(w http.ResponseWriter, r *http.Request) {
	t := a.templates["members"]
	view := MembersFromState(a.state)
	t.Execute(w, view)
}

func (a *Attorney) MemberHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path
	name = strings.Replace(name, "/member/", "", 1)
	t := a.templates["member"]
	view := MemberDetailFromState(a.state, name)
	t.Execute(w, view)
}
