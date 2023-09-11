package state

import (
	"errors"
	"time"

	"github.com/lienkolabs/swell/crypto"
	"github.com/lienkolabs/synergy/social/actions"
)

type Greeting struct {
	Action       *actions.GreetCheckinEvent
	EphemeralKey crypto.Token
}

type Event struct {
	Collective   *Collective
	StartAt      time.Time
	EstimatedEnd time.Time
	Description  string
	Venue        string
	Open         bool
	Public       bool
	Hash         crypto.Hash
	Managers     *UnamedCollective // default é qualquer um do coletivo
	Votes        []actions.Vote
	Checkin      map[crypto.Token]*Greeting
	Live         bool
}

func (p *Event) IncorporateVote(vote actions.Vote, state *State) error {
	if err := IsNewValidVote(vote, p.Votes, p.Hash); err != nil {
		return err
	}
	p.Votes = append(p.Votes, vote)
	if p.Live {
		return nil
	}
	if !p.Collective.Consensus(p.Hash, p.Votes) {
		return nil
	}
	// new consensus
	state.IndexConsensus(vote.Hash, true)
	p.Live = true
	//if state.index != nil {
	//	state.index.AddEventToCollective(p, p.Collective)
	//}
	state.Proposals.Delete(p.Hash)
	if _, ok := state.Events[p.Hash]; !ok {
		state.Events[p.Hash] = p
		return nil
	}
	return errors.New("already live")
}

type EventUpdate struct {
	Event           *Event
	StartAt         *time.Time
	EstimatedEnd    *time.Time
	Description     *string
	Venue           *string
	Open            *bool
	Public          *bool
	ManagerMajority *byte
	Votes           []actions.Vote
	Hash            crypto.Hash
	Updated         bool
	Reasons         string
}

func (p *EventUpdate) IncorporateVote(vote actions.Vote, state *State) error {
	if err := IsNewValidVote(vote, p.Votes, p.Hash); err != nil {
		return err
	}
	p.Votes = append(p.Votes, vote)
	if p.Updated {
		return nil
	}
	if !p.Event.Managers.Consensus(p.Hash, p.Votes) {
		return nil
	}
	// new consensus, update event details
	state.IndexConsensus(vote.Hash, true)
	p.Updated = true
	state.Proposals.Delete(p.Hash)
	if event := p.Event; event != nil {
		if p.StartAt != nil {
			event.StartAt = *p.StartAt
		}
		if p.EstimatedEnd != nil {
			event.EstimatedEnd = *p.EstimatedEnd
		}
		if p.Description != nil {
			event.Description = *p.Description
		}
		if p.Venue != nil {
			event.Venue = *p.Venue
		}
		if p.Open != nil {
			event.Open = *p.Open
		}
		if p.Public != nil {
			event.Public = *p.Public
		}
		if p.ManagerMajority != nil {
			p.Event.Managers.Majority = int(*p.ManagerMajority)
		}
		return nil
	}
	return errors.New("event not found")
}

type CancelEvent struct {
	Event *Event
	Hash  crypto.Hash
	Votes []actions.Vote
}

func (p *CancelEvent) IncorporateVote(vote actions.Vote, state *State) error {
	if err := IsNewValidVote(vote, p.Votes, p.Hash); err != nil {
		return err
	}
	p.Votes = append(p.Votes, vote)
	if !p.Event.Live {
		return nil
	}
	if !p.Event.Managers.Consensus(p.Hash, p.Votes) {
		return nil
	}
	// new consensus, update event details
	state.IndexConsensus(vote.Hash, true)
	p.Event.Live = false
	//if state.index != nil {
	//	state.index.RemoveEventFromCollective(p.Event, p.Event.Collective)
	//}
	state.Proposals.Delete(p.Hash)
	return nil
}

type EventCheckinGreet struct {
	Event  *Event
	Hash   crypto.Hash
	Greets []actions.GreetCheckinEvent
}

func (p *EventCheckinGreet) IncorporateGreet(greet actions.GreetCheckinEvent, state *State) error {

	if greet.EventHash != p.Hash {
		return errors.New("invalid hash")
	}
	for _, cast := range p.Greets {
		if cast.CheckedIn == greet.CheckedIn {
			return errors.New("checkin already greeted")
		}
	}
	p.Greets = append(p.Greets, greet)

	if !p.Event.Live {
		return nil
	}

	state.Proposals.Delete(p.Hash)
	return nil
}
