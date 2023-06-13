package api

import (
	"os"

	"github.com/lienkolabs/swell/crypto"
	"github.com/lienkolabs/synergy/social/actions"
)

const FilePart = 100000

type Signer interface {
	Epoch() uint64
	Token() crypto.Token
}

func InstructDraft(draft *CreateDraft, a *Attorney) (*actions.Draft, []*actions.MultipartMedia) {
	action := actions.Draft{
		Epoch:         a.epoch,
		Author:        a.author,
		OnBehalfOf:    draft.OnBehalfOf,
		Reasons:       draft.Reasons,
		Title:         draft.Title,
		Keywords:      draft.Keywords,
		Description:   draft.Description,
		ContentType:   draft.ContentType,
		PreviousDraft: draft.PreviousDraft,
	}
	if draft.Policy != nil {
		action.Policy = &actions.Policy{
			Majority:      draft.Policy.Majority,
			SuperMajority: draft.Policy.SuperMajority,
		}
	}
	if draft.CoAuthors != nil {
		action.CoAuthors = make([]crypto.Token, len(draft.CoAuthors))
		for n, coAuthor := range draft.CoAuthors {
			action.CoAuthors[n] = coAuthor
		}
	}
	if draft.References != nil {
		action.References = make([]crypto.Hash, len(draft.References))
		for n, reference := range draft.References {
			action.References[n] = DecodeHash(reference)
		}
	}
	bytes, err := os.ReadFile(draft.FilePath)
	if err != nil {
		return nil, nil
	}
	action.ContentHash = crypto.Hasher(bytes)
	if len(bytes) > 254*100000 {
		return nil, nil
	}
	action.NumberOfParts = byte(len(bytes)/100000) + 1
	if action.NumberOfParts == 1 {
		action.Content = bytes
		return &action, nil
	}
	action.Content = bytes[0:FilePart]
	multiPart := make([]*actions.MultipartMedia, action.NumberOfParts-2)
	for n := 1; n < int(action.NumberOfParts); n++ {
		multiPart[n-1] = &actions.MultipartMedia{
			Hash: action.ContentHash,
			Part: byte(n + 1),
			Of:   action.NumberOfParts,
			Data: bytes[n*FilePart : (n+1)*FilePart],
		}
	}
	return &action, multiPart
}
