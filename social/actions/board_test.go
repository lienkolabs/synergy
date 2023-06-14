package actions

import (
	"reflect"
	"testing"

	"github.com/lienkolabs/swell/crypto"
)

var (
	exampleArray = [2]string{"test", "board"}

	board = &CreateBoard{
		Epoch:       10,
		Author:      crypto.Token{},
		Reasons:     "create board test",
		OnBehalfOf:  "coletivo_teste",
		Name:        "first_board",
		Description: "create board test",
		Keywords:    exampleArray[:],
		PinMajority: 1,
	}

	uBoard = &UpdateBoard{
		Epoch:       11,
		Author:      crypto.Token{},
		Reasons:     "update board test",
		Board:       "first_board",
		Description: "creat board test description updated",
		Keywords:    exampleArray[:],
		PinMajority: 2,
	}

	pin = &Pin{
		Epoch:   12,
		Author:  crypto.Token{},
		Reasons: "pin on board test",
		Board:   "first_board",
		Draft:   crypto.Hash{},
		Pin:     true,
	}

	editor = &BoardEditor{
		Epoch:   13,
		Author:  crypto.Token{},
		Reasons: "include board editor test",
		Board:   "first_board",
		Editor:  crypto.Token{},
		Insert:  true,
	}
)

func TestCreateBoard(t *testing.T) {
	b := ParseCreateBoard(board.Serialize())
	if b == nil {
		t.Error("Coult not parse actions CreateBoard")
		return
	}
	if !reflect.DeepEqual(b, board) {
		t.Error("Parse and Serialize not working for actions CreateBoard ")
	}
}

func TestUpdateBoard(t *testing.T) {
	b := ParseUpdateBoard(uBoard.Serialize())
	if b == nil {
		t.Error("Could not parse actions UpdateBoard")
		return
	}
	if !reflect.DeepEqual(b, uBoard) {
		t.Error("Parse and Serialize not working for actions UpdateBoard")
	}
}

func TestPin(t *testing.T) {
	p := ParsePin(pin.Serialize())
	if p == nil {
		t.Error("Could not parse actions Pin")
		return
	}
	if !reflect.DeepEqual(p, pin) {
		t.Error("Parse and Serialize not working for actions Pin")
	}
}

func TestBoardEditor(t *testing.T) {
	e := ParseBoardEditor(editor.Serialize())
	if e == nil {
		t.Error("Could not parse actions BoardEditor")
		return
	}
	if !reflect.DeepEqual(e, editor) {
		t.Error("Parse and Serialize not working for actions BoardEditor")
	}
}
