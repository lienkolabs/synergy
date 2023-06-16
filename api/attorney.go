package api

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/lienkolabs/swell/crypto"
	"github.com/lienkolabs/swell/util"
	"github.com/lienkolabs/synergy/social"
	"github.com/lienkolabs/synergy/social/actions"
	"github.com/lienkolabs/synergy/social/state"
)

var templateFiles []string = []string{
	"members", "collectives", "collective",
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
		mux.HandleFunc("/members", attorney.membersHandler)
		mux.HandleFunc("/collectives", attorney.collectivesHandler)

		err := http.ListenAndServe(fmt.Sprintf(":%v", port), mux)
		if err != nil {
			log.Fatal(err)
		}
	}()

	return nil
}

// endpoint "/collectives" vai ser respondido por esta função
func (a *Attorney) collectivesHandler(w http.ResponseWriter, r *http.Request) {
	t := a.templates["colletives"]
	view := ColletivesFromState(a.gateway.State)
	t.Execute(w, view)
}

func (a *Attorney) membersHandler(w http.ResponseWriter, r *http.Request) {
	t := a.templates["members"]
	view := MembersFromState(a.gateway.State)
	t.Execute(w, view)
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
