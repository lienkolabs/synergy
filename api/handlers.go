package api

import (
	// "fmt"
	// "log"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"
	"time"

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

func (a *Attorney) MediaHandler(w http.ResponseWriter, r *http.Request) {
	hashtext := r.URL.Path
	hashtext = strings.Replace(hashtext, "/media/", "", 1)
	hash := crypto.DecodeHash(hashtext)

	file, ok := a.state.Media[hash]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("file not found"))
		return
	}
	title := hashtext
	var ext string
	if edit, ok := a.state.Edits[hash]; ok {
		ext = edit.EditType
	} else if draft, ok := a.state.Drafts[hash]; ok {
		ext = draft.DraftType
	} else if draft, ok := a.state.Proposals.Draft[hash]; ok {
		ext = draft.DraftType
	} else if edit, ok := a.state.Proposals.Edit[hash]; ok {
		ext = edit.EditType
	}
	name := fmt.Sprintf("%v", title, ext)
	//cd := mime.FormatMediaType("attachment", map[string]string{"filename": name})
	//w.Header().Set("Content-Disposition", cd)
	//w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeContent(w, r, name, time.Now(), bytes.NewReader(file))
}

func (a *Attorney) NewEditHandler(w http.ResponseWriter, r *http.Request) {
	var hash crypto.Hash
	if err := r.ParseForm(); err == nil {
		hash = crypto.DecodeHash(r.FormValue("draftHash"))
		fmt.Println(crypto.EncodeHash(hash))
		t := a.templates["edit"]
		if view := NewEdit(a.state, hash); view != nil {
			if err := t.Execute(w, view); err != nil {
				log.Println(err)
			}
			return
		}
	}
}

func (a *Attorney) NewDraft2Handler(w http.ResponseWriter, r *http.Request) {
	var hash crypto.Hash
	if err := r.ParseForm(); err == nil {
		hash = crypto.DecodeHash(r.FormValue("previousVersion"))
	}
	t := a.templates["newdraft2"]
	view := NewDraftVerion(a.state, hash)
	if err := t.Execute(w, view); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) EditViewHandler(w http.ResponseWriter, r *http.Request) {
	hash := getHash(r.URL.Path, "/editview/")
	t := a.templates["editview"]
	view := EditDetailFromState(a.state, hash, a.author)
	if err := t.Execute(w, view); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) NewDraftHandler(w http.ResponseWriter, r *http.Request) {
	var hash crypto.Hash
	if err := r.ParseForm(); err == nil {
		hash = crypto.DecodeHash(r.FormValue("previousVersion"))
	}
	t := a.templates["newdraft"]
	view := NewDraftVerion(a.state, hash)
	if err := t.Execute(w, view); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) BoardsHandler(w http.ResponseWriter, r *http.Request) {
	t := a.templates["boards"]
	view := BoardsFromState(a.state)
	if err := t.Execute(w, view); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) BoardHandler(w http.ResponseWriter, r *http.Request) {
	boardName := r.URL.Path
	boardName = strings.Replace(boardName, "/board/", "", 1)
	t := a.templates["board"]
	view := BoardDetailFromState(a.state, boardName, a.author)
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
	view := CollectiveDetailFromState(a.state, name, a.author)
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
	hash := getHash(r.URL.Path, "/edits/")
	t := a.templates["edits"]
	view := EditsFromState(a.state, hash)
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

func getHash(path string, root string) crypto.Hash {
	path = strings.Replace(path, root, "", 1)
	hash := crypto.DecodeHash(path)
	return hash
}

func (a *Attorney) RequestMemberShipVoteHandler(w http.ResponseWriter, r *http.Request) {
	hash := getHash(r.URL.Path, "/requestmembership/")
	view := RequestMembershipFromState(a.state, hash)
	t := a.templates["requestmembershipvote"]
	if err := t.Execute(w, view); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) VotesHandler(w http.ResponseWriter, r *http.Request) {
	view := VotesFromState(a.state, a.author)
	t := a.templates["votes"]
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

func (a *Attorney) CreateBoardHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	board := r.FormValue("collective")
	if _, ok := a.state.Collective(board); !ok {
		fmt.Fprintf(w, "Collective not found")
		return
	}
	t := a.templates["createboard"]
	if err := t.Execute(w, board); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) VoteCreateBoardHandler(w http.ResponseWriter, r *http.Request) {
	hash := getHash(r.URL.Path, "/votecreateboard/")
	t := a.templates["votecreateboard"]
	view := PendingBoardFromState(a.state, hash)
	if err := t.Execute(w, view); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) UpdateCollectiveHandler(w http.ResponseWriter, r *http.Request) {
	collective := strings.Replace(r.URL.Path, "/updatecollective/", "", 1)
	t := a.templates["updatecollective"]
	view := CollectiveToUpdateFromState(a.state, collective)
	if err := t.Execute(w, view); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) VoteUpdateCollectiveHandler(w http.ResponseWriter, r *http.Request) {
	hash := getHash(r.URL.Path, "/voteupdatecollective/")
	t := a.templates["voteupdatecollective"]
	view := CollectiveUpdateFromState(a.state, hash, a.author)
	if err := t.Execute(w, view); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) UpdateBoardHandler(w http.ResponseWriter, r *http.Request) {
	board := strings.Replace(r.URL.Path, "/updateboard/", "", 1)
	t := a.templates["updateboard"]
	view := BoardToUpdateFromState(a.state, board)
	if err := t.Execute(w, view); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) VoteUpdateBoardHandler(w http.ResponseWriter, r *http.Request) {
	hash := getHash(r.URL.Path, "/votecreateboard/")
	t := a.templates["voteupdateboard"]
	view := BoardUpdateFromState(a.state, hash)
	if err := t.Execute(w, view); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) UpdateEventHandler(w http.ResponseWriter, r *http.Request) {
	hash := getHash(r.URL.Path, "/updateevent/")
	t := a.templates["updateevent"]
	view := EventUpdateDetailFromState(a.state, hash, a.author)
	if err := t.Execute(w, view); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) CreateEventHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	board := r.FormValue("collective")
	if _, ok := a.state.Collective(board); !ok {
		fmt.Fprintf(w, "Collective not found")
		return
	}
	t := a.templates["createevent"]
	if err := t.Execute(w, board); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) VoteUpdateEventHandler(w http.ResponseWriter, r *http.Request) {
	hash := getHash(r.URL.Path, "/voteupdateevent/")
	t := a.templates["voteupdateevent"]
	view := EventUpdateFromState(a.state, hash, a.author)
	if err := t.Execute(w, view); err != nil {
		log.Println(err)
	}
}
