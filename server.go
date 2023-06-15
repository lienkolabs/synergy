package main

import (
	"log"
	"net/http"
	"time"

	"github.com/lienkolabs/swell/crypto"
	"github.com/lienkolabs/synergy/api"
	"github.com/lienkolabs/synergy/social"
	"github.com/lienkolabs/synergy/social/actions"
)

func main4() {
	user, _ := crypto.RandomAsymetricKey()
	users := map[crypto.Token]string{user: "Ruben"}
	state := social.TestGenesisState(users)
	gateway := social.SelfGateway(state)

	<-time.After(time.Second)
	draft := &actions.Draft{
		Epoch:         18,
		Author:        user,
		Reasons:       "draft test",
		OnBehalfOf:    "",
		CoAuthors:     []crypto.Token{},
		Title:         "first_draft",
		Keywords:      []string{"Testando"},
		Description:   "draft test",
		ContentType:   "txt",
		ContentHash:   crypto.Hasher([]byte{0, 1, 2, 3, 4, 5}),
		NumberOfParts: 1,
		Content:       []byte{0, 1, 2, 3, 4, 5},
		PreviousDraft: crypto.ZeroHash,
		References:    []crypto.Hash{},
	}
	gateway.Action(draft.Serialize())
	<-time.After(5 * time.Second)
	gateway.Stop()

}

func main() {
	fs := http.FileServer(http.Dir("./api/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	token, _ := crypto.RandomAsymetricKey()
	handles := &api.Handles{
		Handle: map[string]crypto.Token{"@rubis": token},
	}

	http.HandleFunc("/api", handles.ApiHandler)

	log.Print("Listening on :3000...")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
