package api

import (
	"fmt"
	"sort"
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

func actionsToActionUpdateView(actions []index.ActionDetails, genesisTime time.Time) []ActionUpdateView {
	updates := make([]ActionUpdateView, 0)
	for _, action := range actions {
		actionTime := genesisTime.Add(time.Second * time.Duration(action.Epoch))
		actionUpdateView := ActionUpdateView{
			Description:         action.Description,
			LastUpdatedTime:     actionTime,
			LastUpdatedInterval: fmt.Sprintf("%s ago", time.Since(actionTime)),
		}

		if action.VoteStatus {
			actionUpdateView.VoteStatus = "approved"
		}
		if len(action.Votes) > 0 {
			actionUpdateView.VoteHash = crypto.EncodeHash(action.Votes[0].Hash)
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
		actions := i.GetRecentActions(crypto.Hasher([]byte(collective)))
		objView := ObjectUpdateView{
			Name:       collective,
			ObjectKind: "collective",
			Updates:    actionsToActionUpdateView(actions, genesisTime),
		}
		objects = append(objects, objView)
	}
	for _, board := range boards {
		actions := i.GetRecentActions(crypto.Hasher([]byte(board)))
		objView := ObjectUpdateView{
			Name:       board,
			ObjectKind: "board",
			Updates:    actionsToActionUpdateView(actions, genesisTime),
		}
		objects = append(objects, objView)
	}
	for _, eventHash := range events {
		if event, ok := s.Events[eventHash]; ok {
			actions := i.GetRecentActions(eventHash)
			objView := ObjectUpdateView{
				Name:       fmt.Sprintf("%s event from %s", event.StartAt.Format(time.RFC822), event.Collective.Name),
				ObjectKind: "event",
				Updates:    actionsToActionUpdateView(actions, genesisTime),
			}
			objects = append(objects, objView)
		}
	}
	if len(objects) == 0 {
		return nil
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
	Description string
	ProposedAt  string
	NVotes      int
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
	pending := i.GetPendingActions(token)
	// if len(pending) == 0 {
	// 	return nil
	// }
	for _, a := range pending {
		item := PendingActionDetailView{
			Description: a.Description,
			ProposedAt:  s.TimeOfEpoch(a.Epoch).Format(time.RFC822),
			NVotes:      len(a.Votes),
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
				Link:    fmt.Sprintf("/board/%s", board.Name),
			})
		}
		for _, edit := range draft.Edits {
			authors := edit.Authors
			if authors == nil {
				continue
			}
			consensusEpoch := authors.ConsensusEpoch(edit.Votes)
			editOnDraft := EditOnDraftView{
				Link: fmt.Sprintf("/edit/%s", edit.Edit),
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
					Link:    fmt.Sprintf("/collective/%s", stamp.Reputation.Name),
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
}

type NewActionsView struct {
	Actions []NewActionView
	Head    HeaderInfo
}

func NewActionsFromState(s *state.State, i *index.Index) *NewActionsView {
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
		if action.Approved != 2 {
			des, _, _, _, category := i.ActionToString(action.Action, action.Approved == 1)
			view.Actions = append(view.Actions, NewActionView{Action: des, Category: category})
		}
	}
	return &view
}

type MyEventView struct {
	Collective  string
	StartAt     string
	Description string
	Venue       string
	Open        bool
	Public      bool
	Greeting    bool
	Hash        string
}

type MyEventsView struct {
	Events  []MyEventView
	Managed []MyEventView
	Head    HeaderInfo
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
	for _, event := range events {
		eventView := MyEventView{
			Collective:  event.Collective.Name,
			StartAt:     event.StartAt.Format(time.RFC822),
			Description: event.Description,
			Venue:       event.Venue,
			Open:        event.Open,
			Public:      event.Public,
			Hash:        crypto.EncodeHash(event.Hash),
		}
		if greeting, ok := event.Checkin[token]; ok {
			if greeting.Action != nil {
				eventView.Greeting = true
			}
		}
		view.Events = append(view.Events, eventView)
	}
	managed := i.EventsOnMember(token)
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
