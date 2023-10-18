package api

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/lienkolabs/breeze/crypto"
	"github.com/lienkolabs/synergy/social/actions"
)

func (a *AttorneyGeneral) CredentialsHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	handle := r.FormValue("handle")
	password := r.FormValue("password")
	cookie := a.CreateSession(handle, password)
	if cookie == "" {
		if err := a.templates.ExecuteTemplate(w, "login.html", "invalid credentials"); err != nil {
			log.Println(err)
		}
	} else {
		http.SetCookie(w, newCookie(cookie))
		header := HeaderInfo{
			UserHandle: handle,
		}
		if err := a.templates.ExecuteTemplate(w, "main.html", header); err != nil {
			fmt.Println("-----------------")
			log.Println(err)
		}
	}
}

func (a *AttorneyGeneral) ApiHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	var actionArray []actions.Action
	var err error
	author := a.Author(r)
	switch r.FormValue("action") {
	case "BoardEditor":
		actionArray, err = BoardEditorForm(r, a.state.MembersIndex, author).ToAction()
	case "CancelEvent":
		actionArray, err = CancelEventForm(r).ToAction()
	case "CheckinEvent":
		actionArray, err = CheckinEventForm(r, a.ephemeralpub).ToAction()
	case "CreateBoard":
		actionArray, err = CreateBoardForm(r).ToAction()
	case "CreateCollective":
		actionArray, err = CreateCollectiveForm(r).ToAction()
	case "CreateEvent":
		actionArray, err = CreateEventForm(r, a.state.MembersIndex, author).ToAction()
	case "GreetCheckinEvent":
		actionArray, err = GreetCheckinEventForm(r, a.state.MembersIndex).ToAction()
	case "ImprintStamp":
		actionArray, err = ImprintStampForm(r).ToAction()
	case "Pin":
		actionArray, err = PinForm(r).ToAction()
	case "React":
		actionArray, err = ReactForm(r).ToAction()
	case "Release":
		actionArray, err = ReleaseDraftForm(r).ToAction()
	case "RemoveMember":
		actionArray, err = RemoveMemberForm(r, a.state.MembersIndex).ToAction()
	case "RequestMembership":
		actionArray, err = RequestMembershipForm(r).ToAction()
	case "UpdateBoard":
		actionArray, err = UpdateBoardForm(r).ToAction()
	case "UpdateCollective":
		actionArray, err = UpdateCollectiveForm(r).ToAction()
	case "UpdateEvent":
		actionArray, err = UpdateEventForm(r, a.state.MembersIndex).ToAction()
	case "Vote":
		actionArray, err = VoteForm(r).ToAction()
	}
	if err == nil && len(actionArray) > 0 {
		a.Send(actionArray, author)
	}
	redirect := fmt.Sprintf("/%v", r.FormValue("redirect"))
	http.Redirect(w, r, redirect, http.StatusSeeOther)
}

func (a *AttorneyGeneral) MainHandler(w http.ResponseWriter, r *http.Request) {
	header := HeaderInfo{
		UserHandle: a.Handle(r),
	}
	if err := a.templates.ExecuteTemplate(w, "main.html", header); err != nil {
		log.Println(err)
	}
}

func (a *AttorneyGeneral) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if err := a.templates.ExecuteTemplate(w, "login.html", ""); err != nil {
		log.Println(err)
	}
}

func (a *AttorneyGeneral) CreateCollectiveHandler(w http.ResponseWriter, r *http.Request) {
	head := HeaderInfo{
		Active:   "CreateCollective",
		Path:     "venture >",
		EndPath:  "create collective",
		Section:  "venture",
		UserName: a.Handle(r),
	}
	if err := a.templates.ExecuteTemplate(w, "createcollective.html", head); err != nil {
		log.Println(err)
	}
}

func (a *AttorneyGeneral) NewDraft2Handler(w http.ResponseWriter, r *http.Request) {
	var hash crypto.Hash
	if err := r.ParseForm(); err == nil {
		hash = crypto.DecodeHash(r.FormValue("previousVersion"))
	}
	view := NewDraftVersion(a.state, hash)
	if view != nil {
		view.Head.UserHandle = a.Handle(r)
		if err := a.templates.ExecuteTemplate(w, "newdraft2.html", view); err != nil {
			log.Println(err)
		} else {
			return
		}
	}
	head := HeaderInfo{
		Error:      "draft not found",
		UserHandle: a.Handle(r),
	}
	if err := a.templates.ExecuteTemplate(w, "main.html", head); err != nil {
		log.Println(err)
	}
}

func (a *AttorneyGeneral) EditViewHandler(w http.ResponseWriter, r *http.Request) {
	author := a.Author(r)
	hash := getHash(r.URL.Path, "/editview/")
	view := EditDetailFromState(a.state, a.indexer, hash, author)
	if err := a.templates.ExecuteTemplate(w, "editview.html", view); err != nil {
		log.Println(err)
	}
}

func (a *AttorneyGeneral) MediaHandler(w http.ResponseWriter, r *http.Request) {
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

func (a *AttorneyGeneral) NewEditHandler(w http.ResponseWriter, r *http.Request) {
	var hash crypto.Hash
	if err := r.ParseForm(); err == nil {
		hash = crypto.DecodeHash(r.FormValue("draftHash"))
		fmt.Println(crypto.EncodeHash(hash))
		if view := NewEdit(a.state, hash); view != nil {
			if err := a.templates.ExecuteTemplate(w, "edit.html", view); err != nil {
				log.Println(err)
			}
			return
		}
	}
}

func (a *AttorneyGeneral) BoardsHandler(w http.ResponseWriter, r *http.Request) {
	view := BoardsFromState(a.state)
	if err := a.templates.ExecuteTemplate(w, "boards.html", view); err != nil {
		log.Println(err)
	}
}

func (a *AttorneyGeneral) BoardHandler(w http.ResponseWriter, r *http.Request) {
	author := a.Author(r)
	boardName := r.URL.Path
	boardName = strings.Replace(boardName, "/board/", "", 1)
	view := BoardDetailFromState(a.state, boardName, author)
	if view == nil {
		w.Write([]byte("board not found"))
	} else if err := a.templates.ExecuteTemplate(w, "board.html", view); err != nil {
		log.Println(err)
	}
}

func (a *AttorneyGeneral) CollectivesHandler(w http.ResponseWriter, r *http.Request) {
	view := ColletivesFromState(a.state)
	if err := a.templates.ExecuteTemplate(w, "collectives.html", view); err != nil {
		log.Println(err)
	}
}

func (a *AttorneyGeneral) CollectiveHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path
	name = strings.Replace(name, "/collective/", "", 1)
	author := a.Author(r)
	view := CollectiveDetailFromState(a.state, a.indexer, name, author)
	if err := a.templates.ExecuteTemplate(w, "collective.html", view); err != nil {
		log.Println(err)
	}
}

func (a *AttorneyGeneral) DraftsHandler(w http.ResponseWriter, r *http.Request) {
	view := DraftsFromState(a.state)
	if err := a.templates.ExecuteTemplate(w, "drafts.html", view); err != nil {
		log.Println(err)
	}
}

func (a *AttorneyGeneral) DraftHandler(w http.ResponseWriter, r *http.Request) {
	hashEncoded := r.URL.Path
	hashEncoded = strings.Replace(hashEncoded, "/draft/", "", 1)
	hash := crypto.DecodeHash(hashEncoded)
	author := a.Author(r)
	view := DraftDetailFromState(a.state, a.indexer, hash, author, a.genesisTime)
	if err := a.templates.ExecuteTemplate(w, "draft.html", view); err != nil {
		log.Println(err)
	}
}

func (a *AttorneyGeneral) EditsHandler(w http.ResponseWriter, r *http.Request) {
	hash := getHash(r.URL.Path, "/edits/")
	view := EditsFromState(a.state, hash)
	if err := a.templates.ExecuteTemplate(w, "edits.html", view); err != nil {
		log.Println(err)
	}
}

func (a *AttorneyGeneral) EventsHandler(w http.ResponseWriter, r *http.Request) {
	view := EventsFromState(a.state)
	if err := a.templates.ExecuteTemplate(w, "events.html", view); err != nil {
		log.Println(err)
	}
}

func (a *AttorneyGeneral) EventHandler(w http.ResponseWriter, r *http.Request) {
	hashEncoded := r.URL.Path
	hashEncoded = strings.Replace(hashEncoded, "/event/", "", 1)
	hash := crypto.DecodeHash(hashEncoded)
	author := a.Author(r)
	view := EventDetailFromState(a.state, a.indexer, hash, author, a.ephemeralprv)
	if err := a.templates.ExecuteTemplate(w, "event.html", view); err != nil {
		log.Println(err)
	}
}

func (a *AttorneyGeneral) RequestMemberShipVoteHandler(w http.ResponseWriter, r *http.Request) {
	hash := getHash(r.URL.Path, "/requestmembership/")
	view := RequestMembershipFromState(a.state, hash)
	if err := a.templates.ExecuteTemplate(w, "requestmembershipvote.html", view); err != nil {
		log.Println(err)
	}
}

func (a *AttorneyGeneral) VotesHandler(w http.ResponseWriter, r *http.Request) {
	author := a.Author(r)
	view := VotesFromState(a.state, a.indexer, author)
	if err := a.templates.ExecuteTemplate(w, "votes.html", view); err != nil {
		log.Println(err)
	}
}

func (a *AttorneyGeneral) MembersHandler(w http.ResponseWriter, r *http.Request) {
	view := MembersFromState(a.state)
	if err := a.templates.ExecuteTemplate(w, "members.html", view); err != nil {
		log.Println(err)
	}
}

func (a *AttorneyGeneral) MemberHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path
	name = strings.Replace(name, "/member/", "", 1)
	view := MemberViewFromState(a.state, a.indexer, name)
	if err := a.templates.ExecuteTemplate(w, "member.html", view); err != nil {
		log.Println(err)
	}
}

func (a *AttorneyGeneral) CreateBoardHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	collective := r.FormValue("collective")
	head := HeaderInfo{
		Active:  "Connections",
		Path:    "venture > connections > collectives > " + collective + " > ",
		EndPath: "create board",
		Section: "venture",
	}
	info := TemplateInfo{
		Head:           head,
		CollectiveName: collective,
	}
	if err := a.templates.ExecuteTemplate(w, "createboard.html", info); err != nil {
		log.Println(err)
	}
}

func (a *AttorneyGeneral) VoteCreateBoardHandler(w http.ResponseWriter, r *http.Request) {
	hash := getHash(r.URL.Path, "/votecreateboard/")
	view := PendingBoardFromState(a.state, hash)
	if err := a.templates.ExecuteTemplate(w, "votecreateboard.html", view); err != nil {
		log.Println(err)
	}
}

func (a *AttorneyGeneral) VoteCreateEventHandler(w http.ResponseWriter, r *http.Request) {
	hash := getHash(r.URL.Path, "/votecreateevent/")
	view := PendingEventFromState(a.state, a.indexer, hash)
	if err := a.templates.ExecuteTemplate(w, "votecreateevent.html", view); err != nil {
		log.Println(err)
	}
}

func (a *AttorneyGeneral) VoteCancelEventHandler(w http.ResponseWriter, r *http.Request) {
	hash := getHash(r.URL.Path, "/votecancelevent/")
	view := CancelEventFromState(a.state, a.indexer, hash)
	if err := a.templates.ExecuteTemplate(w, "votecancelevent.html", view); err != nil {
		log.Println(err)
	}
}

func (a *AttorneyGeneral) UpdateCollectiveHandler(w http.ResponseWriter, r *http.Request) {
	collective := strings.Replace(r.URL.Path, "/updatecollective/", "", 1)
	view := CollectiveToUpdateFromState(a.state, collective)
	if err := a.templates.ExecuteTemplate(w, "updatecollective.html", view); err != nil {
		log.Println(err)
	}
}

func (a *AttorneyGeneral) VoteUpdateCollectiveHandler(w http.ResponseWriter, r *http.Request) {
	author := a.Author(r)
	hash := getHash(r.URL.Path, "/voteupdatecollective/")
	view := CollectiveUpdateFromState(a.state, hash, author)
	if err := a.templates.ExecuteTemplate(w, "voteupdatecollective.html", view); err != nil {
		log.Println(err)
	}
}

func (a *AttorneyGeneral) UpdateBoardHandler(w http.ResponseWriter, r *http.Request) {
	board := strings.Replace(r.URL.Path, "/updateboard/", "", 1)
	view := BoardToUpdateFromState(a.state, board)
	if err := a.templates.ExecuteTemplate(w, "updateboard.html", view); err != nil {
		log.Println(err)
	}
}

func (a *AttorneyGeneral) VoteUpdateBoardHandler(w http.ResponseWriter, r *http.Request) {
	hash := getHash(r.URL.Path, "/votecreateboard/")
	view := BoardUpdateFromState(a.state, hash)
	if err := a.templates.ExecuteTemplate(w, "voteupdateboard.html", view); err != nil {
		log.Println(err)
	}
}

func (a *AttorneyGeneral) UpdateEventHandler(w http.ResponseWriter, r *http.Request) {
	author := a.Author(r)
	hash := getHash(r.URL.Path, "/updateevent/")
	view := EventUpdateDetailFromState(a.state, a.indexer, hash, author)
	if err := a.templates.ExecuteTemplate(w, "updateevent.html", view); err != nil {
		log.Println(err)
	}
}

func (a *AttorneyGeneral) CreateEventHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	collective := r.FormValue("collective")
	head := HeaderInfo{
		Active:  "Connections",
		Path:    "venture > connections > collectives > " + collective + " > ",
		EndPath: "create event",
		Section: "venture",
	}
	info := TemplateInfo{
		Head:           head,
		CollectiveName: collective,
	}
	if err := a.templates.ExecuteTemplate(w, "createevent.html", info); err != nil {
		log.Println(err)
	}
}

func (a *AttorneyGeneral) VoteUpdateEventHandler(w http.ResponseWriter, r *http.Request) {
	author := a.Author(r)
	hash := getHash(r.URL.Path, "/voteupdateevent/")
	view := EventUpdateFromState(a.state, hash, author)
	if err := a.templates.ExecuteTemplate(w, "voteupdateevent.html", view); err != nil {
		log.Println(err)
	}
}

func (a *AttorneyGeneral) ConnectionsHandler(w http.ResponseWriter, r *http.Request) {
	author := a.Author(r)
	view := ConnectionsFromState(a.state, a.indexer, author, a.genesisTime)
	if err := a.templates.ExecuteTemplate(w, "connections.html", view); err != nil {
		log.Println(err)
	}
}

func (a *AttorneyGeneral) UpdatesHandler(w http.ResponseWriter, r *http.Request) {
	author := a.Author(r)
	view := UpdatesViewFromState(a.state, a.indexer, author, a.genesisTime)
	if err := a.templates.ExecuteTemplate(w, "updates.html", view); err != nil {
		log.Println(err)
	}
}

func (a *AttorneyGeneral) PendingActionsHandler(w http.ResponseWriter, r *http.Request) {
	author := a.Author(r)
	view := PendingActionsFromState(a.state, a.indexer, author, a.genesisTime)
	if err := a.templates.ExecuteTemplate(w, "pending.html", view); err != nil {
		log.Println(err)
	}
}

func (a *AttorneyGeneral) MyMediaHandler(w http.ResponseWriter, r *http.Request) {
	author := a.Author(r)
	view := MyMediaFromState(a.state, a.indexer, author)
	if err := a.templates.ExecuteTemplate(w, "mymedia.html", view); err != nil {
		log.Println(err)
	}
}

func (a *AttorneyGeneral) MyEventsHandler(w http.ResponseWriter, r *http.Request) {
	author := a.Author(r)
	view := MyEventsFromState(a.state, a.indexer, author)
	if err := a.templates.ExecuteTemplate(w, "myevents.html", view); err != nil {
		log.Println(err)
	}
}

func (a *AttorneyGeneral) NewsHandler(w http.ResponseWriter, r *http.Request) {
	view := NewActionsFromState(a.state, a.indexer, a.genesisTime)
	if err := a.templates.ExecuteTemplate(w, "news.html", view); err != nil {
		log.Println(err)
	}
}

func (a *AttorneyGeneral) DetailedVoteHandler(w http.ResponseWriter, r *http.Request) {
	hash := getHash(r.URL.Path, "/detailedvote/")
	view := DetailedVoteFromState(a.state, a.indexer, hash, a.genesisTime)
	if err := a.templates.ExecuteTemplate(w, "detailedvote.html", view); err != nil {
		log.Println(err)
	}
}