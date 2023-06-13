package state

import "github.com/lienkolabs/synergy/social/actions"

const (
	DraftProposal byte = iota
)

type Proposal interface {
	IncorporateVote(vote actions.VoteAction, state *State) error
}
