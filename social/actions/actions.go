package actions

import "github.com/lienkolabs/swell/crypto"

const (
	AVote byte = iota
	ACreateCollective
	AUpdateCollective
	ARequestMembership
	ARemoveMember
	ADraft
	AEdit
	AMultipartMedia
	ACreateBoard
	AUpdateBoard
	APin
	ABoardEditor
	AReleaseDraft
	AImprintStamp
	AReact
	ASignIn
	ACreateEvent
	ACancelEvent
	AUpdateEvent
	ACheckinEvent
	AAcceptCheckinEvent
	AUnknown
)

type Action interface {
	Serialize() []byte
}

func ActionKind(data []byte) byte {
	if len(data) < 8+crypto.TokenSize+1 {
		return AUnknown
	}
	actionByte := data[8+crypto.TokenSize]
	if actionByte >= AUnknown {
		return AUnknown
	}
	return actionByte
}
