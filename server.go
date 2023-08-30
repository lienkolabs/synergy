package main

import (
	"fmt"

	"github.com/lienkolabs/swell/crypto"
	"github.com/lienkolabs/synergy/api"
	"github.com/lienkolabs/synergy/social"
	"github.com/lienkolabs/synergy/social/index"
)

func main() {
	N := 3
	users := make(map[crypto.Token]string)
	userToken := make([]crypto.Token, N)
	for n := 0; n < N; n++ {
		userToken[n], _ = crypto.RandomAsymetricKey()
		users[userToken[n]] = fmt.Sprintf("user_%v", n)
	}
	indexer := index.NewIndex()
	state := social.TestGenesisState(users, indexer)
	gateway := social.SelfGateway(state) // simulador de blockchain

	_, attorneySecret := crypto.RandomAsymetricKey()
	for n := 0; n < N; n++ {
		api.NewAttorneyServer(attorneySecret, userToken[n], 3000+n, gateway, indexer)
	}

	for true {

	}
}
