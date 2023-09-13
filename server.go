package main

import (
	"fmt"

	"github.com/lienkolabs/breeze/crypto"
	"github.com/lienkolabs/synergy/api"
	"github.com/lienkolabs/synergy/social"
	"github.com/lienkolabs/synergy/social/index"
)

var pks []crypto.PrivateKey = []crypto.PrivateKey{
	{118, 35, 197, 163, 215, 20, 35, 190, 110, 151, 246, 231, 86, 177, 156, 89, 122, 69, 28, 233, 185, 150, 126, 169, 237, 173, 83, 120, 145, 238, 242, 137,
		171, 216, 111, 131, 116, 217, 38, 148, 28, 178, 174, 63, 166, 4, 50, 6, 20, 133, 15, 153, 41, 252, 164, 165, 2, 127, 163, 204, 24, 24, 188, 240},
	{152, 224, 227, 154, 131, 1, 186, 147, 73, 37, 4, 253, 11, 148, 195, 67, 86, 85, 28, 162, 78, 239, 168, 42, 204, 222, 144, 41, 186, 246, 250, 57, 125, 202,
		107, 133, 63, 39, 136, 246, 120, 222, 29, 73, 106, 213, 95, 132, 50, 130, 162, 42, 95, 159, 10, 246, 213, 217, 160, 125, 181, 194, 37, 174},
	{125, 86, 238, 128, 237, 4, 143, 47, 214, 72, 71, 47, 72, 45, 214, 45, 178, 98, 105, 154, 171, 151, 73, 183, 234, 120, 128, 38, 174, 253, 105, 162, 189,
		253, 40, 134, 214, 5, 229, 224, 171, 175, 152, 114, 72, 167, 9, 215, 75, 171, 3, 255, 30, 255, 110, 127, 9, 3, 129, 24, 230, 246, 109, 184},
}

func main() {
	N := 3
	users := make(map[crypto.Token]string)
	userToken := make([]crypto.Token, N)
	indexer := index.NewIndex()
	for n := 0; n < N; n++ {
		userToken[n] = pks[n].PublicKey()
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
