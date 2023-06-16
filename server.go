package main

import (
	"fmt"
	"net/http"

	"github.com/lienkolabs/swell/crypto"
	"github.com/lienkolabs/synergy/api"
	"github.com/lienkolabs/synergy/social"
)

func main() {
	user, _ := crypto.RandomAsymetricKey()
	fmt.Printf("%v", user)
	users := map[crypto.Token]string{user: "Ruben"}
	state := social.TestGenesisState(users)

	gateway := social.SelfGateway(state) // simulador de blockchain

	fs := http.FileServer(http.Dir("./api/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	_, attorneySecret := crypto.RandomAsymetricKey()
	api.NewAttorneyServer(attorneySecret, user, 3000, gateway)
	for true {

	}
}
