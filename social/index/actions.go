package index

import (
	"fmt"

	"github.com/lienkolabs/swell/crypto"

	"github.com/lienkolabs/synergy/social/actions"
)

func (i *Index) ActionToObjects(action actions.Action) []crypto.Hash {
	switch v := action.(type) {
	case *actions.ImprintStamp:
		return []crypto.Hash{crypto.Hasher([]byte(v.OnBehalfOf)), v.Hash}
	case *actions.CreateEvent:
		return []crypto.Hash{crypto.Hasher([]byte(v.OnBehalfOf))}
	case *actions.CancelEvent:
		return []crypto.Hash{v.Hash}
	case *actions.UpdateEvent:
		return []crypto.Hash{v.EventHash}
	case *actions.CheckinEvent:
		return []crypto.Hash{v.EventHash}
	case *actions.GreetCheckinEvent:
		return []crypto.Hash{v.EventHash}
	case *actions.CreateBoard:
		return []crypto.Hash{crypto.Hasher([]byte(v.OnBehalfOf))}
	case *actions.UpdateBoard:
		hash := crypto.Hasher([]byte(v.Board))
		return []crypto.Hash{hash}
	case *actions.Pin:
		return []crypto.Hash{crypto.Hasher([]byte(v.Board)), v.Draft}
	case *actions.BoardEditor:
		return []crypto.Hash{crypto.Hasher([]byte(v.Board))}
	case *actions.Draft:
		if v.OnBehalfOf != "" {
			return []crypto.Hash{crypto.Hasher([]byte(v.OnBehalfOf))}
		}
		return nil
	case *actions.ReleaseDraft:
		return []crypto.Hash{v.ContentHash}
	case *actions.Edit:
		if v.OnBehalfOf != "" {
			return []crypto.Hash{crypto.Hasher([]byte(v.OnBehalfOf)), v.EditedDraft}
		}
		return []crypto.Hash{v.EditedDraft}
	case *actions.React:
		if v.OnBehalfOf != "" {
			return []crypto.Hash{crypto.Hasher([]byte(v.OnBehalfOf)), v.Hash}
		}
		return []crypto.Hash{v.Hash}
	case *actions.CreateCollective:
		return []crypto.Hash{crypto.ZeroHash}
	case *actions.UpdateCollective:
		return []crypto.Hash{crypto.Hasher([]byte(v.OnBehalfOf))}
	case *actions.RequestMembership:
		return []crypto.Hash{crypto.Hasher([]byte(v.Collective))}
	case *actions.RemoveMember:
		return []crypto.Hash{crypto.Hasher([]byte(v.OnBehalfOf))}
	case *actions.Signin:
		return []crypto.Hash{crypto.ZeroHash}
	}
	return nil
}

func (i *Index) ActionToString(action actions.Action, status bool) (string, uint64) {
	switch v := action.(type) {
	case *actions.ImprintStamp:
		if draft, ok := i.state.Drafts[v.Hash]; ok {
			if status {
				return fmt.Sprintf("%v stamped released draft %v", v.OnBehalfOf, draft.Title), v.Epoch
			} else {
				handle := i.state.Members[crypto.HashToken(v.Author)]
				return fmt.Sprintf("%v proposed stamp on behalf of %v for the released draft %v", handle, v.OnBehalfOf, draft.Title), v.Epoch
			}
		}
		return "", 0
	case *actions.CreateEvent:
		return "", 0
	case *actions.CancelEvent:
		return "", 0
	case *actions.UpdateEvent:
		return "", 0
	case *actions.CheckinEvent:
		return "", 0
	case *actions.GreetCheckinEvent:
		return "", 0
	case *actions.CreateBoard:
		return "", 0
	case *actions.UpdateBoard:
		return "", 0
	case *actions.Pin:
		return "", 0
	case *actions.BoardEditor:
		return "", 0
	case *actions.Draft:
		return "", 0
	case *actions.ReleaseDraft:
		return "", 0
	case *actions.Edit:
		return "", 0
	case *actions.React:
		return "", 0
	case *actions.CreateCollective:
		return "", 0
	case *actions.UpdateCollective:
		return "", 0
	case *actions.RequestMembership:
		return "", 0
	case *actions.RemoveMember:
		return "", 0
	case *actions.Signin:
		return "", 0
	}
	return "", 0
}
