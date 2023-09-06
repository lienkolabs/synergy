package api

import (
	"fmt"
	"time"

	"github.com/lienkolabs/swell/crypto"
	"github.com/lienkolabs/swell/crypto/dh"
	"github.com/lienkolabs/synergy/social/index"
	"github.com/lienkolabs/synergy/social/state"
)

type EventsView struct {
	Hash        string
	Live        bool
	Description string
	StartAt     time.Time
	Collective  string
	Public      bool
}

type EventVoteAction struct {
	Kind   string // create, cancel, update
	Hash   string
	Update string
}

type EventsListView struct {
	Events []EventsView
	Head   HeaderInfo
}

type VoteUpdateEventView struct {
	Description     string
	OldDescription  string
	StartAt         string
	OldStartAt      string
	EstimatedEnd    string
	OldEstimatedEnd string
	Venue           string
	OldVenue        string
	Open            string
	OldOpen         string
	Public          string
	OldPublic       string
	Hash            string
	Reasons         string
	Collective      string
	Managing        bool
	VoteHash        string
	Head            HeaderInfo
}

func yesorno(b *bool) string {
	if b == nil {
		return ""
	}
	if *b {
		return "yes"
	}
	return "no"
}

func EventUpdateFromState(s *state.State, hash crypto.Hash, token crypto.Token) VoteUpdateEventView {
	update, ok := s.Proposals.UpdateEvent[hash]
	if !ok {
		return VoteUpdateEventView{}
	}
	old := update.Event
	head := HeaderInfo{
		Active:  "MyEvents",
		Path:    "venture > my events > ",
		EndPath: "update event " + old.StartAt.Format("2006-01-02") + " by " + old.Collective.Name,
		Section: "explore",
	}
	vote := VoteUpdateEventView{
		OldDescription:  old.Description,
		OldStartAt:      old.StartAt.String(),
		OldEstimatedEnd: old.EstimatedEnd.String(),
		OldVenue:        old.Venue,
		OldOpen:         yesorno(&old.Open),
		OldPublic:       yesorno(&old.Public),
		Open:            yesorno(update.Open),
		Public:          yesorno(update.Public),
		Hash:            crypto.EncodeHash(old.Hash),
		Reasons:         update.Reasons,
		Collective:      old.Collective.Name,
		VoteHash:        crypto.EncodeHash(hash),
		Head:            head,
	}
	if update.Description != nil {
		vote.Description = *update.Description
	}
	if update.StartAt != nil {
		vote.StartAt = update.StartAt.String()
	}
	if update.EstimatedEnd != nil {
		vote.StartAt = update.StartAt.String()
	}
	if update.Venue != nil {
		vote.StartAt = *update.Venue
	}
	if old.Managers.IsMember(token) {
		vote.Managing = true
		vote.Head = HeaderInfo{
			Active:  "MyEvents",
			Path:    "venture > my events > ",
			EndPath: old.StartAt.Format("2006-01-02") + " by " + old.Collective.Name,
			Section: "venture",
		}
	} else {
		vote.Head = HeaderInfo{
			Active:  "Events",
			Path:    "explore > events > ",
			EndPath: old.StartAt.Format("2006-01-02") + " by " + old.Collective.Name,
			Section: "explore",
		}
	}

	return vote
}

type EventDetailView struct {
	Live            bool
	Description     string
	StartAt         time.Time
	EstimatedEnd    time.Time
	Collective      string
	Venue           string
	Open            bool
	Public          bool
	ManagerMajority int
	Managers        []MemberDetailView
	Checkedin       []CheckInDetails
	Votes           []EventVoteAction
	Managing        bool
	Hash            string
	Greeted         []MemberDetailView
	MyGreeting      string
	Head            HeaderInfo
}

func EventsFromState(state *state.State) EventsListView {
	head := HeaderInfo{
		Active:  "Events",
		Path:    "explore > ",
		EndPath: "events",
		Section: "explore",
	}
	view := EventsListView{
		Head:   head,
		Events: make([]EventsView, 0),
	}
	for _, event := range state.Events {
		itemView := EventsView{
			Hash: crypto.EncodeHash(event.Hash),
			Live: event.Live,

			Description: event.Description,
			StartAt:     event.StartAt,
			Collective:  event.Collective.Name,
			Public:      event.Public,
		}
		view.Events = append(view.Events, itemView)
	}
	return view
}

type CheckInDetails struct {
	Handle       string
	EphemeralKey string
}

func EventDetailFromState(s *state.State, hash crypto.Hash, token crypto.Token, ephemeral crypto.PrivateKey) *EventDetailView {
	event, ok := s.Events[hash]
	if !ok {
		event = s.Proposals.GetEvent(hash)
		if event == nil {
			return nil
		}
	}
	view := EventDetailView{
		Live:            event.Live,
		Description:     event.Description,
		StartAt:         event.StartAt,
		Collective:      event.Collective.Name,
		EstimatedEnd:    event.EstimatedEnd,
		Venue:           event.Venue,
		Open:            event.Open,
		Public:          event.Public,
		Checkedin:       make([]CheckInDetails, 0),
		ManagerMajority: event.Managers.Majority,
		Managers:        make([]MemberDetailView, 0),
		Votes:           make([]EventVoteAction, 0),
		Managing:        event.Managers.IsMember(token),
		Hash:            crypto.EncodeHash(hash),
	}
	if view.Managing {
		view.Head = HeaderInfo{
			Active:  "MyEvents",
			Path:    "venture > my events > ",
			EndPath: event.StartAt.Format("2006-01-02") + " by " + event.Collective.Name,
			Section: "venture",
		}
	} else {
		view.Head = HeaderInfo{
			Active:  "Events",
			Path:    "explore > events > ",
			EndPath: event.StartAt.Format("2006-01-02") + " by " + event.Collective.Name,
			Section: "explore",
		}
	}
	for token, _ := range event.Managers.ListOfMembers() {
		handle, ok := s.Members[crypto.Hasher(token[:])]
		if ok {
			view.Managers = append(view.Managers, MemberDetailView{handle})
		}
	}
	for token, greet := range event.Checkin {
		if handle, ok := s.Members[crypto.Hasher(token[:])]; ok {
			if greet != nil && greet.Action != nil {
				view.Greeted = append(view.Greeted, MemberDetailView{handle})
				// if its me, de-crypt message
				if greet.Action.CheckedIn.Equal(token) {
					fmt.Println("-0")
					dhCipher := dh.ConsensusCipher(ephemeral, greet.Action.EphemeralToken)
					if secretKey, err := dhCipher.Open(greet.Action.SecretKey); err == nil {
						fmt.Println("-1")
						cipher := crypto.CipherFromKey(secretKey)
						if content, err := cipher.Open(greet.Action.PrivateContent); err == nil {
							view.MyGreeting = string(content)
							fmt.Println("-3", view.MyGreeting)
						}
					}
				}
			} else {
				bytes, _ := greet.EphemeralKey.MarshalText()
				view.Checkedin = append(view.Checkedin, CheckInDetails{Handle: handle, EphemeralKey: string(bytes)})
			}
		}
	}
	pending := s.Proposals.GetVotes(token)
	if len(pending) > 0 {
		for pendingHash := range pending {
			vote := EventVoteAction{
				Hash: crypto.EncodeHash(pendingHash),
			}
			switch s.Proposals.Kind(pendingHash) {
			case state.CreateEventProposal:
				// collective vote
				vote.Kind = "Create"
				view.Votes = append(view.Votes, vote)
			case state.UpdateEventProposal:
				// managers vote
				vote.Kind = "Update"
				vote.Update = "Update"
				view.Votes = append(view.Votes, vote)
			case state.CancelEventProposal:
				// managers vote
				vote.Kind = "Cancel"
				view.Votes = append(view.Votes, vote)
			}
		}
	}
	return &view
}

func EventUpdateDetailFromState(s *state.State, hash crypto.Hash, token crypto.Token) *EventDetailView {
	event, ok := s.Events[hash]
	if !ok {
		event = s.Proposals.GetEvent(hash)
		if event == nil {
			return nil
		}
	}
	head := HeaderInfo{
		Active:  "MyEvents",
		Path:    "venture > my events > " + event.StartAt.Format("2006-01-02") + " by " + event.Collective.Name + " > ",
		EndPath: "update",
		Section: "venture",
	}
	view := EventDetailView{
		StartAt:         event.StartAt,
		Live:            event.Live,
		Description:     event.Description,
		Collective:      event.Collective.Name,
		EstimatedEnd:    event.EstimatedEnd,
		Venue:           event.Venue,
		Open:            event.Open,
		Public:          event.Public,
		ManagerMajority: event.Managers.Majority,
		Managers:        make([]MemberDetailView, 0),
		Votes:           make([]EventVoteAction, 0),
		Managing:        event.Managers.IsMember(token),
		Hash:            crypto.EncodeHash(hash),
		Head:            head,
	}
	for token, _ := range event.Managers.ListOfMembers() {
		handle, ok := s.Members[crypto.Hasher(token[:])]
		if ok {
			view.Managers = append(view.Managers, MemberDetailView{handle})
		}
	}
	pending := s.Proposals.GetVotes(token)
	if len(pending) > 0 {
		for pendingHash := range pending {
			vote := EventVoteAction{
				Hash: crypto.EncodeHash(pendingHash),
			}
			if vote.Kind != "Update" {
				view.Votes = append(view.Votes, vote)
			}
		}
	}
	return &view
}

// Members template struct

type MembersView struct {
	Hash   string
	Handle string
}

type MembersListView struct {
	Members []MembersView
	Head    HeaderInfo
}

type MemberDetailView struct {
	Handle string
}

type MemberDetailViewPage struct {
	Detail MemberDetailView
	Head   HeaderInfo
}

func MembersFromState(state *state.State) MembersListView {
	head := HeaderInfo{
		Active:  "Members",
		Path:    "explore > ",
		EndPath: "members",
		Section: "explore",
	}
	view := MembersListView{
		Head:    head,
		Members: make([]MembersView, 0),
	}
	for hash, member := range state.Members {
		hashText, _ := hash.MarshalText()
		itemView := MembersView{
			Hash:   string(hashText),
			Handle: member,
		}
		view.Members = append(view.Members, itemView)
	}
	return view
}

func MemberDetailFromState(state *state.State, handle string) *MemberDetailViewPage {
	_, ok := state.MembersIndex[handle]
	if !ok {
		return nil
	}
	detail := MemberDetailView{
		Handle: handle,
	}
	head := HeaderInfo{
		Active:  "Members",
		Path:    "explore > members > ",
		EndPath: handle,
		Section: "explore",
	}
	view := MemberDetailViewPage{
		Detail: detail,
		Head:   head,
	}
	return &view
}

// Central Connections

type LastAction struct {
	Type        string
	Handle      string
	TimeOfInstr time.Time
}

type LastReference struct {
	Author      string
	TimeOfInstr time.Time
}

type CentralCollectives struct {
	Name     string
	NBoards  int
	NStamps  int
	NEvents  int
	LastSelf LastAction
	LastAny  LastAction
}

type CentralBoards struct {
	Name     string
	NPins    int
	NEditors int
	LastSelf LastAction
	LastAny  LastAction
}

type CentralEvents struct {
	DateCol      string //data e horario mais nome do coletivo
	NCheckins    int
	NPenCheckins int
	LastSelf     LastAction
	LastAny      LastAction
}

type CentralEdits struct {
	Title       string
	CreatedAt   time.Time
	NReferences int
	LastRef     LastReference
}

type CentralConnectionsListView struct {
	Head         HeaderInfo
	Collectives  []CentralCollectives
	NCollectives int
	Boards       []CentralBoards
	NBoards      int
	Events       []CentralEvents
	NEvents      int
	Edits        []CentralEdits
	NEdits       int
}

func CentralConnectionsFromState(state *state.State, indexer *index.Index, token crypto.Token) CentralConnectionsListView {
	head := HeaderInfo{
		Active:  "CentralConnections",
		Path:    "venture > ",
		EndPath: "central connections",
		Section: "venture",
	}
	view := CentralConnectionsListView{
		Head:        head,
		Collectives: make([]CentralCollectives, 0),
		Boards:      make([]CentralBoards, 0),
		Events:      make([]CentralEvents, 0),
		Edits:       make([]CentralEdits, 0),
	}
	// Check collectives user is a member of and get their info

	// for _, collective := range state.Collectives {
	memberscol := indexer.CollectivesOnMember(token)
	for _, collective := range memberscol {
		// if collective.IsMember(token) {
		nboards := len(indexer.BoardsOnCollective(collective))
		nstamps := 0
		for _, release := range state.Releases {
			for _, stamp := range release.Stamps {
				if stamp.Reputation.Name == collective.Name {
					nstamps++
				}
			}
		}
		nevents := 0
		for _, event := range state.Events {
			if event.Collective.Name == collective.Name {
				nevents++
			}
		}
		item := CentralCollectives{
			Name:    collective.Name,
			NBoards: nboards,
			NStamps: nstamps,
			NEvents: nevents,
		}
		lastaction := indexer.LastMemberActionOnCollective(token, collective.Name)
		if lastaction != nil {
			fmt.Println(lastaction.Description)
			item.LastSelf = LastAction{
				Type:        lastaction.Description,
				Handle:      state.Members[crypto.HashToken(token)],
				TimeOfInstr: state.TimeOfEpoch(lastaction.Epoch),
			}
		}
		view.Collectives = append(view.Collectives, item)
		// }
	}
	view.NCollectives = len(view.Collectives)
	// Check boards user is an editor at and get their info
	for _, board := range state.Boards {
		if board.Editors.IsMember(token) {
			item := CentralBoards{
				Name:     board.Name,
				NPins:    len(board.Pinned),
				NEditors: len(board.Editors.ListOfMembers()),
			}
			view.Boards = append(view.Boards, item)
		}
	}
	view.NBoards = len(view.Boards)
	// Check events user is a manager on and get their info
	for _, event := range state.Events {
		if event.Managers.IsMember(token) {
			eventname := event.StartAt.Format("2006-01-02") + " by " + event.Collective.Name
			ncheckins := len(event.Checkin)
			ngreets := 0
			for _, greet := range event.Checkin {
				if greet != nil && greet.Action != nil {
					ngreets++
				}
			}
			item := CentralEvents{
				DateCol:   eventname,
				NCheckins: ncheckins,
				// NPenCheckins: ncheckins - ngreets,
			}
			view.Events = append(view.Events, item)
		}
	}
	view.NEvents = len(view.Events)
	// Check edits user has proposed and get their info
	for _, edit := range state.Edits {
		if edit.Authors.IsMember(token) {
			item := CentralEdits{
				Title: edit.Draft.Title,
			}
			view.Edits = append(view.Edits, item)
		}
	}
	view.NEdits = len(view.Edits)
	return view
}
