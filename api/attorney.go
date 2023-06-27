package api

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/lienkolabs/swell/crypto"
	"github.com/lienkolabs/swell/util"
	"github.com/lienkolabs/synergy/social"
	"github.com/lienkolabs/synergy/social/actions"
	"github.com/lienkolabs/synergy/social/state"
)

var templateFiles []string = []string{
	"boards", "board", "collectives", "collective", "draft", "drafts", "edits", "events",
	"event", "member", "members", "votes", "requestmembershipvote", "newdraft2", "edit",
	"createboard", "votecreateboard", "updateboard", "voteupdateboard", "updateevent",
	"updatecollective", "voteupdatecollective", "createevent", "voteupdateevent",
}

type Attorney struct {
	author    crypto.Token
	pk        crypto.PrivateKey
	wallet    crypto.PrivateKey
	pending   map[crypto.Hash]actions.Action
	epoch     uint64
	gateway   *social.Gateway
	state     *state.State
	templates map[string]*template.Template
}

func NewAttorneyServer(pk crypto.PrivateKey, token crypto.Token, port int, gateway *social.Gateway) *Attorney {
	attorney := Attorney{
		author:  token,
		pk:      pk,
		wallet:  pk,
		pending: make(map[crypto.Hash]actions.Action),
		gateway: gateway,
		state:   gateway.State,
		epoch:   0,
	}
	blockEvent := gateway.Register()
	send := make(chan actions.Action)
	go func() {
		for {
			select {
			case epoch := <-blockEvent:
				attorney.epoch = epoch
			case action := <-send:
				gateway.Action(attorney.DressAction(action))
			}
		}
	}()

	go func() {
		attorney.templates = make(map[string]*template.Template)
		for _, file := range templateFiles {
			filePath := fmt.Sprintf("./api/templates/%v.html", file)
			if t, err := template.ParseFiles(filePath); err != nil {
				log.Fatal(err)
			} else {
				attorney.templates[file] = t
			}
		}

		mux := http.NewServeMux()

		fs := http.FileServer(http.Dir("./api/static"))
		mux.Handle("/static/", http.StripPrefix("/static/", fs))

		mux.HandleFunc("/api", attorney.ApiHandler)
		mux.HandleFunc("/boards", attorney.BoardsHandler)
		mux.HandleFunc("/board/", attorney.BoardHandler)
		mux.HandleFunc("/collectives", attorney.CollectivesHandler)
		mux.HandleFunc("/collective/", attorney.CollectiveHandler)
		mux.HandleFunc("/drafts", attorney.DraftsHandler)
		mux.HandleFunc("/draft/", attorney.DraftHandler)
		mux.HandleFunc("/edits", attorney.EditsHandler)
		mux.HandleFunc("/events", attorney.EventsHandler)
		mux.HandleFunc("/event/", attorney.EventHandler)
		mux.HandleFunc("/members", attorney.MembersHandler)
		mux.HandleFunc("/member/", attorney.MemberHandler)
		mux.HandleFunc("/votes/", attorney.VotesHandler)
		mux.HandleFunc("/requestmembership/", attorney.RequestMemberShipVoteHandler)
		mux.HandleFunc("/newdraft", attorney.NewDraft2Handler)
		mux.HandleFunc("/edit", attorney.NewEditHandler)
		mux.HandleFunc("/media/", attorney.MediaHandler)
		mux.HandleFunc("/uploadfile", attorney.UploadHandler)
		mux.HandleFunc("/createboard", attorney.CreateBoardHandler)
		mux.HandleFunc("/votecreateboard/", attorney.VoteCreateBoardHandler)
		mux.HandleFunc("/updateboard/", attorney.UpdateBoardHandler)
		mux.HandleFunc("/voteupdateboard/", attorney.UpdateBoardHandler)
		mux.HandleFunc("/updatecollective/", attorney.UpdateCollectiveHandler)
		mux.HandleFunc("/voteupdatecollective/", attorney.VoteUpdateCollectiveHandler)
		mux.HandleFunc("/updateevent/", attorney.UpdateEventHandler)
		mux.HandleFunc("/createevent", attorney.CreateEventHandler)
		mux.HandleFunc("/voteupdateevent/", attorney.VoteUpdateEventHandler)
		// mux.HandleFunc("/member/votes", attorney.VotesHandler)

		srv := &http.Server{
			Addr:         fmt.Sprintf(":%v", port),
			Handler:      mux,
			WriteTimeout: 2 * time.Second,
		}
		srv.ListenAndServe()

	}()

	return nil
}

func (a *Attorney) Send(all []actions.Action) {
	for _, action := range all {
		dressed := a.DressAction(action)
		a.gateway.Action(dressed)
	}
}

// Dress a giving action with current epoch, attorney´s author
// attorneys signature, attorneys wallet and wallet signature
func (a *Attorney) DressAction(action actions.Action) []byte {
	bytes := action.Serialize()
	dress := []byte{0}
	util.PutUint64(a.epoch, &dress)
	util.PutToken(a.author, &dress)
	dress = append(dress, 0, 1, 0, 0) // axe void synergy
	dress = append(dress, bytes[8+crypto.TokenSize:]...)
	util.PutToken(a.pk.PublicKey(), &dress)
	signature := a.pk.Sign(dress)
	util.PutSignature(signature, &dress)
	util.PutToken(a.wallet.PublicKey(), &dress)
	util.PutUint64(0, &dress) // zero fee
	signature = a.wallet.Sign(dress)
	util.PutSignature(signature, &dress)
	return dress
}

func (a *Attorney) Confirmed(hash crypto.Hash) {
	delete(a.pending, hash)
}
