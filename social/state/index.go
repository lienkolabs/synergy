package state

import "github.com/lienkolabs/swell/crypto"

type Indexer interface {
	AddBoardToCollective(*Board, *Collective)
	RemoveBoardFromCollective(*Board, *Collective)

	AddStampToCollective(*Stamp, *Collective)

	AddEventToCollective(*Event, *Collective)
	RemoveEventFromCollective(*Event, *Collective)

	IndexConsensus(crypto.Hash, bool)
}
