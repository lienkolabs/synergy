package actions

import (
	"reflect"
	"testing"

	"github.com/lienkolabs/swell/crypto"
)

var (
	content = make([]byte, 10)

	draft = &Draft{
		Epoch:         18,
		Author:        crypto.Token{},
		Reasons:       "draft test",
		OnBehalfOf:    "first_collective",
		CoAuthors:     []crypto.Token{},
		Policy:        policy,
		Title:         "first_draft",
		Keywords:      exampleArray[:],
		Description:   "draft test",
		ContentType:   "txt",
		ContentHash:   crypto.Hash{},
		NumberOfParts: 1,
		Content:       content,
		PreviousDraft: crypto.Hash{},
		References:    []crypto.Hash{},
	}
)

func TestDraft(t *testing.T) {
	d := ParseDraft(draft.Serialize())
	if d == nil {
		t.Error("Could not parse actions Draft")
		return
	}
	if !reflect.DeepEqual(d, draft) {
		t.Error("Parse and Serialize not working for AcceptJoinAudience")
	}
}

// func Test
