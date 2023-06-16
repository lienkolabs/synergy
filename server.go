package main

import (
	"log"
	"net/http"

	"github.com/lienkolabs/swell/crypto"
	"github.com/lienkolabs/synergy/api"
	"github.com/lienkolabs/synergy/social"
)

func main() {
	user, _ := crypto.RandomAsymetricKey()
	users := map[crypto.Token]string{user: "Ruben"}
	state := social.TestGenesisState(users)

	gateway := social.SelfGateway(state) // simulador de blockchain

	fs := http.FileServer(http.Dir("./api/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	_, attorneySecret := crypto.RandomAsymetricKey()
	attorneyProxy := api.NewAttorneyServer(attorneySecret, user, 3200, gateway)
	stateView := api.NewStateView(state)

	handles := &api.Handles{
		Handle:   map[string]crypto.Token{"Ruben": user},
		Attorney: attorneyProxy,
	}
	http.HandleFunc("/api", handles.ApiHandler)

	http.HandleFunc("/members", stateView.MembersHandler)
	http.HandleFunc("/collectives", stateView.CollectivesHandler)

	log.Print("Listening on :3000...")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
