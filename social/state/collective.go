package state

import (
	"errors"

	"github.com/lienkolabs/swell/crypto"
	"github.com/lienkolabs/synergy/social/actions"
)

type Collective struct {
	Name        string
	Members     map[crypto.Token]struct{}
	Description string
	Policy      actions.Policy
}

func CreateCollectiveToState(create *actions.CreateCollective, s *State) error {
	if _, ok := s.Members[create.Author]; !ok {
		return errors.New("not a member of synergy")
	}
	if _, ok := s.Collectives[create.Name]; ok {
		return errors.New("collective already exists")
	}
	if create.Policy.Majority < 0 || create.Policy.Majority > 100 || create.Policy.SuperMajority < 0 || create.Policy.SuperMajority > 100 {
		return errors.New("invalid policy")
	}
	s.Collectives[create.Name] = &Collective{
		Name:        create.Name,
		Members:     map[crypto.Token]struct{}{},
		Description: create.Description,
		Policy: actions.Policy{
			Majority:      create.Policy.Majority,
			SuperMajority: create.Policy.SuperMajority,
		},
	}
	return nil
}

func (c *Collective) ListOfMembers() map[crypto.Token]struct{} {
	return nil
}

func (c *Collective) Photo() *Collective {
	cloned := Collective{
		Name:    c.Name,
		Members: make(map[crypto.Token]struct{}),
		Policy: actions.Policy{
			Majority:      c.Policy.Majority,
			SuperMajority: c.Policy.SuperMajority,
		},
	}
	for member, _ := range c.Members {
		cloned.Members[member] = struct{}{}
	}
	return &cloned
}

func (c *Collective) IncludeMember(token crypto.Token) {
	c.Members[token] = struct{}{}
}

func (c *Collective) RemoveMember(token crypto.Token) {
	delete(c.Members, token)
}

func (c *Collective) ChangeMajority(majority int) {
	c.Policy.Majority = majority
}

func (c *Collective) Consensus(hash crypto.Hash, votes []actions.Vote) bool {
	required := len(c.Members)*c.Policy.Majority/100 + 1
	if required > len(c.Members) {
		required = len(c.Members)
	}
	return consensus(c.Members, required, hash, votes)
}

func (c *Collective) SuperConsensus(hash crypto.Hash, votes []actions.Vote) bool {
	required := len(c.Members)*c.Policy.SuperMajority/100 + 1
	if required > len(c.Members) {
		required = len(c.Members)
	}
	return consensus(c.Members, required, hash, votes)
}

func (c *Collective) IsMember(token crypto.Token) bool {
	_, ok := c.Members[token]
	return ok
}

type UnamedCollective struct {
	Members  map[crypto.Token]struct{}
	Majority int
}

func (c *UnamedCollective) ListOfMembers() map[crypto.Token]struct{} {
	return c.Members
}

func (c *UnamedCollective) Consensus(hash crypto.Hash, votes []actions.Vote) bool {
	required := len(c.Members)*c.Majority/100 + 1
	if required > len(c.Members) {
		required = len(c.Members)
	}
	return consensus(c.Members, required, hash, votes)
}

func (c *UnamedCollective) IsMember(token crypto.Token) bool {
	_, ok := c.Members[token]
	return ok
}

func (c *UnamedCollective) IncludeMember(token crypto.Token) {
	c.Members[token] = struct{}{}
}

func (c *UnamedCollective) RemoveMember(token crypto.Token) {
	delete(c.Members, token)
}

func (c *UnamedCollective) ChangeMajority(majority int) {
	c.Majority = majority
}

type PendingUpdate struct {
	Update       *actions.UpdateCollective
	Collective   *Collective
	Hash         crypto.Hash
	ChangePolicy bool
	Votes        []actions.Vote
}

func UpdateCollectiveToState(update *actions.UpdateCollective, s *State) error {
	collective, ok := s.Collectives[update.OnBehalfOf]
	if !ok {
		return errors.New("unkown collective")
	}
	if !collective.IsMember(update.Author) {
		return errors.New("not a member of collective")
	}
	hash := crypto.Hasher(update.Serialize()) // proposal hash = hash of instruction
	vote := actions.Vote{
		Epoch:   update.Epoch,
		Author:  update.Author,
		Reasons: "commit",
		Hash:    hash,
		Approve: true,
	}

	if update.Policy != nil {
		if update.Policy.Majority < 0 || update.Policy.Majority > 100 || update.Policy.SuperMajority < 0 || update.Policy.SuperMajority > 100 {
			return errors.New("invalid policy")
		}
		if collective.SuperConsensus(hash, []actions.Vote{vote}) {
			if update.Description != "" {
				collective.Description = update.Description
			}
			collective.Policy = actions.Policy{
				Majority:      update.Policy.Majority,
				SuperMajority: update.Policy.SuperMajority,
			}
			return nil
		}
	} else {
		if collective.Consensus(hash, []actions.Vote{vote}) {
			if update.Description != "" {
				collective.Description = update.Description
			}
			return nil
		}
	}

	pending := PendingUpdate{
		Update: update,
		// consensus is based on the collective composition at the moment
		// of incorporation of instruction
		Collective: collective.Photo(),
		Hash:       hash,
		Votes:      []actions.Vote{vote},
	}
	if update.Policy != nil {
		pending.ChangePolicy = true
	}
	s.Proposals[hash] = &pending
	s.setDeadline(update.Epoch+ProposalDeadline, hash)
	return nil
}

func (p *PendingUpdate) IncorporateVote(vote actions.Vote, state *State) error {
	if err := isValidVote(p.Hash, vote, p.Votes); err != nil {
		return err
	}
	p.Votes = append(p.Votes, vote)
	if p.ChangePolicy {
		if !p.Collective.SuperConsensus(p.Hash, p.Votes) {
			return nil
		}
	} else {
		if !p.Collective.Consensus(p.Hash, p.Votes) {
			return nil
		}
	}
	// exclude pending update from live proposals because of consensus
	delete(state.Proposals, p.Hash)
	// update collective
	if p.Update.Description != "" {
		p.Collective.Description = p.Update.Description
	}
	if p.ChangePolicy {
		p.Collective.Policy = actions.Policy{
			Majority:      p.Update.Policy.Majority,
			SuperMajority: p.Update.Policy.SuperMajority,
		}
	}
	return nil
}

type PendingRequestMembership struct {
	Request    *actions.RequestMembership
	Collective *Collective
	Hash       crypto.Hash
	Votes      []actions.Vote
}

func (p *PendingRequestMembership) IncorporateVote(vote actions.Vote, state *State) error {
	if err := isValidVote(p.Hash, vote, p.Votes); err != nil {
		return err
	}
	p.Votes = append(p.Votes, vote)
	if !p.Collective.Consensus(p.Hash, p.Votes) {
		return nil
	}
	delete(state.Proposals, p.Hash)
	collective, ok := state.Collectives[p.Collective.Name]
	if !ok {
		return errors.New("collective not found")
	}
	collective.Members[p.Request.Author] = struct{}{}
	return nil
}

type PendingRemoveMember struct {
	Remove     *actions.RemoveMember
	Collective *Collective
	Hash       crypto.Hash
	Votes      []actions.Vote
}

func (p *PendingRemoveMember) IncorporateVote(vote actions.Vote, state *State) error {
	if err := isValidVote(p.Hash, vote, p.Votes); err != nil {
		return err
	}
	p.Votes = append(p.Votes, vote)
	if !p.Collective.Consensus(p.Hash, p.Votes) {
		return nil
	}
	delete(state.Proposals, p.Hash)
	collective, ok := state.Collectives[p.Collective.Name]
	if !ok {
		return errors.New("collective not found")
	}
	delete(collective.Members, p.Remove.Member)
	return nil
}
