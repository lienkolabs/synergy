package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/lienkolabs/swell/crypto"
	"github.com/lienkolabs/synergy/social/actions"
)

type Handles struct {
	Handle   map[string]crypto.Token
	Attorney *Attorney
}

func (h *Handles) ApiHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	var actionArray []actions.Action
	var err error
	switch r.FormValue("action") {
	case "AcceptCheckinEvent":
		actionArray, err = AcceptCheckinEventForm(r, h.Handle).ToAction()
	case "BoardEditor":
		actionArray, err = BoardEditorForm(r, h.Handle).ToAction()
	case "CancelEvent":
		actionArray, err = CancelEventForm(r).ToAction()
	case "CheckinEvent":
		actionArray, err = CheckinEventForm(r).ToAction()
	case "CreateBorad":
		actionArray, err = CreateBoradForm(r).ToAction()
	case "CreateCollective":
		actionArray, err = CreateCollectiveForm(r).ToAction()
	case "CreateEvent":
		actionArray, err = CreateEventForm(r, h.Handle).ToAction()
	case "Draft":
		actionArray, err = DraftForm(r, h.Handle).ToAction()
	case "Edit":
		actionArray, err = EditForm(r, h.Handle).ToAction()
	case "ImprintStamp":
		actionArray, err = ImprintStampForm(r).ToAction()
	case "Pin":
		actionArray, err = PinForm(r).ToAction()
	case "React":
		actionArray, err = ReactForm(r).ToAction()
	case "ReleaseDraft":
		actionArray, err = ReleaseDraftForm(r).ToAction()
	case "RemoveMember":
		actionArray, err = RemoveMemberForm(r, h.Handle).ToAction()
	case "RequestMembership":
		actionArray, err = RequestMembershipForm(r).ToAction()
	case "UpdateBoard":
		actionArray, err = UpdateBoardForm(r).ToAction()
	case "UpdateCollective":
		actionArray, err = UpdateCollectiveForm(r).ToAction()
	case "UpdateEvent":
		actionArray, err = UpdateEventForm(r, h.Handle).ToAction()
	case "Vote":
		actionArray, err = VoteForm(r).ToAction()
	}
	if err != nil && len(actionArray) > 0 {
		h.Attorney.Send(actionArray)
	}
	http.Redirect(w, r, "/static/index.html", http.StatusSeeOther)
}

func FormToI(r *http.Request, field string) int {
	value, _ := strconv.Atoi(r.FormValue(field))
	return value
}

func FormToHash(r *http.Request, field string) crypto.Hash {
	var hash crypto.Hash
	hash.UnmarshalText([]byte(r.FormValue(field)))
	return hash
}

func FormToToken(r *http.Request, field string, handles map[string]crypto.Token) crypto.Token {
	token := handles[r.FormValue(field)]
	return token
}

func FormToTokenArray(r *http.Request, field string, handles map[string]crypto.Token) []crypto.Token {
	h := strings.Split(r.FormValue(field), ",")
	tokens := make([]crypto.Token, 0)
	for _, handle := range h {
		if token, ok := handles[handle]; ok {
			tokens = append(tokens, token)
		}
	}
	return tokens
}

func FormToHashArray(r *http.Request, field string) []crypto.Hash {
	h := strings.Split(r.FormValue(field), ",")
	hashes := make([]crypto.Hash, 0)
	for _, caption := range h {
		var hash crypto.Hash
		if hash.UnmarshalText([]byte(caption)) == nil {
			hashes = append(hashes, hash)
		}
	}
	return hashes
}

func FormToStringArray(r *http.Request, field string) []string {
	words := strings.Split(r.FormValue(field), ",")
	return words
}

func FormToBool(r *http.Request, field string) bool {
	return r.FormValue(field) == "true"
}

func FormToPolicy(r *http.Request) Policy {
	return Policy{
		Majority:      FormToI(r, "policyMajority"),
		SuperMajority: FormToI(r, "policySupermajority"),
	}
}

func FormToTime(r *http.Request, field string) time.Time {
	t, _ := time.Parse(time.DateTime, r.FormValue(field))
	return t
}

func AcceptCheckinEventForm(r *http.Request, handles map[string]crypto.Token) AcceptCheckinEvent {
	action := AcceptCheckinEvent{
		Action:    "AcceptCheckinEvent",
		ID:        FormToI(r, "id"),
		Reasons:   r.FormValue("reasons"),
		EventHash: FormToHash(r, "eventHash"),
		CheckedIn: FormToToken(r, "checkedIn", handles),
	}
	return action
}

func BoardEditorForm(r *http.Request, handles map[string]crypto.Token) BoardEditor {
	action := BoardEditor{
		Action:  "BoardEditor",
		ID:      FormToI(r, "id"),
		Reasons: r.FormValue("reasons"),
		Board:   r.FormValue("board"),
		Insert:  FormToBool(r, "insert"),
	}
	return action
}

func CancelEventForm(r *http.Request) CancelEvent {
	action := CancelEvent{
		Action:  "CancelEvent",
		Reasons: r.FormValue("reasons"),
		ID:      FormToI(r, "id"),
		Hash:    FormToHash(r, "hash"),
	}
	return action
}

func CheckinEventForm(r *http.Request) CheckinEvent {
	action := CheckinEvent{
		Action:    "CheckinEvent",
		ID:        FormToI(r, "id"),
		Reasons:   r.FormValue("reasons"),
		EventHash: FormToHash(r, "eventhash"),
	}
	return action
}

func CreateBoradForm(r *http.Request) CreateBoard {
	action := CreateBoard{
		Action:      "CreateBoard",
		ID:          FormToI(r, "id"),
		Reasons:     r.FormValue("reasons"),
		OnBehalfOf:  r.FormValue("onBehalfOf"),
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		Keywords:    strings.Split(r.FormValue("keywords"), ","),
		PinMajority: FormToI(r, "pinMajority"),
	}
	return action
}

func CreateCollectiveForm(r *http.Request) CreateCollective {
	action := CreateCollective{
		Action:      "CreateCollective",
		ID:          FormToI(r, "id"),
		Reasons:     r.FormValue("reasons"),
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		Policy:      FormToPolicy(r),
	}
	return action
}

func CreateEventForm(r *http.Request, handles map[string]crypto.Token) CreateEvent {
	action := CreateEvent{
		Action:       "CreateEvent",
		ID:           FormToI(r, "id"),
		Reasons:      r.FormValue("reasons"),
		OnBehalfOf:   r.FormValue("onBehalfOf"),
		StartAt:      FormToTime(r, "startAt"),
		EstimatedEnd: FormToTime(r, "estimatedEnd"),
		Description:  r.FormValue("description"),
		Venue:        r.FormValue("venue"),
		Open:         FormToBool(r, "open"),
		Public:       FormToBool(r, "publiic"),
		Managers:     FormToTokenArray(r, "managers", handles),
	}
	return action

}

func DraftForm(r *http.Request, handles map[string]crypto.Token) Draft {
	action := Draft{
		Action:        "Draft",
		ID:            FormToI(r, "id"),
		Reasons:       r.FormValue("reasons"),
		OnBehalfOf:    r.FormValue("onBehalfOf"),
		CoAuthors:     FormToTokenArray(r, "coAuthors", handles),
		Title:         r.FormValue("title"),
		Keywords:      FormToStringArray(r, "keywords"),
		ContentType:   r.FormValue("contentType"),
		FilePath:      r.FormValue("filePath"),
		PreviousDraft: FormToHash(r, "PreviousDraft"),
		References:    FormToHashArray(r, "references"),
	}
	return action
}

func EditForm(r *http.Request, handles map[string]crypto.Token) Edit {
	action := Edit{
		Action:      "Edit",
		ID:          FormToI(r, "id"),
		Reasons:     r.FormValue("reasons"),
		OnBehalfOf:  r.FormValue("onBehalfOf"),
		CoAuthors:   FormToTokenArray(r, "coAuthors", handles),
		EditedDraft: FormToHash(r, "editedDraft"),
		ContentType: r.FormValue("contentType"),
		FilePath:    r.FormValue("filePath"),
	}
	return action
}

func ImprintStampForm(r *http.Request) ImprintStamp {
	action := ImprintStamp{
		Action:     "ImprintStamp",
		ID:         FormToI(r, "id"),
		Reasons:    r.FormValue("reasons"),
		OnBehalfOf: r.FormValue("onBehalfOf"),
		Hash:       FormToHash(r, "hash"),
	}
	return action
}

func PinForm(r *http.Request) Pin {
	action := Pin{
		Action:  "Pin",
		ID:      FormToI(r, "id"),
		Reasons: r.FormValue("reasons"),
		Board:   r.FormValue("board"),
		Draft:   FormToHash(r, "draft"),
		Pin:     FormToBool(r, "pin"),
	}
	return action
}

func ReactForm(r *http.Request) React {
	action := React{
		Action:     "React",
		ID:         FormToI(r, "id"),
		Reasons:    r.FormValue("reasons"),
		OnBehalfOf: r.FormValue("onBehalfOf"),
		Hash:       FormToHash(r, "hash"),
		Reaction:   byte(FormToI(r, "reaction")),
	}
	return action
}

func ReleaseDraftForm(r *http.Request) ReleaseDraft {
	action := ReleaseDraft{
		Action:      "ReleaseDraft",
		ID:          FormToI(r, "id"),
		Reasons:     r.FormValue("reasons"),
		ContentHash: FormToHash(r, "contentHash"),
	}
	return action
}

func RemoveMemberForm(r *http.Request, handles map[string]crypto.Token) RemoveMember {
	action := RemoveMember{
		Action:     "RemoveMember",
		ID:         FormToI(r, "id"),
		Reasons:    r.FormValue("reasons"),
		OnBehalfOf: r.FormValue("onBehalfOf"),
		Member:     FormToToken(r, "hash", handles),
	}
	return action
}

func RequestMembershipForm(r *http.Request) RequestMembership {
	action := RequestMembership{
		Action:     "RequestMembership",
		ID:         FormToI(r, "id"),
		Reasons:    r.FormValue("reasons"),
		Collective: r.FormValue("collective"),
		Include:    FormToBool(r, "include"),
	}
	return action
}

func UpdateBoardForm(r *http.Request) UpdateBoard {
	action := UpdateBoard{
		Action:  "UpdateBoard",
		ID:      FormToI(r, "id"),
		Reasons: r.FormValue("reasons"),
	}
	if s := r.FormValue("description"); s != "" {
		action.Description = &s
	}
	if s := r.FormValue("keywords"); s != "" {
		keywords := FormToStringArray(r, "keywords")
		action.Keywords = &keywords
	}
	if s := r.FormValue("pinMajority"); s != "" {
		majorty := FormToI(r, "pinMajorty")
		action.PinMajority = &majorty
	}
	return action
}

func UpdateCollectiveForm(r *http.Request) UpdateCollective {
	action := UpdateCollective{
		Action:     "RequestMembership",
		ID:         FormToI(r, "id"),
		Reasons:    r.FormValue("reasons"),
		OnBehalfOf: r.FormValue("onBehalfOf"),
	}
	if s := r.FormValue("description"); s != "" {
		action.Description = &s
	}
	if s := r.FormValue("policyMajority"); s != "" {
		policy := FormToPolicy(r)
		action.Policy = &policy
	}
	return action
}

func UpdateEventForm(r *http.Request, handles map[string]crypto.Token) UpdateEvent {
	action := UpdateEvent{
		Action:    "UpdateEvent",
		ID:        FormToI(r, "id"),
		Reasons:   r.FormValue("reasons"),
		EventHash: FormToHash(r, "eventHash"),
	}
	if s := r.FormValue("description"); s != "" {
		action.Description = &s
	}
	if s := r.FormValue("venue"); s != "" {
		action.Venue = &s
	}
	if s := r.FormValue("open"); s != "" {
		open := FormToBool(r, "open")
		action.Open = &open
	}
	if s := r.FormValue("public"); s != "" {
		public := FormToBool(r, "public")
		action.Public = &public
	}
	if s := r.FormValue("managers"); s != "" {
		managers := FormToTokenArray(r, "managers", handles)
		action.Managers = &managers
	}
	return action
}

func VoteForm(r *http.Request) Vote {
	action := Vote{
		Action:  "Vote",
		ID:      FormToI(r, "id"),
		Reasons: r.FormValue("reasons"),
		Hash:    FormToHash(r, "hash"),
		Approve: FormToBool(r, "approve"),
	}
	return action
}
