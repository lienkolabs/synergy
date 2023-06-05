package synergy

import "github.com/lienkolabs/swell/crypto"

type ReactionInstruction struct {
	Epoch    uint64
	Author   crypto.Token
	Hash     crypto.Hash
	Reaction byte
}
