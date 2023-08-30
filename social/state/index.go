package state

type Indexer interface {
	AddEventToCollective(*Event, *Collective)
	RemoveEventFromCollective(*Event, *Collective)
}
