package synergy

const (
	DraftProposal byte = iota
)

type Proposal interface {
	IncorporateVote(vote VoteInstruction, state *State) error
}
