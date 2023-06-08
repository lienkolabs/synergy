package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"os"

	"github.com/lienkolabs/swell/crypto"
	"github.com/lienkolabs/synergy"
)

var ZeroHash crypto.Hash

type Policy struct {
	Majority      int `json:"majority"`
	SuperMajority int `json:"superMajority"`
}

type Instruction struct {
	Instruction string `json:"instruction"`
}

func JSONType(data []byte) string {
	var inst Instruction
	json.Unmarshal(data, &inst)
	return inst.Instruction
}

func ParseCreateDraft(data []byte) *CreateDraft {
	var create CreateDraft
	if json.Unmarshal(data, &create) != nil {
		return nil
	}
	return &create
}

type CreateDraft struct {
	Instruction   string   `json:"instruction"`
	OnBehalfOf    string   `json:"onBeahlfOf,omitempty"`
	CoAuthors     []string `json:"coAuthors,omitempty"`
	Policy        Policy   `json:"policy"`
	Reasons       string   `json:"reasons"`
	Title         string   `json:"title"`
	Keywords      []string `json:"keywords"`
	Description   string   `json:"description"`
	ContentType   string   `json:"contentType"`
	FilePath      string   `json:"filePath"`
	PreviousDraft string   `json:"previousDraft,,omitempty"`
	References    []string `json:"references,omitempty"`
}

type CreateEdit struct {
	OnBehalfOf string   `json:"onBeahlfOf,omitempty"`
	CoAuthors  []string `json:"coAuthors,omitempty"`
	Reasons    string   `json:"reasons"`
	Draft      string   `json:"draft"`
	EditType   string   `json:"editType"`
	FilePath   string   `json:"filePath"`
}

type VoteInstruction struct {
	Instruction string `json:"instruction"`
	Reasons     string `json:"reasons,omitempty"`
	Hash        string `json:"hash"`
	Approve     bool
}

type CreateCollective struct {
	Instruction string  `json:"instruction"`
	Name        string  `json:"name"`
	Weight      int     `json:"weight"`
	Policy      *Policy `json:"policy,omitempty"`
}

type CreateBoard struct {
	Instruction string   `json:"instruction"`
	Name        string   `json:"name"`
	Weight      int      `json:"weight"`
	Policy      *Policy  `json:"policy,omitempty"`
	Keywords    []string `json:"keywords,omitempty"`
}

type CreateJournal struct {
	Instruction string   `json:"instruction"`
	Name        string   `json:"name"`
	Weight      int      `json:"weight"`
	Policy      *Policy  `json:"policy,omitempty"`
	Keywords    []string `json:"keywords,omitempty"`
}

type Release struct {
	Instruction string `json:"instruction"`
	DraftHash   string `json:"draft"`
	Journal     string `json:"journal"`
}

type Publish struct {
	Instruction string `json:"instruction"`
	Journal     string `json:"journal"`
}

type UpdateMembers struct {
	Instruction string `json:"instruction"`
	Collective  string `json:"collective"`
	Token       string `json:"token"`
	Weight      int    `json:"weight"`
	Include     bool   `json:"include"`
}

type PinBoard struct {
	Instruction string `json:"instruction"`
	Board       string `json:"board"`
	Draft       string `json:"draft"`
	Pin         bool   `json:"pin"`
}

type UpdateBoardEditors struct {
	Instruction string `json:"instruction"`
	Board       string `json:"board"`
	Token       string `json:"token"`
	Weight      int    `json:"weight"`
	Include     bool   `json:"include"`
}

func DraftInstructionToJSON(d *synergy.AlternativeDraftInstruction) string {
	c := CreateDraft{
		Instruction: "Create Draft",
		OnBehalfOf:  d.OnBehalfOf,
		Policy: Policy{
			Majority:      d.Policy.Majority,
			SuperMajority: d.Policy.SuperMajority,
		},
		Reasons:     d.Reasons,
		Title:       d.Title,
		Keywords:    d.Keywords,
		Description: d.Description,
		ContentType: d.ContentType,
	}
	if d.PreviousDraft != ZeroHash {
		c.PreviousDraft = hex.EncodeToString(d.PreviousDraft[:])
	}
	if d.CoAuthors != nil && len(d.CoAuthors) > 0 {
		c.CoAuthors = make([]string, len(d.CoAuthors))
		for n, author := range d.CoAuthors {
			c.CoAuthors[n] = hex.EncodeToString(author[:])
		}
	}
	if d.References != nil && len(d.References) > 0 {
		c.References = make([]string, len(d.References))
		for n, reference := range d.References {
			c.References[n] = hex.EncodeToString(reference[:])
		}
	}
	data, err := json.Marshal(c)
	if err == nil {
		var out bytes.Buffer
		json.Indent(&out, data, "", "\t")
		out.WriteTo(os.Stdout)
		return ""
	}
	return ""
}

func DecodeHash(value string) crypto.Hash {
	var hash crypto.Hash
	bytes, err := hex.DecodeString(value)
	if err == nil && len(bytes) == crypto.Size {
		copy(hash[:], bytes)
	}
	return hash
}

func DecodeToken(value string) crypto.Token {
	var token crypto.Token
	bytes, err := hex.DecodeString(value)
	if err == nil && len(bytes) == crypto.Size {
		copy(token[:], bytes)
	}
	return token
}
