package state

type Indexer interface {
	AddBoardToCollective(*Board, *Collective)
	RemoveBoardFromCollective(*Board, *Collective)

	AddStampToCollective(*Stamp, *Collective)

	AddEventToCollective(*Event, *Collective)
	RemoveEventFromCollective(*Event, *Collective)
}
