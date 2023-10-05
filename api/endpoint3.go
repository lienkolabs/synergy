package api

import (
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/lienkolabs/breeze/crypto"
	"github.com/lienkolabs/synergy/social/index"
	"github.com/lienkolabs/synergy/social/state"
)

type UpdatesView struct {
	Objects []ObjectUpdateView
	Head    HeaderInfo
}

func (u *UpdatesView) Len() int {
	return len(u.Objects)
}

func (u *UpdatesView) Less(i, j int) bool {
	return u.Objects[i].LastUpdated().Before(u.Objects[j].LastUpdated())
}

func (u *UpdatesView) Swap(i, j int) {
	u.Objects[i], u.Objects[j] = u.Objects[j], u.Objects[i]
}

type ObjectUpdateView struct {
	Name       string
	ObjectKind string
	Updates    []ActionUpdateView
}

func (o ObjectUpdateView) LastUpdated() time.Time {
	var last time.Time
	for _, update := range o.Updates {
		if update.LastUpdatedTime.After(last) {
			last = update.LastUpdatedTime
		}
	}
	return last
}

type ActionUpdateView struct {
	Description         string
	VoteStatus          string
	VoteHash            string
	LastUpdatedInterval string
	LastUpdatedTime     time.Time
}

func actionsToActionUpdateView(actions []index.ActionDetails, genesisTime time.Time, token crypto.Token) []ActionUpdateView {
	updates := make([]ActionUpdateView, 0)
	for _, action := range actions {
		actionTime := genesisTime.Add(time.Second * time.Duration(action.Epoch))
		actionUpdateView := ActionUpdateView{
			Description:         action.Description,
			LastUpdatedTime:     actionTime,
			LastUpdatedInterval: PrettyDuration(time.Since(actionTime)),
		}

		if action.VoteStatus {
			actionUpdateView.VoteStatus = "approved"
		} else {
			actionUpdateView.VoteStatus = "pending vote"
		}
		if len(action.Votes) > 0 && (!action.VoteStatus) {
			hasCast := false
			for _, vote := range action.Votes {
				if vote.Author.Equal(token) {
					hasCast = true
					break
				}
			}
			if !hasCast {
				actionUpdateView.VoteHash = crypto.EncodeHash(action.Votes[0].Hash)
			}
		}
		updates = append(updates, actionUpdateView)
	}
	return updates
}

func UpdatesViewFromState(s *state.State, i *index.Index, token crypto.Token, genesisTime time.Time) *UpdatesView {
	head := HeaderInfo{
		Active:  "Updates",
		Path:    "venture > ",
		EndPath: "updates",
		Section: "venture",
	}
	collectives := i.CollectivesOnMember(token)
	boards := i.BoardsOnMember(token)
	events := i.EventsOnMember(token)

	objects := make([]ObjectUpdateView, 0)
	for _, collective := range collectives {
		actions := i.GetRecentActionsWithLinks(crypto.Hasher([]byte(collective)))
		updates := actionsToActionUpdateView(actions, genesisTime, token)
		if len(updates) > 0 {
			objView := ObjectUpdateView{
				Name:       collective,
				ObjectKind: "collective",
				Updates:    updates,
			}
			objects = append(objects, objView)
		}
	}
	for _, board := range boards {
		actions := i.GetRecentActionsWithLinks(crypto.Hasher([]byte(board)))
		updates := actionsToActionUpdateView(actions, genesisTime, token)
		if len(updates) > 0 {
			objView := ObjectUpdateView{
				Name:       board,
				ObjectKind: "board",
				Updates:    updates,
			}
			objects = append(objects, objView)
		}
	}
	for _, eventHash := range events {
		if event, ok := s.Events[eventHash]; ok {
			actions := i.GetRecentActionsWithLinks(eventHash)
			updates := actionsToActionUpdateView(actions, genesisTime, token)
			if len(updates) > 0 {
				objView := ObjectUpdateView{
					Name:       fmt.Sprintf("%s event from %s", event.StartAt.Format(time.RFC822), event.Collective.Name),
					ObjectKind: "event",
					Updates:    updates,
				}
				objects = append(objects, objView)
			}
		}
	}
	if len(objects) == 0 {
		return &UpdatesView{
			Objects: objects,
			Head:    head,
		}
	}
	updatesView := &UpdatesView{
		Objects: objects,
		Head:    head,
	}
	sort.Sort(updatesView)
	return updatesView
}

type PendingActionsView struct {
	Pending []PendingActionDetailView
	Head    HeaderInfo
}

type PendingActionDetailView struct {
	Description  string
	ProposedAt   string
	VotesApprove int
	VotesReject  int
	VotesNeeeded int
	VoteHash     string
}

func PendingActionsFromState(s *state.State, i *index.Index, token crypto.Token, genesisTime time.Time) *PendingActionsView {
	head := HeaderInfo{
		Active:  "Pending",
		Path:    "venture > ",
		EndPath: "pending actions",
		Section: "venture",
	}
	view := PendingActionsView{
		Pending: make([]PendingActionDetailView, 0),
		Head:    head,
	}
	pendingActions := i.GetPendingActionsDetailed(token)
	if len(pendingActions) == 0 {
		return nil
	}
	for _, pending := range pendingActions {
		proposed := genesisTime.Add(time.Duration(pending.Epoch) * time.Second)
		item := PendingActionDetailView{
			Description:  pending.Description,
			ProposedAt:   PrettyDuration(time.Since(proposed)),
			VotesNeeeded: len(pending.Pool.Voters) * pending.Pool.Majority / 100,
			VoteHash:     crypto.EncodeHash(pending.Pool.Votes[0].Hash),
		}
		for _, vote := range pending.Pool.Votes {
			if vote.Approve {
				item.VotesApprove += 1
			} else {
				item.VotesReject += 1
			}
		}
		view.Pending = append(view.Pending, item)
	}
	return &view
}

type MyEditView struct {
	DraftTitle  string
	DraftHash   string
	Hash        string
	PublishedAt string
	AuthorType  string
}

type EditOnDraftView struct {
	Caption string
	Link    string
	Time    string
}

type MyDraftView struct {
	Title       string
	Hash        string
	PublishedAt string
	Release     string
	Pinned      []CaptionLink
	Edit        []EditOnDraftView
	Stamps      []CaptionLink
	AuthorType  string
}

type MyMediaView struct {
	Drafts []MyDraftView
	Edits  []MyEditView
	Head   HeaderInfo
}

func MyMediaFromState(s *state.State, i *index.Index, token crypto.Token) *MyMediaView {
	head := HeaderInfo{
		Active:  "MyMedia",
		Path:    "venture > my media > ",
		EndPath: "drafts",
		Section: "venture",
	}
	myMedia := &MyMediaView{
		Drafts: make([]MyDraftView, 0),
		Edits:  make([]MyEditView, 0),
		Head:   head,
	}
	drafts := i.MemberToDraft[token]
	for _, draft := range drafts {
		authorship := ""
		if draft.Authors.CollectiveName() == "" {
			if len(draft.Authors.ListOfMembers()) > 1 {
				authorship = "as coauthor"
			} else {
				authorship = "as author"
			}
		} else {
			authorship = "on behalf of " + draft.Authors.CollectiveName() + " collective"
		}
		myDraftView := MyDraftView{
			Title:      draft.Title,
			AuthorType: authorship,
			Hash:       crypto.EncodeHash(draft.DraftHash),
			Pinned:     make([]CaptionLink, 0),
			Edit:       make([]EditOnDraftView, 0),
			Stamps:     make([]CaptionLink, 0),
		}
		consensusEpoch := draft.Authors.ConsensusEpoch(draft.Votes)
		if consensusEpoch > 0 {
			myDraftView.PublishedAt = s.TimeOfEpoch(consensusEpoch).Format(time.RFC822)
		}
		for _, board := range draft.Pinned {
			myDraftView.Pinned = append(myDraftView.Pinned, CaptionLink{
				Caption: board.Name,
				Link:    url.QueryEscape(board.Name),
			})
		}
		for _, edit := range draft.Edits {
			authors := edit.Authors
			if authors == nil {
				continue
			}
			consensusEpoch := authors.ConsensusEpoch(edit.Votes)
			editOnDraft := EditOnDraftView{
				Link: crypto.EncodeHash(edit.Edit),
				Time: s.TimeOfEpoch(consensusEpoch).Format(time.RFC822),
			}
			if authors.CollectiveName() != "" {
				editOnDraft.Caption = fmt.Sprintf("on behalf of %s", authors.CollectiveName())
			} else {
				for author, _ := range authors.ListOfMembers() {
					if handle, ok := s.Members[crypto.HashToken(author)]; ok {
						if N := len(authors.ListOfMembers()); N > 1 {
							editOnDraft.Caption = fmt.Sprintf("by %s and %d others", handle, N-1)
						} else {
							editOnDraft.Caption = fmt.Sprintf("by %s", handle)
						}
					}
				}
			}
			myDraftView.Edit = append(myDraftView.Edit, editOnDraft)
		}
		release := s.Releases[draft.DraftHash]
		if release != nil {
			for _, stamp := range release.Stamps {
				myDraftView.Stamps = append(myDraftView.Stamps, CaptionLink{
					Caption: stamp.Reputation.Name,
					Link:    url.QueryEscape(stamp.Reputation.Name),
				})
			}
			relaseEpoch := release.Draft.Authors.ConsensusEpoch(release.Votes)
			if relaseEpoch > 0 {
				myDraftView.Release = s.TimeOfEpoch(relaseEpoch).Format(time.RFC822)
			}
		}
		myMedia.Drafts = append(myMedia.Drafts, myDraftView)

	}
	edits := i.MemberToEdit[token]
	for _, edit := range edits {
		authorship := ""
		if edit.Authors.CollectiveName() == "" {
			if len(edit.Authors.ListOfMembers()) > 1 {
				authorship = "as coauthor"
			} else {
				authorship = "as author"
			}
		} else {
			authorship = "on behalf of " + edit.Authors.CollectiveName() + " collective"
		}
		myEditView := MyEditView{
			DraftTitle: edit.Draft.Title,
			DraftHash:  crypto.EncodeHash(edit.Draft.DraftHash),
			Hash:       crypto.EncodeHash(edit.Edit),
			AuthorType: authorship,
		}
		consensusEpoch := edit.Authors.ConsensusEpoch(edit.Votes)
		if consensusEpoch > 0 {
			myEditView.PublishedAt = s.TimeOfEpoch(consensusEpoch).Format(time.RFC822)
		}
		myMedia.Edits = append(myMedia.Edits, myEditView)
	}
	return myMedia
}

type NewActionView struct {
	Action   string
	Category string
	Duration string
}

type NewActionsView struct {
	NewStuff  int
	Updates   int
	Awareness int
	People    int
	Actions   []NewActionView
	Head      HeaderInfo
}

func NewActionsFromState(s *state.State, i *index.Index, genesisTime time.Time) *NewActionsView {
	head := HeaderInfo{
		Active:  "News",
		Path:    "explore > ",
		EndPath: "news",
		Section: "explore",
	}
	view := NewActionsView{
		Actions: make([]NewActionView, 0),
		Head:    head,
	}
	for _, action := range i.RecentActions {
		if action.Approved == 1 {
			des, category, epoch := i.ActionToFormatedString(action.Action)
			if len(des) > 0 {
				actionTime := genesisTime.Add(time.Second * time.Duration(epoch))
				duration := PrettyDuration(time.Since(actionTime))
				view.Actions = append(view.Actions, NewActionView{Action: fmt.Sprintf("<span>%v</span>", des), Category: strings.ReplaceAll(category, " ", "_"), Duration: duration})
				switch category {
				case "new stuff":
					view.NewStuff += 1
				case "updates":
					view.Updates += 1
				case "awareness":
					view.Awareness += 1
				case "people":
					view.People += 1
				}
			}
		}
	}
	return &view
}

type MyEventView struct {
	Collective           string
	StartAt              string
	Description          string
	Venue                string
	Open                 bool
	Public               bool
	Greeting             bool
	Attendee             []CaptionLink
	AttendeeCount        int
	GreetingCount        int
	GreetingPendingCount int
	Hash                 string
}

type MyEventsView struct {
	TodayCount    int
	NextWeekCount int
	FurtherCount  int
	Events        []MyEventView
	Managed       []MyEventView
	Head          HeaderInfo
}

func MyEventsFromState(s *state.State, i *index.Index, token crypto.Token) *MyEventsView {
	head := HeaderInfo{
		Active:  "MyEvents",
		Path:    "venture > my events > ",
		EndPath: "attending",
		Section: "venture",
	}
	view := MyEventsView{
		Events: make([]MyEventView, 0),
		Head:   head,
	}
	events := i.MemberToCheckin[token]
	managed := i.EventsOnMember(token)
	hashes := make(map[crypto.Hash]struct{})
	for _, event := range events {
		hashes[event.Hash] = struct{}{}
	}
	for _, hash := range managed {
		if _, ok := hashes[hash]; !ok {
			if event, ok := s.Events[hash]; ok {
				events = append(events, event)
			}
		}
	}
	for _, event := range events {
		if time.Until(event.StartAt) < -12*time.Hour {
			continue
		}
		if time.Until(event.StartAt) < 24*time.Hour {
			view.TodayCount += 1
		} else if time.Until(event.StartAt) < 7*24*time.Hour {
			view.NextWeekCount += 1
		} else {
			view.FurtherCount += 1
		}
		eventView := MyEventView{
			Collective:    event.Collective.Name,
			StartAt:       event.StartAt.Format(time.RFC822),
			Description:   event.Description,
			Venue:         event.Venue,
			Open:          event.Open,
			Public:        event.Public,
			Hash:          crypto.EncodeHash(event.Hash),
			AttendeeCount: len(event.Checkin),
		}
		if greeting, ok := event.Checkin[token]; ok {
			if greeting.Action != nil {
				eventView.Greeting = true
			}
			eventView.GreetingCount += 1
		}
		for _, checkin := range event.Checkin {
			if checkin != nil {
				token := checkin.Action.Author
				if handle, ok := s.Members[crypto.HashToken(token)]; ok {
					eventView.Attendee = append(eventView.Attendee, CaptionLink{
						Caption: handle,
						Link:    url.QueryEscape(handle),
					})
				}
			}
		}
		eventView.GreetingPendingCount = len(event.Checkin) - eventView.GreetingCount
		view.Events = append(view.Events, eventView)
	}
	for _, hash := range managed {
		if event, ok := s.Events[hash]; ok {
			eventView := MyEventView{
				Collective:  event.Collective.Name,
				StartAt:     event.StartAt.Format(time.RFC822),
				Description: event.Description,
				Venue:       event.Venue,
				Open:        event.Open,
				Public:      event.Public,
				Hash:        crypto.EncodeHash(event.Hash),
			}
			view.Managed = append(view.Managed, eventView)
		}
	}
	return &view
}

type DetailedVote struct {
	Author  CaptionLink
	Approve bool
	Reasons string
}

type DetailedPool struct {
	Approve  []DetailedVote
	Reject   []DetailedVote
	Needed   int
	NotVoted []CaptionLink
	Head     HeaderInfo
}

func DetailedVoteFromState(s *state.State, i *index.Index, hash crypto.Hash) *DetailedPool {
	pool := s.Proposals.Pooling(hash)
	if pool == nil {
		return nil
	}
	fmt.Println(pool.Voters)
	detailed := DetailedPool{
		Needed:   pool.Majority * len(pool.Voters) / 100,
		Approve:  make([]DetailedVote, 0),
		Reject:   make([]DetailedVote, 0),
		NotVoted: make([]CaptionLink, 0),
	}
	detailed.Head = HeaderInfo{
		Active:  "Pending",
		Path:    "venture >",
		EndPath: "create collective",
		Section: "venture",
	}
	for _, vote := range pool.Votes {
		author := s.Members[crypto.HashToken(vote.Author)]
		voteDetailed := DetailedVote{
			Author:  CaptionLink{Caption: author, Link: url.QueryEscape(author)},
			Approve: vote.Approve,
			Reasons: vote.Reasons,
		}
		if vote.Approve {
			detailed.Approve = append(detailed.Approve, voteDetailed)
		} else {
			detailed.Reject = append(detailed.Reject, voteDetailed)
		}
		delete(pool.Voters, vote.Author)
	}
	for voter, _ := range pool.Voters {
		author := s.Members[crypto.HashToken(voter)]
		detailed.NotVoted = append(detailed.NotVoted, CaptionLink{Caption: author, Link: url.QueryEscape(author)})
	}
	return &detailed
}
