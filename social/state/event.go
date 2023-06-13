package state

import (
	"errors"
	"time"

	"github.com/lienkolabs/swell/crypto"
	"github.com/lienkolabs/synergy/social/actions"
)

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
	p.Live = true
	delete(state.Proposals, p.Hash)
	if _, ok := state.Events[p.Hash]; !ok {
		state.Events[p.Hash] = p
		return nil
	}
	return errors.New("already live")
}

type EventUpdate struct {
	Collective   *Collective
	Event        *Event
	StartAt      *time.Time
	EstimatedEnd *time.Time
	Description  *string
	Venue        *string
	Open         *bool
	Public       *bool
	Votes        []actions.Vote
	Hash         crypto.Hash
	Updated      bool
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
	p.Updated = true
	delete(state.Proposals, p.Hash)
	if event, ok := state.Events[p.Hash]; ok {
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
		return nil
	}
	return errors.New("event not found")
}
