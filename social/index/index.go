package index

import (
	"github.com/lienkolabs/swell/crypto"
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

type Index struct {
	//
	indexedMembers      map[crypto.Token]string           // token to handle
	memberToAction      map[crypto.Token][]*indexedAction // ação e se foi aprovada ou se está pendente
	pendingIndexActions map[crypto.Hash]crypto.Token

	// central connections member connections
	memberToCollective map[crypto.Token][]*state.Collective
	memberToBoard      map[crypto.Token][]*state.Board
	memberToEvent      map[crypto.Token][]*state.Event

	//memberToEdit       map[string][]*state.Edit
	//memberToDraft      map[s]

	// central connections collectives card
	collectiveToBoards map[*state.Collective][]*state.Board
	collectiveToStamps map[*state.Collective][]*state.Stamp
	collectiveToEvents map[*state.Collective][]*state.Event
	// collectiveLastAction map[*state.Collective][]lastaction

	// central connections edit card
	editToDrafts map[*state.Edit][]*state.Draft
}

func (i *Index) BoardsOnCollective(collective *state.Collective) []*state.Board {
	return i.collectiveToBoards[collective]
}

func (i *Index) IndexAction(action actions.Action) {
	author := action.Authored()
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
			newAction.approved = 1
		case *actions.RequestMembership:
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

// Collective's boards

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

// Collective's stamps

func (i *Index) AddStampToCollective(stamp *state.Stamp, collective *state.Collective) {
	if stamps, ok := i.collectiveToStamps[collective]; ok {
		i.collectiveToStamps[collective] = append(stamps, stamp)
	} else {
		i.collectiveToStamps[collective] = []*state.Stamp{stamp}
	}
}

// Collective's events

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

// Edit's drafts

func (i *Index) AddDraftToEdit(draft *state.Draft, edit *state.Edit) {
	if drafts, ok := i.editToDrafts[edit]; ok {
		i.editToDrafts[edit] = append(drafts, draft)
	} else {
		i.editToDrafts[edit] = []*state.Draft{draft}
	}
}

func NewIndex() *Index {
	return &Index{
		// central connections
		memberToCollective: make(map[crypto.Token][]*state.Collective),
		memberToBoard:      make(map[crypto.Token][]*state.Board),
		memberToEvent:      make(map[crypto.Token][]*state.Event),
		//memberToEdit:       make(map[string][]*state.Edit),
		collectiveToBoards: make(map[*state.Collective][]*state.Board),
		collectiveToStamps: make(map[*state.Collective][]*state.Stamp),
		collectiveToEvents: make(map[*state.Collective][]*state.Event),
		// collectiveLastAction: make(map[*state.Collective][]lastaction),
		editToDrafts: make(map[*state.Edit][]*state.Draft),

		indexedMembers:      make(map[crypto.Token]string),
		memberToAction:      make(map[crypto.Token][]*indexedAction),
		pendingIndexActions: make(map[crypto.Hash]crypto.Token),
	}
}
