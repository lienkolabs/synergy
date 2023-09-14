package index

import (
	"fmt"

	"github.com/lienkolabs/breeze/crypto"
	"github.com/lienkolabs/synergy/social/actions"
	"github.com/lienkolabs/synergy/social/state"
)

type LastAction struct {
	Author      crypto.Token
	Description string
	Epoch       uint64
}

type indexedAction struct {
	action   actions.Action
	hash     crypto.Hash
	approved byte
}

const ActionsCacheCount = 10

type ActionDetails struct {
	Description string
	ObjectHash  string
	Author      crypto.Token
	Votes       []actions.Vote
	VoteStatus  bool
	Epoch       uint64
}

type RecentActions struct {
	actions []actions.Action
}

func NewRecentActions(action actions.Action) *RecentActions {
	return &RecentActions{actions: []actions.Action{action}}
}

func (r *RecentActions) Append(action actions.Action) {
	if len(r.actions) == ActionsCacheCount {
		r.actions = append(r.actions[1:], action)
	} else {
		r.actions = append(r.actions, action)
	}
}

func (r *RecentActions) Last() actions.Action {
	if len(r.actions) == 0 {
		return nil
	}
	return r.actions[len(r.actions)-1]
}

type Index struct {
	//
	indexedMembers      map[crypto.Token]string           // token to handle
	memberToAction      map[crypto.Token][]*indexedAction // ação e se foi aprovada ou se está pendente
	pendingIndexActions map[crypto.Hash]crypto.Token

	indexVotes          map[crypto.Token]*SetOfHashes
	indexCompletedVotes map[crypto.Hash][]actions.Vote

	// central connections member connections
	memberToCollective map[crypto.Token][]string
	memberToBoard      map[crypto.Token][]string
	memberToEvent      map[crypto.Token][]crypto.Hash
	//memberToEdit
	//memberToDraft

	objectHashToActionHash map[crypto.Hash]*RecentActions // object to recent actions

	// central connections collectives card
	collectiveToBoards map[*state.Collective][]*state.Board
	collectiveToStamps map[*state.Collective][]*state.Stamp
	collectiveToEvents map[*state.Collective][]*state.Event
	// collectiveLastAction map[*state.Collective][]lastaction

	// central connections edit card
	editToDrafts map[*state.Edit][]*state.Draft

	state *state.State

	stateProposals *state.Proposals
}

func NewIndex() *Index {
	return &Index{
		// central connections
		memberToCollective: make(map[crypto.Token][]string),
		memberToBoard:      make(map[crypto.Token][]string),
		memberToEvent:      make(map[crypto.Token][]crypto.Hash),
		//memberToEdit:       make(map[string][]*state.Edit),
		collectiveToBoards: make(map[*state.Collective][]*state.Board),
		collectiveToStamps: make(map[*state.Collective][]*state.Stamp),
		collectiveToEvents: make(map[*state.Collective][]*state.Event),
		// collectiveLastAction: make(map[*state.Collective][]lastaction),
		editToDrafts: make(map[*state.Edit][]*state.Draft),

		indexedMembers:      make(map[crypto.Token]string),
		memberToAction:      make(map[crypto.Token][]*indexedAction),
		pendingIndexActions: make(map[crypto.Hash]crypto.Token),

		indexVotes:          make(map[crypto.Token]*SetOfHashes),
		indexCompletedVotes: make(map[crypto.Hash][]actions.Vote),

		objectHashToActionHash: make(map[crypto.Hash]*RecentActions),
	}
}

func (i *Index) SetState(s *state.State) {
	i.state = s
	i.stateProposals = s.Proposals
}

func (i *Index) ActionStatus(action actions.Action) ([]actions.Vote, bool) {
	hash := action.Hashed()
	if votes, ok := i.indexCompletedVotes[hash]; ok {
		return votes, true
	}
	return i.stateProposals.Votes(hash), false
}

// Objects related to a given collective

func (i *Index) BoardsOnCollective(collective *state.Collective) []*state.Board {
	return i.collectiveToBoards[collective]
}

func (i *Index) StampsOnCollective(collective *state.Collective) []*state.Stamp {
	return i.collectiveToStamps[collective]
}

func (i *Index) EventsOnCollective(collective *state.Collective) []*state.Event {
	return i.collectiveToEvents[collective]
}

// Objects related to a given member

func (i *Index) CollectivesOnMember(member crypto.Token) []string {
	return i.memberToCollective[member]
}

func (i *Index) BoardsOnMember(member crypto.Token) []string {
	return i.memberToBoard[member]
}

func (i *Index) EventsOnMember(member crypto.Token) []crypto.Hash {
	return i.memberToEvent[member]
}

func (i *Index) AddMemberToIndex(token crypto.Token, handle string) {
	i.indexedMembers[token] = handle
}

func (i Index) GetLastAction(objectHash crypto.Hash) *ActionDetails {
	recent := i.objectHashToActionHash[objectHash]
	if recent == nil || len(recent.actions) == 0 {
		return nil
	}
	// TODO: check consensus status
	des, hash, author, epoch := i.ActionToString(recent.actions[len(recent.actions)-1], true)
	return &ActionDetails{
		Description: des,
		Author:      author,
		ObjectHash:  hash,
		Epoch:       epoch,
	}
}

func (i Index) GetRecentActions(objectHash crypto.Hash) []ActionDetails {
	recent := i.objectHashToActionHash[objectHash]
	if recent == nil {
		return nil
	}
	details := make([]ActionDetails, len(recent.actions))
	for n, r := range recent.actions {
		// TODO: check consensus status
		status := true
		des, _, _, epoch := i.ActionToString(r, status)
		votes, status := i.ActionStatus(r)
		details[n] = ActionDetails{
			Description: des,
			Votes:       votes,
			VoteStatus:  status,
			Epoch:       epoch,
		}
	}
	return details
}

func (i *Index) IndexAction(action actions.Action) {
	author := action.Authored()
	objects := i.ActionToObjects(action)
	for _, object := range objects {
		if recent, ok := i.objectHashToActionHash[object]; ok {
			recent.Append(action)
		} else {
			i.objectHashToActionHash[object] = NewRecentActions(action)
		}
	}
	//hash := action.Hashed()
	if _, ok := i.indexedMembers[author]; ok {
		newAction := indexedAction{
			action:   action,
			hash:     action.Hashed(),
			approved: 0,
		}
		switch action.(type) {
		case *actions.GreetCheckinEvent:
			newAction.approved = 1
		case *actions.CheckinEvent:
			newAction.approved = 1
		case *actions.React:
			newAction.approved = 1
		case *actions.Signin:
			newAction.approved = 1
		case *actions.Vote:
			newAction.approved = 1
		case *actions.CreateCollective:
			i.IndexConsensusAction(action)
			newAction.approved = 1
		}
		if indexedActions, ok := i.memberToAction[author]; ok {
			i.memberToAction[author] = append(indexedActions, &newAction)
		} else {
			i.memberToAction[author] = []*indexedAction{&newAction}
		}
		if newAction.approved == 0 {
			hash := action.Hashed()
			i.pendingIndexActions[hash] = author
		}
	}
}

func (i *Index) isIndexedMember(token crypto.Token) bool {
	_, ok := i.indexedMembers[token]
	return ok
}

func appendOrCreate[T any](values []T, value T) []T {
	if values == nil {
		return []T{value}
	}
	return append(values, value)
}

func removeItem[T comparable](values []T, value T) []T {
	for n, item := range values {
		if item == value {
			if n == len(values)-1 {
				return values[0:n]
			}
			return append(values[0:n], values[n+1:]...)
		}
	}
	return values
}

func (i *Index) IndexConsensusAction(action actions.Action) {
	switch v := action.(type) {
	case *actions.CreateCollective:
		fmt.Println("index", v)
		if i.isIndexedMember(v.Author) {
			i.memberToCollective[v.Author] = appendOrCreate[string](i.memberToCollective[v.Author], v.Name)
		}
	case *actions.RequestMembership:
		if i.isIndexedMember(v.Author) {
			i.memberToCollective[v.Author] = appendOrCreate[string](i.memberToCollective[v.Author], v.Collective)
		}
	case *actions.RemoveMember:
		if i.isIndexedMember(v.Member) {
			i.memberToCollective[v.Member] = removeItem[string](i.memberToCollective[v.Author], v.OnBehalfOf)
		}
	case *actions.CreateBoard:
		if i.isIndexedMember(v.Author) {
			i.memberToBoard[v.Author] = appendOrCreate[string](i.memberToBoard[v.Author], v.Name)
		}
	case *actions.CreateEvent:
		if i.isIndexedMember(v.Author) {
			hash := crypto.Hasher(v.Serialize())
			i.memberToEvent[v.Author] = appendOrCreate[crypto.Hash](i.memberToEvent[v.Author], hash)
		}
	}
}

func (i *Index) IndexConsensus(hash crypto.Hash, approved bool) {
	author, ok := i.pendingIndexActions[hash]
	if !ok {
		return
	}
	delete(i.pendingIndexActions, hash)
	indexActions, ok := i.memberToAction[author]
	if !ok {
		return
	}
	for _, action := range indexActions {
		if action.hash.Equal(hash) {
			if approved {
				i.IndexConsensusAction(action.action)
				action.approved = 1
			} else {
				action.approved = 2
			}
		}
	}
}

func (i *Index) LastMemberActionOnCollective(member crypto.Token, collective string) *LastAction {
	allActions, ok := i.memberToAction[member]
	if !ok {
		return nil
	}
	for n := len(allActions) - 1; n >= 0; n-- {
		switch v := allActions[n].action.(type) {
		case *actions.CreateBoard:
			if v.OnBehalfOf == collective {
				return &LastAction{
					Author:      member,
					Description: "create board",
					Epoch:       v.Epoch,
				}
			}
		case *actions.Draft:
			if v.OnBehalfOf == collective && v.Authored().Equal(member) {
				return &LastAction{
					Author:      member,
					Description: "submit draft",
					Epoch:       v.Epoch,
				}
			}
		case *actions.Edit:
			if v.OnBehalfOf == collective && v.Authored().Equal(member) {
				return &LastAction{
					Author:      member,
					Description: "submit edit",
					Epoch:       v.Epoch,
				}
			}

		case *actions.CreateEvent:
			if v.OnBehalfOf == collective && v.Authored().Equal(member) {
				return &LastAction{
					Author:      member,
					Description: "create event",
					Epoch:       v.Epoch,
				}
			}
		case *actions.RemoveMember:
			if v.OnBehalfOf == collective && v.Authored().Equal(member) {
				return &LastAction{
					Author:      member,
					Description: "remove member",
					Epoch:       v.Epoch,
				}
			}
		}
	}
	return nil
}

func (i *Index) AddBoardToCollective(board *state.Board, collective *state.Collective) {
	if boards, ok := i.collectiveToBoards[collective]; ok {
		i.collectiveToBoards[collective] = append(boards, board)
	} else {
		i.collectiveToBoards[collective] = []*state.Board{board}
	}
}

func (i *Index) RemoveBoardFromCollective(board *state.Board, collective *state.Collective) {
	if boards, ok := i.collectiveToBoards[collective]; ok {
		for n, e := range boards {
			if e == board {
				removed := boards[0:n]
				if n < len(boards)-1 {
					removed = append(removed, boards[n+1:]...)
				}
				i.collectiveToBoards[collective] = removed
			}
		}
	}
}

func (i *Index) AddStampToCollective(stamp *state.Stamp, collective *state.Collective) {
	if stamps, ok := i.collectiveToStamps[collective]; ok {
		i.collectiveToStamps[collective] = append(stamps, stamp)
	} else {
		i.collectiveToStamps[collective] = []*state.Stamp{stamp}
	}
}

func (i *Index) AddEventToCollective(event *state.Event, collective *state.Collective) {
	if events, ok := i.collectiveToEvents[collective]; ok {
		i.collectiveToEvents[collective] = append(events, event)
	} else {
		i.collectiveToEvents[collective] = []*state.Event{event}
	}
}

func (i *Index) RemoveEventFromCollective(event *state.Event, collective *state.Collective) {
	if events, ok := i.collectiveToEvents[collective]; ok {
		for n, e := range events {
			if e == event {
				removed := events[0:n]
				if n < len(events)-1 {
					removed = append(removed, events[n+1:]...)
				}
				i.collectiveToEvents[collective] = removed
			}
		}
	}
}

/*
func (i *Index) RemoveMemberFromCollective(collective *state.Collective, member crypto.Token) {
	delete(i.memberToCollective, member)
}

func (i *Index) AddMemberToCollective(collective *state.Collective, member crypto.Token) {
	if collectives, ok := i.memberToCollective[member]; ok {
		i.memberToCollective[member] = append(collectives, collective)
	} else {
		i.memberToCollective[member] = []*state.Collective{collective}
	}
}

func (i *Index) AddEditorToBoard(board *state.Board, editor crypto.Token) {
	if boards, ok := i.memberToBoard[editor]; ok {
		i.memberToBoard[editor] = append(boards, board)
	} else {
		i.memberToBoard[editor] = []*state.Board{board}
	}
}

func (i *Index) RemoveEditorFromBoard(board *state.Board, editor crypto.Token) {
	delete(i.memberToBoard, editor)
}
*/

// Collective's boards

/*func (i *Index) AddBoardToCollective(board *state.Board, collective *state.Collective) {if boards, ok := i.collectiveToBoards[collective]; ok {
		i.collectiveToBoards[collective] = append(boards, board)
	} else {
		i.collectiveToBoards[collective] = []*state.Board{board}
	}
}

func (i *Index) RemoveBoardFromCollective(board *state.Board, collective *state.Collective) {
	if boards, ok := i.collectiveToBoards[collective]; ok {
		for n, e := range boards {
			if e == board {
				removed := boards[0:n]
				if n < len(boards)-1 {
					removed = append(removed, boards[n+1:]...)
				}
				i.collectiveToBoards[collective] = removed
			}
		}
	}
}

// Collective's stamps


// Collective's events


// Edit's drafts

func (i *Index) AddDraftToEdit(draft *state.Draft, edit *state.Edit) {
	if drafts, ok := i.editToDrafts[edit]; ok {
		i.editToDrafts[edit] = append(drafts, draft)
	} else {
		i.editToDrafts[edit] = []*state.Draft{draft}
	}
}
*/
