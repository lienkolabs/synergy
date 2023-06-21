package api

import (
	// "fmt"
	// "log"
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/lienkolabs/swell/crypto"
	"github.com/lienkolabs/synergy/social/state"
)

// State

// var templatesNames = []string{
// 	"boards", "board", "collectives", "collective", "draft", "drafts", "edits", "events",
// 	"event", "member", "members",
// }

type StateView struct {
	State     *state.State
	Templates map[string]*template.Template
}

func (a *Attorney) BoardsHandler(w http.ResponseWriter, r *http.Request) {
	t := a.templates["boards"]
	view := BoardsFromState(a.state)
	if err := t.Execute(w, view); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) BoardHandler(w http.ResponseWriter, r *http.Request) {
	boardHash := r.URL.Path
	boardHash = strings.Replace(boardHash, "/board/", "", 1)
	t := a.templates["board"]
	hash := crypto.DecodeHash(boardHash)
	fmt.Println(hash, boardHash)
	view := BoardDetailFromState(a.state, hash)
	if view == nil {
		w.Write([]byte("board not found"))
	} else if err := t.Execute(w, view); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) CollectivesHandler(w http.ResponseWriter, r *http.Request) {
	t := a.templates["collectives"]
	view := ColletivesFromState(a.state)
	if err := t.Execute(w, view); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) CollectiveHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path
	name = strings.Replace(name, "/collective/", "", 1)
	t := a.templates["collective"]
	view := CollectiveDetailFromState(a.state, name)
	if err := t.Execute(w, view); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) DraftsHandler(w http.ResponseWriter, r *http.Request) {
	t := a.templates["drafts"]
	view := DraftsFromState(a.state)
	if err := t.Execute(w, view); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) DraftHandler(w http.ResponseWriter, r *http.Request) {
	hashEncoded := r.URL.Path
	hashEncoded = strings.Replace(hashEncoded, "/draft/", "", 1)
	hash := crypto.DecodeHash(hashEncoded)
	t := a.templates["draft"]
	view := DraftDetailFromState(a.state, hash, a.author)
	if err := t.Execute(w, view); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) EditsHandler(w http.ResponseWriter, r *http.Request) {
	t := a.templates["edits"]
	view := EditsFromState(a.state)
	if err := t.Execute(w, view); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) EventsHandler(w http.ResponseWriter, r *http.Request) {
	t := a.templates["events"]
	view := EventsFromState(a.state)
	if err := t.Execute(w, view); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) EventHandler(w http.ResponseWriter, r *http.Request) {
	hashEncoded := r.URL.Path
	hashEncoded = strings.Replace(hashEncoded, "/event/", "", 1)
	hash := crypto.DecodeHash(hashEncoded)
	t := a.templates["event"]
	view := EventDetailFromState(a.state, hash, a.author)
	if err := t.Execute(w, view); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) MembersHandler(w http.ResponseWriter, r *http.Request) {
	t := a.templates["members"]
	view := MembersFromState(a.state)
	if err := t.Execute(w, view); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) MemberHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path
	name = strings.Replace(name, "/member/", "", 1)
	t := a.templates["member"]
	view := MemberDetailFromState(a.state, name)
	if err := t.Execute(w, view); err != nil {
		log.Println(err)
	}
}
