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
	indexer := index.NewIndex()
	for n := 0; n < N; n++ {
		userToken[n], _ = crypto.RandomAsymetricKey()
		users[userToken[n]] = fmt.Sprintf("user_%v", n)
		indexer.AddMemberToIndex(userToken[n], users[userToken[n]])
	}
	state := social.TestGenesisState(users, indexer)
	indexer.SetState(state)
	gateway := social.SelfGateway(state) // simulador de blockchain

	_, attorneySecret := crypto.RandomAsymetricKey()
	for n := 0; n < N; n++ {
		api.NewAttorneyServer(attorneySecret, userToken[n], 3000+n, gateway, indexer)
	}

	for true {

	}
}
