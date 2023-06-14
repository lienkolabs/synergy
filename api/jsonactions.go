package api

import (
	"encoding/json"
	"time"

	"github.com/lienkolabs/swell/crypto"
	"github.com/lienkolabs/synergy/social/actions"
)

/*
	actions
		AcceptCheckinEvent
		BoardEditor
		CancelEvent
		CheckinEvent
		CreateBoard
		CreateCollective
		CreateEvent
		Draft
		Edit
		ImprintStamp
		Pin
		React
		ReleaseDraft
		RemoveMember
		RequestMembership
		UpdateBoard
		UpdateCollective
		UpdateEvent
		Vote
*/

type Policy struct {
	Majority      int `json:"majority"`
	SuperMajority int `json:"superMajority"`
}

type Action struct {
	Action string `json:"action"`
	ID     int    `json:"id"`
}

func JSONType(data []byte) string {
	var a Action
	json.Unmarshal(data, &a)
	return a.Action
}

type AcceptCheckinEvent struct {
	Action    string       `json:"action"`
	ID        int          `json:"id"`
	Reasons   string       `json:"reasons"`
	EventHash crypto.Hash  `json:"eventHash"`
	CheckedIn crypto.Token `json:"checkedIn"`
}

func (a AcceptCheckinEvent) ToAction() actions.Action {
	action := actions.AcceptCheckinEvent{
		Reasons:   a.Reasons,
		EventHash: a.EventHash,
		CheckedIn: a.CheckedIn,
	}
	return &action
}

type BoardEditor struct {
	Action  string       `json:"action"`
	ID      int          `json:"id"`
	Reasons string       `json:"reasons"`
	Board   string       `json:"board"`
	Editor  crypto.Token `json:"editor"`
	Insert  bool         `json:"insert"`
}

func (a BoardEditor) ToAction() actions.Action {
	action := actions.BoardEditor{
		Reasons: a.Reasons,
		Board:   a.Board,
		Editor:  a.Editor,
		Insert:  a.Insert,
	}
	return &action
}

type CancelEvent struct {
	Action  string      `json:"action"`
	ID      int         `json:"id"`
	Reasons string      `json:"reasons"`
	Hash    crypto.Hash `json:"hash"`
}

func (a CancelEvent) ToAction() actions.Action {
	action := actions.CancelEvent{
		Reasons: a.Reasons,
		Hash:    a.Hash,
	}
	return &action
}

type CheckinEvent struct {
	Action    string      `json:"action"`
	ID        int         `json:"id"`
	Reasons   string      `json:"reasons"`
	EventHash crypto.Hash `json:"eventHash"`
}

func (a CheckinEvent) ToAction() actions.Action {
	action := actions.CheckinEvent{
		Reasons:   a.Reasons,
		EventHash: a.EventHash,
	}
	return &action
}

type CreateBoard struct {
	Action      string   `json:"action"`
	ID          int      `json:"id"`
	Reasons     string   `json:"reasons"`
	OnBehalfOf  string   `json:"onBehalfOf"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Keywords    []string `json:"keywords"`
	PinMajority int      `json:"pinMajority"`
}

func (a CreateBoard) ToAction() actions.Action {
	action := actions.CreateBoard{
		Reasons:     a.Reasons,
		OnBehalfOf:  a.OnBehalfOf,
		Name:        a.Name,
		Description: a.Description,
		Keywords:    a.Keywords,
		PinMajority: byte(a.PinMajority),
	}
	return &action
}

type CreateCollective struct {
	Action      string `json:"action"`
	ID          int    `json:"id"`
	Reasons     string `json:"reasons"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Policy      Policy `json:"policy"`
}

func (a CreateCollective) ToAction() actions.Action {
	action := actions.CreateCollective{
		Reasons:     a.Reasons,
		Name:        a.Name,
		Description: a.Description,
		Policy:      actions.Policy(a.Policy),
	}
	return &action
}

type CreateEvent struct {
	Action       string         `json:"action"`
	ID           int            `json:"id"`
	Reasons      string         `json:"reasons"`
	OnBehalfOf   string         `json:"onBehalfOf"`
	StartAt      time.Time      `json:"startAt"`
	EstimatedEnd time.Time      `json:"estimatedEnd"`
	Description  string         `json:"description"`
	Venue        string         `json:"venue"`
	Open         bool           `json:"open"`
	Public       bool           `json:"public"`
	Managers     []crypto.Token `json:"managers,omitempty"`
}

func (a CreateEvent) ToAction() actions.Action {
	action := actions.CreateEvent{
		Reasons:      a.Reasons,
		OnBehalfOf:   a.OnBehalfOf,
		StartAt:      a.StartAt,
		EstimatedEnd: a.EstimatedEnd,
		Description:  a.Description,
		Venue:        a.Venue,
		Open:         a.Open,
		Public:       a.Public,
		Managers:     a.Managers,
	}
	return &action
}

type Draft struct {
	Action        string         `json:"action"`
	ID            int            `json:"id"`
	Reasons       string         `json:"reasons"`
	OnBehalfOf    string         `json:"onBeahlfOf,omitempty"`
	CoAuthors     []crypto.Token `json:"coAuthors,omitempty"`
	Policy        *Policy        `json:"policy"`
	Title         string         `json:"title"`
	Keywords      []string       `json:"keywords"`
	Description   string         `json:"description"`
	ContentType   string         `json:"contentType"`
	FilePath      string         `json:"filePath"`
	PreviousDraft crypto.Hash    `json:"previousDraft,,omitempty"`
	References    []string       `json:"references,omitempty"`
}

type Edit struct {
	Action      string         `json:"action"`
	ID          int            `json:"id"`
	Reasons     string         `json:"reasons"`
	OnBehalfOf  string         `json:"onBeahlfOf,omitempty"`
	CoAuthors   []crypto.Token `json:"coAuthors,omitempty"`
	EditedDraft crypto.Hash    `json:"editedDraft"`
	ContentType string         `json:"contentType"`
	FilePath    string         `json:"filePath"`
}

type ImprintStamp struct {
	Action     string      `json:"action"`
	ID         int         `json:"id"`
	Reasons    string      `json:"reasons"`
	OnBehalfOf string      `json:"onBeahlfOf,omitempty"`
	Hash       crypto.Hash `json:"hash"`
}

func (a ImprintStamp) ToAction() actions.Action {
	action := actions.ImprintStamp{
		Reasons:    a.Reasons,
		OnBehalfOf: a.OnBehalfOf,
		Hash:       a.Hash,
	}
	return &action
}

type Pin struct {
	Action  string      `json:"action"`
	ID      int         `json:"id"`
	Reasons string      `json:"reasons"`
	Board   string      `json:"board"`
	Draft   crypto.Hash `json:"draft"`
	Pin     bool        `json:"pin"`
}

func (a Pin) ToAction() actions.Action {
	action := actions.Pin{
		Reasons: a.Reasons,
		Board:   a.Board,
		Pin:     a.Pin,
	}
	return &action
}

type React struct {
	Action     string      `json:"action"`
	ID         int         `json:"id"`
	Reasons    string      `json:"reasons"`
	OnBehalfOf string      `json:"onBeahlfOf,omitempty"`
	Hash       crypto.Hash `json:"hash"`
	Reaction   byte        `json:"reaction"`
}

func (a React) ToAction() actions.Action {
	action := actions.React{
		Reasons:    a.Reasons,
		OnBehalfOf: a.OnBehalfOf,
		Hash:       a.Hash,
		Reaction:   a.Reaction,
	}
	return &action
}

type ReleaseDraft struct {
	Action      string      `json:"action"`
	ID          int         `json:"id"`
	Reasons     string      `json:"reasons"`
	ContentHash crypto.Hash `json:"contentHash"`
}

func (a ReleaseDraft) ToAction() actions.Action {
	action := actions.ReleaseDraft{
		Reasons:     a.Reasons,
		ContentHash: a.ContentHash,
	}
	return &action
}

type RemoveMember struct {
	Action     string       `json:"action"`
	ID         int          `json:"id"`
	Reasons    string       `json:"reasons"`
	OnBehalfOf string       `json:"onBeahlfOf,omitempty"`
	Member     crypto.Token `json:"member"`
}

func (a RemoveMember) ToAction() actions.Action {
	action := actions.RemoveMember{
		Reasons:    a.Reasons,
		OnBehalfOf: a.OnBehalfOf,
		Member:     a.Member,
	}
	return &action
}

type RequestMembership struct {
	Action     string `json:"action"`
	ID         int    `json:"id"`
	Reasons    string `json:"reasons"`
	Collective string `json:"collective"`
	Include    bool   `json:"include"`
}

func (a RequestMembership) ToAction() actions.Action {
	action := actions.RequestMembership{
		Reasons:    a.Reasons,
		Collective: a.Collective,
		Include:    a.Include,
	}
	return &action
}

type UpdateBoard struct {
	Action      string    `json:"action"`
	ID          int       `json:"id"`
	Reasons     string    `json:"reasons"`
	Board       string    `json:"board"`
	Description *string   `json:"description,omitempty"`
	Keywords    *[]string `json:"keywords,omitempty"`
	PinMajority *int      `json:"pinMajority"`
}

func (a UpdateBoard) ToAction() actions.Action {
	action := actions.UpdateBoard{
		Reasons:     a.Reasons,
		Board:       a.Board,
		Description: *a.Description,
		Keywords:    *a.Keywords,
		PinMajority: byte(*a.PinMajority),
	}
	return &action
}

type UpdateCollective struct {
	Action      string  `json:"action"`
	ID          int     `json:"id"`
	Reasons     string  `json:"reasons"`
	OnBehalfOf  string  `json:"onBehalfOf"`
	Description *string `json:"description,omitempty"`
	Policy      *Policy `json:"policy,omitempty"`
}

func (a UpdateCollective) ToAction() actions.Action {
	action := actions.UpdateCollective{
		Reasons:     a.Reasons,
		OnBehalfOf:  a.OnBehalfOf,
		Description: *a.Description,
	}
	if a.Policy != nil {
		action.Policy = &actions.Policy{
			Majority:      a.Policy.Majority,
			SuperMajority: a.Policy.SuperMajority,
		}
	}
	return &action
}

type UpdateEvent struct {
	Action      string `json:"action"`
	ID          int    `json:"id"`
	Reasons     string `json:"reasons"`
	EventHash   crypto.Hash
	Description *string         `json:"description,omitempty"`
	Venue       *string         `json:"venue,omitempty"`
	Open        *bool           `json:"open,omitempty"`
	Public      *bool           `json:"public,omitempty"`
	Managers    *[]crypto.Token `json:"managers,omitempty"`
}

type Vote struct {
	Action  string `json:"action"`
	ID      int    `json:"id"`
	Reasons string `json:"reasons,omitempty"`
	Hash    string `json:"hash"`
	Approve bool   `json:"approve"`
}

type DraftObject struct {
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	OnBehalfOf      string   `json:"onBehalfOf,omitempty"`
	Authors         []string `json:"authors,omitempty"`
	DraftType       string   `json:"draftType"`
	DraftHash       string   `json:"draftHash"`
	PreviousVersion string   `json:"previousVersion.omitempty"`
}
