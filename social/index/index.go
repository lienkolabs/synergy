package index

import (
	"github.com/lienkolabs/synergy/social/state"
)

type lastaction struct {
	Author     string
	ActionName string
	TimeWindow string
}

type membercollective struct {
	Collective *state.Collective
	LastAction lastaction
}

type memberboard struct {
	Board      *state.Board
	LastAction lastaction
}

type memberevent struct {
	Event      *state.Event
	LastAction lastaction
}

type Index struct {
	// central connections member connections
	memberToCollective map[string][]membercollective
	memberToBoard      map[string][]memberboard
	memberToEvent      map[string][]memberevent
	memberToEdit       map[string][]*state.Edit

	// central connections collectives card
	collectiveToBoards map[*state.Collective][]*state.Board
	collectiveToStamps map[*state.Collective][]*state.Stamp
	collectiveToEvents map[*state.Collective][]*state.Event
	// collectiveLastAction map[*state.Collective][]lastaction

	// central connections edit card
	editToDrafts map[*state.Edit][]*state.Draft
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
		memberToCollective: make(map[string][]membercollective),
		memberToBoard:      make(map[string][]memberboard),
		memberToEvent:      make(map[string][]memberevent),
		memberToEdit:       make(map[string][]*state.Edit),
		collectiveToBoards: make(map[*state.Collective][]*state.Board),
		collectiveToStamps: make(map[*state.Collective][]*state.Stamp),
		collectiveToEvents: make(map[*state.Collective][]*state.Event),
		// collectiveLastAction: make(map[*state.Collective][]lastaction),
		editToDrafts: make(map[*state.Edit][]*state.Draft),
	}
}
