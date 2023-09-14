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
	}
	sort.Sort(updatesView)
	return updatesView
}

type PendingActionsView struct {
	Pending []PendingActionDetailView
}

type PendingActionDetailView struct {
	Description string
	ProposedAt  string
	NVotes      int
}

func PendingActionsFromState(s *state.State, i *index.Index, token crypto.Token, genesisTime time.Time) *PendingActionsView {
	pending := i.GetPendingActions(token)
	if len(pending) == 0 {
		return nil
	}
	view := PendingActionsView{
		Pending: make([]PendingActionDetailView, 0),
	}
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
