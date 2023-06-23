package api

import (
	"fmt"
	"strings"

	"github.com/lienkolabs/swell/crypto"
	"github.com/lienkolabs/synergy/social/state"
)

// Drafts template struct

type DraftsView struct {
	Title       string
	Authors     []string
	Hash        string
	Description string
	Keywords    []string
}

type DraftsListView struct {
	Drafts []DraftsView
}

type AuthorDetail struct {
	Name       string
	Collective bool
}

type ReferenceDetail struct {
	Title  string
	Author string
	Date   string
}

// co-autor, stamp, pin, version, release
type DraftVoteAction struct {
	Kind       string
	OnBehalfOf string // collective or board editor
	Hash       string
}

type DraftDetailView struct {
	Title       string
	Date        string
	Description string
	Keywords    []string
	Hash        string
	//Content      string
	Authors      []AuthorDetail
	References   []ReferenceDetail
	PreviousHash string
	Pinned       []string
	Edited       bool
	Released     bool
	Stamps       []string
	Votes        []DraftVoteAction
	Policy       Policy
	Authorship   bool
}

func membersToHandles(members map[crypto.Token]struct{}, state *state.State) []string {
	handles := make([]string, 0)
	for member, _ := range members {
		handle, ok := state.Members[crypto.Hasher(member[:])]
		if ok {
			handles = append(handles, handle)
		}
	}
	return handles
}

func hashesToString(hashes []crypto.Hash) []string {
	output := make([]string, 0)
	for _, hash := range hashes {
		text, err := hash.MarshalText()
		if err != nil {
			output = append(output, string(text))
		}
	}
	return output
}

func PinList(pin []*state.Board) []string {
	list := make([]string, 0)
	if len(pin) == 0 {
		return list
	}
	for _, p := range pin {
		list = append(list, p.Name)
	}
	return list
}

func StampList(stamps []*state.Stamp) []string {
	list := make([]string, 0)
	if len(stamps) == 0 {
		return list
	}
	for _, p := range stamps {
		list = append(list, p.Reputation.Name)
	}
	return list
}

func References(r []crypto.Hash, s *state.State) []ReferenceDetail {
	references := make([]ReferenceDetail, 0)
	for _, hash := range r {
		if draft, ok := s.Drafts[hash]; ok {
			reference := ReferenceDetail{
				Title:  draft.Title,
				Author: authorsEtAll(draft.Authors, s),
				Date:   fmt.Sprintf("%v", draft.Date.Year()),
			}
			references = append(references, reference)
		}
	}
	return references
}

func authorsEtAll(c state.Consensual, s *state.State) string {
	authors := AuthorList(c, s)
	if len(authors) == 0 {
		return ""
	}
	N := len(authors)
	tail := ""
	if len(authors) > 3 {
		N = 3
		tail = " et al."
	}
	authorlist := make([]string, N)
	for n := 0; n < N; n++ {
		authorlist[n] = authors[n].Name
	}

	return fmt.Sprintf("%v%v", strings.Join(authorlist, ","), tail)
}

func AuthorList(c state.Consensual, s *state.State) []AuthorDetail {
	if c == nil {
		return []AuthorDetail{}
	}
	name := c.CollectiveName()
	if name != "" {
		author := AuthorDetail{
			Name:       name,
			Collective: true,
		}
		return []AuthorDetail{author}
	}
	authors := make([]AuthorDetail, 0)
	for token, _ := range c.ListOfMembers() {
		if handle, ok := s.Members[crypto.Hasher(token[:])]; ok {
			authors = append(authors, AuthorDetail{Name: handle})
		}
	}
	return authors
}

func DraftsFromState(state *state.State) DraftsListView {
	view := DraftsListView{
		Drafts: make([]DraftsView, 0),
	}
	for _, draft := range state.Drafts {
		hash, _ := draft.DraftHash.MarshalText()
		itemView := DraftsView{
			Title:       draft.Title,
			Hash:        string(hash),
			Authors:     membersToHandles(draft.Authors.ListOfMembers(), state),
			Description: draft.Description,
			Keywords:    draft.Keywords,
		}
		view.Drafts = append(view.Drafts, itemView)
	}
	return view
}

func DraftDetailFromState(s *state.State, hash crypto.Hash, token crypto.Token) *DraftDetailView {
	draft, ok := s.Drafts[hash]
	if !ok {
		draft, ok = s.Proposals.Draft[hash]
		if !ok {
			return nil
		}
	}
	hashText, _ := hash.MarshalText()
	view := DraftDetailView{
		Title:       draft.Title,
		Description: draft.Description,
		Keywords:    draft.Keywords,
		Authors:     AuthorList(draft.Authors, s),
		References:  References(draft.References, s),
		Pinned:      PinList(draft.Pinned),
		Votes:       make([]DraftVoteAction, 0),
		Authorship:  draft.Authors.IsMember(token),
		Hash:        string(hashText),
	}
	pending := s.Proposals.GetVotes(token)
	if len(pending) > 0 {
		for pendingHash := range pending {
			pendingHashText, _ := pendingHash.MarshalText()
			vote := DraftVoteAction{
				OnBehalfOf: s.Proposals.OnBehalfOf(pendingHash),
				Hash:       string(pendingHashText),
			}
			switch s.Proposals.Kind(pendingHash) {
			case state.DraftProposal:
				if pending, ok := s.Proposals.Draft[pendingHash]; ok && pending.DraftHash.Equal(hash) {
					vote.Kind = "Authorship"
					view.Votes = append(view.Votes, vote)
				}
			case state.ReleaseDraftProposal:
				if pending, ok := s.Proposals.ReleaseDraft[pendingHash]; ok && pending.Draft.DraftHash.Equal(hash) {
					vote.Kind = "Release"
					view.Votes = append(view.Votes, vote)
				}

			case state.PinProposal:
				if pending, ok := s.Proposals.Pin[pendingHash]; ok && pending.Hash.Equal(hash) {
					vote.Kind = "Pin"
					view.Votes = append(view.Votes, vote)
				}
			case state.ImprintStampProposal:
				if pending, ok := s.Proposals.Pin[pendingHash]; ok && pending.Hash.Equal(hash) {
					vote.Kind = "Stamp"
					view.Votes = append(view.Votes, vote)
				}
			}
		}
	}
	if release, ok := s.Releases[draft.DraftHash]; ok {
		view.Stamps = StampList(release.Stamps)
		view.Released = true
	}
	if len(draft.Edits) > 0 {
		view.Edited = true
	}
	if draft.PreviousVersion != nil {
		text, _ := draft.PreviousVersion.DraftHash.MarshalText()
		view.PreviousHash = string(text)
	}
	view.Policy.Majority, view.Policy.SuperMajority = draft.Authors.GetPolicy()
	return &view
}

// Edits template struct

type EditsView struct {
	Title   string
	Authors []string
	Reasons string
}

type EditsListView struct {
	Edits []EditsView
}

func EditsFromState(state *state.State) EditsListView {
	view := EditsListView{
		Edits: make([]EditsView, 0),
	}
	for _, edit := range state.Edits {
		itemView := EditsView{
			Title:   edit.Draft.Title,
			Authors: membersToHandles(edit.Authors.ListOfMembers(), state),
			Reasons: edit.Reasons,
		}
		view.Edits = append(view.Edits, itemView)
	}
	return view
}

// Votes template struct

type VotesView struct {
	Action  string
	Scope   string
	Hash    string
	Handler string
}

type VotesListView struct {
	Votes []VotesView
}

type VoteDetailView struct {
	Hash string
}

func VotesFromState(s *state.State, token crypto.Token) VotesListView {
	view := VotesListView{
		Votes: make([]VotesView, 0),
	}
	votes := s.Proposals.GetVotes(token)
	fmt.Println(len(votes))
	for hash := range votes {
		hashText, _ := hash.MarshalText()
		itemView := VotesView{
			Action: s.Proposals.KindText(hash),
			Scope:  s.Proposals.OnBehalfOf(hash),
			Hash:   string(hashText),
		}
		switch s.Proposals.Kind(hash) {
		case state.RequestMembershipProposal:
			itemView.Handler = "requestmembership"
		case state.DraftProposal:
			itemView.Handler = "draft"
		case state.PinProposal:
			itemView.Handler = "draft"
			prop := s.Proposals.Pin[hash]
			itemView.Hash = crypto.EncodeHash(prop.Draft.DraftHash)

		case state.ImprintStampProposal:
			itemView.Handler = "draft"
			prop := s.Proposals.ImprintStamp[hash]
			itemView.Hash = crypto.EncodeHash(prop.Release.Draft.DraftHash)

		case state.ReleaseDraftProposal:
			itemView.Handler = "draft"
			prop := s.Proposals.ReleaseDraft[hash]
			itemView.Hash = crypto.EncodeHash(prop.Draft.DraftHash)

		}
		view.Votes = append(view.Votes, itemView)
	}
	return view
}

type RequestMembershipView struct {
	Collective string
	Handle     string
	Hash       string
	Reasons    string
	Majority   string
}

func RequestMembershipFromState(s *state.State, hash crypto.Hash) *RequestMembershipView {
	vote, ok := s.Proposals.RequestMembership[hash]
	if !ok {
		return nil
	}
	token := vote.Request.Author
	handle, ok := s.Members[crypto.Hasher(token[:])]
	if !ok {
		return nil
	}
	hashText, _ := vote.Hash.MarshalText()
	majority, _ := vote.Collective.GetPolicy()
	return &RequestMembershipView{
		Collective: vote.Collective.Name,
		Handle:     handle,
		Hash:       string(hashText),
		Reasons:    vote.Request.Reasons,
		Majority:   fmt.Sprintf("%v", majority),
	}
}

type EditVersion struct {
	DraftHash string
}

func NewEdit(s *state.State, hash crypto.Hash) *EditVersion {
	draft, ok := s.Drafts[hash]
	if !ok {
		return nil
	}
	return &EditVersion{
		DraftHash: crypto.EncodeHash(draft.DraftHash),
	}
}

type DraftVersion struct {
	OnBehalfOf    string
	Policy        Policy
	Title         string
	Keywords      string
	Description   string
	PreviousDraft string
	References    string
}

func NewDraftVerion(s *state.State, hash crypto.Hash) *DraftVersion {
	draft, ok := s.Drafts[hash]
	if !ok {
		return &DraftVersion{}
	}
	majority, supermajority := draft.Authors.GetPolicy()
	return &DraftVersion{
		OnBehalfOf:    draft.Authors.CollectiveName(),
		Policy:        Policy{Majority: majority, SuperMajority: supermajority},
		Title:         draft.Title,
		Keywords:      strings.Join(draft.Keywords, ","),
		Description:   draft.Description,
		PreviousDraft: crypto.EncodeHash(hash),
	}
}

// func VoteDetailFromState(state *state.State, hash string) *VoteDetailView {
// 	ok := state.Vote(hash)
// 	if !ok {
// 		return nil
// 	}
// 	view := VoteDetailView{
// 		Hash: hash,
// 	}
// 	return &view
// }
