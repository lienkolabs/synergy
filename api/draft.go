package main

import (
	"os"

	"github.com/lienkolabs/swell/crypto"
	"github.com/lienkolabs/synergy"
)

const FilePart = 100000

type Signer interface {
	Epoch() uint64
	Token() crypto.Token
}

func InstructDraft(draft *CreateDraft, signer Signer) (*synergy.AlternativeDraftInstruction, []*synergy.MultipartMedia) {
	inst := synergy.AlternativeDraftInstruction{
		Epoch:      signer.Epoch(),
		Author:     signer.Token(),
		OnBehalfOf: draft.OnBehalfOf,
		Policy: synergy.Policy{
			Majority:      draft.Policy.Majority,
			SuperMajority: draft.Policy.SuperMajority,
		},
		Reasons:       draft.Reasons,
		Title:         draft.Title,
		Keywords:      draft.Keywords,
		Description:   draft.Description,
		ContentType:   draft.ContentType,
		PreviousDraft: DecodeHash(draft.PreviousDraft),
	}
	if draft.CoAuthors != nil {
		inst.CoAuthors = make([]crypto.Token, len(draft.CoAuthors))
		for n, coAuthor := range draft.CoAuthors {
			inst.CoAuthors[n] = DecodeToken(coAuthor)
		}
	}
	if draft.References != nil {
		inst.References = make([]crypto.Hash, len(draft.References))
		for n, reference := range draft.References {
			inst.References[n] = DecodeHash(reference)
		}
	}
	bytes, err := os.ReadFile(draft.FilePath)
	if err != nil {
		return nil, nil
	}
	inst.ContentHash = crypto.Hasher(bytes)
	if len(bytes) > 254*100000 {
		return nil, nil
	}
	inst.NumberOfParts = byte(len(bytes)/100000) + 1
	if inst.NumberOfParts == 1 {
		inst.Content = bytes
		return &inst, nil
	}
	inst.Content = bytes[0:FilePart]
	multiPart := make([]*synergy.MultipartMedia, inst.NumberOfParts-2)
	for n := 1; n < int(inst.NumberOfParts); n++ {
		multiPart[n-1] = &synergy.MultipartMedia{
			Hash: inst.ContentHash,
			Part: byte(n + 1),
			Of:   inst.NumberOfParts,
			Data: bytes[n*FilePart : (n+1)*FilePart],
		}
	}
	return &inst, multiPart
}
