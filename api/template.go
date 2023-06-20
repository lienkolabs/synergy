package api

import (
	// "fmt"
	// "log"
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
	t := a.templates["collectives"]
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
