package api

import (
	"fmt"
	"strings"

	"github.com/lienkolabs/swell/crypto"
	"github.com/lienkolabs/synergy/social/state"
)

type HeaderInfo struct {
	Active  string
	Path    string
	EndPath string
	Section string
}

// Drafts template struct

type DraftsView struct {
	Title       string
	Authors     []AuthorDetail
	Hash        string
	Description string
	Keywords    []string
}

type DraftsListView struct {
	Drafts []DraftsView
	Head   HeaderInfo
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
	Head         HeaderInfo
}

type EditDetailedView struct {
	DraftTitle string
	DraftHash  string
	Reasons    string
	Hash       string
	Authors    []AuthorDetail
	Votes      []DraftVoteAction
	Head       HeaderInfo
}

func EditDetailFromState(s *state.State, hash crypto.Hash, token crypto.Token) *EditDetailedView {
	edit, ok := s.Edits[hash]
	if !ok {
		edit, ok = s.Proposals.Edit[hash]
		if !ok {
			return nil
		}
	}
	head := HeaderInfo{
		Active:  "MyDrafts",
		Path:    "venture > my drafts > " + edit.Draft.Title + " > ",
		EndPath: "edits",
		Section: "venture",
	}
	view := EditDetailedView{
		DraftTitle: edit.Draft.Title,
		DraftHash:  crypto.EncodeHash(edit.Draft.DraftHash),
		Reasons:    edit.Reasons,
		Hash:       crypto.EncodeHash(edit.Edit),
		Authors:    AuthorList(edit.Authors, s),
		Votes:      make([]DraftVoteAction, 0),
		Head:       head,
	}
	pending := s.Proposals.GetVotes(token)
	if len(pending) > 0 {
		for pendingHash := range pending {
			if s.Proposals.Kind(pendingHash) == state.EditProposal {
				vote := DraftVoteAction{
					Kind:       "Authorship",
					OnBehalfOf: s.Proposals.OnBehalfOf(pendingHash),
					Hash:       crypto.EncodeHash(pendingHash),
				}
				view.Votes = append(view.Votes, vote)
			}
		}
	}
	return &view
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
	head := HeaderInfo{
		Active:  "Drafts",
		Path:    "explore > ",
		EndPath: "drafts",
		Section: "explore",
	}
	view := DraftsListView{
		Head:   head,
		Drafts: make([]DraftsView, 0),
	}
	for _, draft := range state.Drafts {
		hash, _ := draft.DraftHash.MarshalText()
		itemView := DraftsView{
			Title:       draft.Title,
			Hash:        string(hash),
			Authors:     AuthorList(draft.Authors, state),
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
	if view.Authorship {
		view.Head = HeaderInfo{
			Active:  "MyDrafts",
			Path:    "venture > my drafts > ",
			EndPath: draft.Title,
			Section: "venture",
		}
	} else {
		view.Head = HeaderInfo{
			Active:  "Drafts",
			Path:    "explore > drafts > ",
			EndPath: draft.Title,
			Section: "explore",
		}
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
				if pending, ok := s.Proposals.Pin[pendingHash]; ok && pending.Draft.DraftHash.Equal(hash) {
					vote.Kind = "Pin"
					view.Votes = append(view.Votes, vote)
				}
			case state.ImprintStampProposal:
				if pending, ok := s.Proposals.ImprintStamp[pendingHash]; ok && pending.Release.Hash.Equal(hash) {
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
	Authors []AuthorDetail
	Reasons string
	Hash    string
}

type EditsListView struct {
	DraftTitle string
	DraftHash  string
	Edits      []EditsView
	Head       HeaderInfo
}

func EditsFromState(s *state.State, drafthash crypto.Hash) EditsListView {
	draft, ok := s.Drafts[drafthash]
	if !ok {
		return EditsListView{}
	}
	head := HeaderInfo{
		Active:  "MyDrafts",
		Path:    "venture > my drafts > " + draft.Title + " > ",
		EndPath: "edits",
		Section: "venture",
	}
	view := EditsListView{
		DraftTitle: draft.Title,
		DraftHash:  crypto.EncodeHash(draft.DraftHash),
		Edits:      make([]EditsView, 0),
		Head:       head,
	}
	for _, edit := range draft.Edits {
		itemView := EditsView{
			Authors: AuthorList(edit.Authors, s),
			Reasons: edit.Reasons,
			Hash:    crypto.EncodeHash(edit.Edit),
		}
		view.Edits = append(view.Edits, itemView)
	}
	return view
}

// Votes template struct

type VotesView struct {
	Action            string
	Scope             string
	Hash              string
	Handler           string
	ObjectType        string
	ObjectLink        string
	ObjectCaption     string
	ComplementType    string
	ComplementLink    string
	ComplementCaption string
}

type VotesListView struct {
	Votes []VotesView
	Head  HeaderInfo
}

type VoteDetailView struct {
	Hash string
}

func VotesFromState(s *state.State, token crypto.Token) VotesListView {
	head := HeaderInfo{
		Active:  "Votes",
		Path:    "venture > ",
		EndPath: "consensus votes",
		Section: "venture",
	}
	view := VotesListView{
		Head:  head,
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
			prop := s.Proposals.RequestMembership[hash]
			handle, ok := s.Members[crypto.Hasher(prop.Request.Author[:])]
			if ok {
				itemView.ObjectCaption = handle
				itemView.ObjectLink = fmt.Sprintf("/member/%v", handle)
				itemView.ObjectType = ""
			}

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

		case state.UpdateCollectiveProposal:
			itemView.Handler = "voteupdatecollective"
			//prop := s.Proposals.UpdateCollective[hash]
			//itemView.Hash = crypto.EncodeHash(prop.Hash)

		case state.RemoveMemberProposal:
			prop := s.Proposals.RemoveMember[hash]
			handle, ok := s.Members[crypto.Hasher(prop.Remove.Member[:])]
			if ok {
				itemView.ObjectCaption = handle
				itemView.ObjectLink = fmt.Sprintf("/member/%v", handle)
				itemView.ObjectType = ""
			}
		case state.EditProposal:
			itemView.Handler = "editview"
		case state.CreateBoardProposal:
			itemView.Handler = "votecreateboard"
		case state.UpdateBoardProposal:
			itemView.Handler = "voteupdateboard"
		case state.BoardEditorProposal:
			prop := s.Proposals.BoardEditor[hash]
			editor, ok := s.Members[crypto.Hasher(prop.Editor[:])]
			if ok {
				itemView.ObjectCaption = editor
				itemView.ObjectLink = fmt.Sprintf("/member/%v", editor)
				if prop.Insert {
					itemView.ObjectType = "include"
				} else {
					itemView.ObjectType = "remove"
				}
			}
			itemView.Scope = ""
			itemView.ComplementCaption = prop.Board.Name
			itemView.ComplementType = "board"
			itemView.ComplementLink = fmt.Sprintf("/board/%v", prop.Board.Name)
		case state.ReactProposal:
		case state.CreateEventProposal:
			itemView.Handler = "event"
		case state.CancelEventProposal:
			itemView.Handler = "event"
			prop := s.Proposals.CancelEvent[hash]
			itemView.Hash = crypto.EncodeHash(prop.Event.Hash)
		case state.UpdateEventProposal:
			itemView.Handler = "voteupdateevent"

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
	Head      HeaderInfo
}

func NewEdit(s *state.State, hash crypto.Hash) *EditVersion {
	draft, ok := s.Drafts[hash]
	if !ok {
		return nil
	}
	head := HeaderInfo{
		Active:  "Draft",
		Path:    "venture > drafts > " + draft.Title + " > venture > ",
		EndPath: "edit",
		Section: "venture",
	}
	return &EditVersion{
		DraftHash: crypto.EncodeHash(draft.DraftHash),
		Head:      head,
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
	Head          HeaderInfo
}

func NewDraftVersion(s *state.State, hash crypto.Hash) *DraftVersion {
	draft, ok := s.Drafts[hash]
	if !ok {
		return &DraftVersion{}
	}
	majority, supermajority := draft.Authors.GetPolicy()
	head := HeaderInfo{
		Active:  "NewDraft",
		Path:    "venture > ",
		EndPath: "new draft",
		Section: "venture",
	}
	return &DraftVersion{
		OnBehalfOf:    draft.Authors.CollectiveName(),
		Policy:        Policy{Majority: majority, SuperMajority: supermajority},
		Title:         draft.Title,
		Keywords:      strings.Join(draft.Keywords, ","),
		Description:   draft.Description,
		PreviousDraft: crypto.EncodeHash(hash),
		Head:          head,
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

// Boards template struct

type CollectiveUpdateView struct {
	Name             string
	OldDescription   string
	Description      string
	OldMajority      int
	Majority         int
	OldSuperMajority int
	SuperMajority    int
	Member           bool
	Hash             string
	Reasons          string
	Head             HeaderInfo
}

func CollectiveToUpdateFromState(s *state.State, name string) *CollectiveUpdateView {
	col, ok := s.Collective(name)
	if !ok {
		return nil
	}
	head := HeaderInfo{
		Active:  "Central",
		Path:    "venture > central > collectives " + name + " > ",
		EndPath: "update collective",
		Section: "venture",
	}
	update := &CollectiveUpdateView{
		Name:             name,
		OldDescription:   col.Description,
		OldMajority:      col.Policy.Majority,
		OldSuperMajority: col.Policy.SuperMajority,
		Head:             head,
	}
	return update
}

func CollectiveUpdateFromState(s *state.State, hash crypto.Hash, token crypto.Token) *CollectiveUpdateView {
	pending, ok := s.Proposals.UpdateCollective[hash]
	if !ok {
		return nil
	}
	live := pending.Collective
	update := &CollectiveUpdateView{
		Name:             live.Name,
		OldDescription:   live.Description,
		OldMajority:      live.Policy.Majority,
		OldSuperMajority: live.Policy.SuperMajority,
		Hash:             crypto.EncodeHash(hash),
		Reasons:          pending.Update.Reasons,
	}
	if pending.Update.Description != nil {
		update.Description = *pending.Update.Description
	}
	if pending.Update.Majority != nil {
		update.Majority = int(*pending.Update.Majority)
	}
	if pending.Update.SuperMajority != nil {
		update.SuperMajority = int(*pending.Update.SuperMajority)
	}
	if live.IsMember(token) {
		update.Member = true
	}

	return update
}

type BoardUpdateView struct {
	Name              string
	Collective        string
	Description       string
	OldDescription    string
	KeywordsString    string
	OldKeywordsString string
	PinMajority       byte
	OldPinMajority    byte
	Reasons           string
	Hash              string
	Head              HeaderInfo
}

func BoardToUpdateFromState(s *state.State, name string) *BoardUpdateView {
	live, ok := s.Board(name)
	if !ok {
		return nil
	}
	head := HeaderInfo{
		Active:  "Central",
		Path:    "central > boards > ",
		EndPath: live.Name,
		Section: "venture",
	}
	update := &BoardUpdateView{
		Name:           live.Name,
		Collective:     live.Collective.Name,
		OldDescription: live.Description,
		OldPinMajority: byte(live.Editors.Majority),
		Head:           head,
	}
	if len(live.Keyword) > 0 {
		update.OldKeywordsString = strings.Join(live.Keyword, ",")
	}
	return update
}

func BoardUpdateFromState(s *state.State, hash crypto.Hash) *BoardUpdateView {
	pending, ok := s.Proposals.UpdateBoard[hash]
	if !ok {
		return nil
	}
	live := pending.Board
	update := &BoardUpdateView{
		Name:           live.Name,
		Collective:     live.Collective.Name,
		OldDescription: live.Description,
		OldPinMajority: byte(live.Editors.Majority),
		Reasons:        pending.Origin.Reasons,
		Hash:           crypto.EncodeHash(pending.Hash),
	}
	if pending.Description != nil {
		update.Description = *pending.Description
	}
	if pending.PinMajority != nil {
		update.PinMajority = *pending.PinMajority
	}

	if len(live.Keyword) > 0 {
		update.OldKeywordsString = strings.Join(live.Keyword, ",")
	}
	if pending.Keywords != nil {
		update.KeywordsString = strings.Join(*pending.Keywords, ",")
	}
	return update
}

type BoardsView struct {
	Name       string
	Hash       string
	Collective string
	Keywords   []string
}

type BoardsListView struct {
	Boards []BoardsView
	Head   HeaderInfo
}

type BoardDetailView struct {
	Name        string
	Description string
	Collective  string
	Keywords    []string
	PinMajority int
	Editors     []string
	Drafts      []DraftsView
	Editorship  bool
	Reasons     string
	Author      string
	Hash        string
	Head        HeaderInfo
}

func BoardsFromState(s *state.State) BoardsListView {
	head := HeaderInfo{
		Active:  "Boards",
		Path:    "explore > ",
		EndPath: "boards",
		Section: "explore",
	}
	view := BoardsListView{
		Head:   head,
		Boards: make([]BoardsView, 0),
	}
	for _, board := range s.Boards {
		hash, _ := board.Hash.MarshalText()
		itemView := BoardsView{
			Name:       board.Name,
			Hash:       string(hash),
			Collective: board.Collective.Name,
			Keywords:   board.Keyword,
		}
		view.Boards = append(view.Boards, itemView)
	}
	return view
}

func PendingBoardFromState(s *state.State, hash crypto.Hash) *BoardDetailView {
	pending, ok := s.Proposals.CreateBoard[hash]
	if !ok {
		return nil
	}
	board := pending.Board
	view := BoardDetailView{
		Name:        board.Name,
		Description: board.Description,
		Collective:  board.Collective.Name,
		Keywords:    board.Keyword,
		PinMajority: board.Editors.Majority,
		Editors:     make([]string, 0),
		Drafts:      make([]DraftsView, 0),
		Reasons:     pending.Origin.Reasons,
		Hash:        crypto.EncodeHash(hash),
	}
	view.Author = s.Members[crypto.Hasher(pending.Origin.Author[:])]
	return &view
}

func BoardDetailFromState(s *state.State, name string, token crypto.Token) *BoardDetailView {
	board, ok := s.Board(name)
	if !ok {
		return nil
	}
	view := BoardDetailView{
		Name:        board.Name,
		Description: board.Description,
		Collective:  board.Collective.Name,
		Keywords:    board.Keyword,
		PinMajority: board.Editors.Majority,
		Editors:     make([]string, 0),
		Drafts:      make([]DraftsView, 0),
		Editorship:  board.Editors.IsMember(token),
	}
	if view.Editorship {
		view.Head = HeaderInfo{
			Active:  "Central",
			Path:    "venture > central > boards > ",
			EndPath: board.Name,
			Section: "venture",
		}
	} else {
		view.Head = HeaderInfo{
			Active:  "Boards",
			Path:    "explore > boards > ",
			EndPath: board.Name,
			Section: "explore",
		}
	}
	for token, _ := range board.Editors.Members {
		handle, ok := s.Members[crypto.Hasher(token[:])]
		if ok {
			view.Editors = append(view.Editors, handle)
		}
	}
	for _, d := range board.Pinned {
		draftView := DraftsView{
			Title:       d.Title,
			Authors:     make([]AuthorDetail, 0),
			Hash:        crypto.EncodeHash(d.DraftHash),
			Description: d.Description,
			Keywords:    d.Keywords,
		}
		view.Drafts = append(view.Drafts, draftView)
	}

	return &view
}

// Collectives template struct

type CollectivesView struct {
	Name         string
	Description  string
	Participants string
}

type CollectivesListView struct {
	Collectives []CollectivesView
	Head        HeaderInfo
}

type CollectiveDetailView struct {
	Name          string
	Description   string
	Majority      int
	SuperMajority int
	Members       []MemberDetailView
	Membership    bool
	Head          HeaderInfo
}

func ColletivesFromState(s *state.State) CollectivesListView {
	head := HeaderInfo{
		Active:  "Collectives",
		Path:    "explore > ",
		EndPath: "collectives",
		Section: "explore",
	}
	view := CollectivesListView{
		Head:        head,
		Collectives: make([]CollectivesView, 0),
	}
	for _, collective := range s.Collectives {
		itemView := CollectivesView{
			Name:         collective.Name,
			Description:  collective.Description,
			Participants: fmt.Sprintf("%v", len(collective.Members)),
		}
		view.Collectives = append(view.Collectives, itemView)
	}
	return view
}

func CollectiveDetailFromState(s *state.State, name string, token crypto.Token) *CollectiveDetailView {
	collective, ok := s.Collective(name)
	if !ok {
		return nil
	}
	view := CollectiveDetailView{
		Name:          collective.Name,
		Description:   collective.Description,
		Majority:      collective.Policy.Majority,
		SuperMajority: collective.Policy.SuperMajority,
		Members:       make([]MemberDetailView, 0),
		Membership:    collective.IsMember(token),
	}
	if view.Membership {
		view.Head = HeaderInfo{
			Active:  "Central",
			Path:    "venture > central > collectives > ",
			EndPath: name,
			Section: "venture",
		}
	} else {
		view.Head = HeaderInfo{
			Active:  "Collectives",
			Path:    "explore > collectives > ",
			EndPath: name,
			Section: "explore",
		}
	}
	for token, _ := range collective.Members {
		handle, ok := s.Members[crypto.Hasher(token[:])]
		if ok {
			view.Members = append(view.Members, MemberDetailView{handle})
		}
	}
	return &view
}
