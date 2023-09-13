package index

import (
	"fmt"

	"github.com/lienkolabs/breeze/crypto"

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
				return fmt.Sprintf("%v proposed %v stamp for %v draft release", handle, v.OnBehalfOf, draft.Title), v.Epoch
			}
		}
		return "", 0
	case *actions.CreateEvent:
		// hash do evento eh o hash da acao do evento
		eventhash := v.Hashed()
		if event, ok := i.state.Events[eventhash]; ok {
			if status {
				return fmt.Sprintf("%v event created on behalf of %v", v.OnBehalfOf, event.Collective.Name), v.Epoch
			} else {
				handle := i.state.Members[crypto.HashToken(v.Author)]
				return fmt.Sprintf("%v proposed a %v event on behalf of %v", handle, v.StartAt.Format("2006-01-02"), v.OnBehalfOf), v.Epoch
			}
		}
		return "", 0
	case *actions.CancelEvent:
		// hash eh o hash do evento original
		if event, ok := i.state.Events[v.Hash]; ok {
			if status {
				return fmt.Sprintf("%v event cancelled on behalf of %v", event.Collective.Name, event.Collective.Name), v.Epoch
			} else {
				handle := i.state.Members[crypto.HashToken(v.Author)]
				return fmt.Sprintf("%v proposed %v event cancellation on behalf of %v", handle, event.StartAt.Format("2006-01-02"), event.Collective.Name), v.Epoch
			}
		}
		return "", 0
	case *actions.UpdateEvent:
		// hash eh o hash do evento original
		if event, ok := i.state.Events[v.EventHash]; ok {
			if status {
				return fmt.Sprintf("%v event update on behalf of %v", event.StartAt.Format("2006-01-02"), event.Collective.Name), v.Epoch
			} else {
				handle := i.state.Members[crypto.HashToken(v.Author)]
				return fmt.Sprintf("%v proposed update for %v event on behalf of %v", handle, event.StartAt.Format("2006-01-02"), event.Collective.Name), v.Epoch
			}
		}
		return "", 0
	case *actions.CheckinEvent:
		// hash eh o hash do evento
		if event, ok := i.state.Events[v.EventHash]; ok {
			if status {
				handle := i.state.Members[crypto.HashToken(v.Author)]
				return fmt.Sprintf("%v checkedin on %v event by %v ", handle, event.StartAt.Format("2006-01-02"), event.Collective.Name), v.Epoch
			} else {
				return "ERRO: CheckInEvent nao esta sendo processado de forma automatica", 0
			}
		}
		return "", 0
	case *actions.GreetCheckinEvent:
		if event, ok := i.state.Events[v.EventHash]; ok {
			if status {
				handle := i.state.Members[crypto.HashToken(v.Author)]
				return fmt.Sprintf("%v greeted checkin on %v event by %v ", handle, event.StartAt.Format("2006-01-02"), event.Collective.Name), v.Epoch
			}
		}
		return "", 0
	case *actions.CreateBoard:
		// hash do board eh o hash do nome do board que esta sendo criado
		boardhash := v.Hashed()
		if board, ok := i.state.Boards[boardhash]; ok {
			if status {
				return fmt.Sprintf("%v board created on behalf of %v", board.Name, v.OnBehalfOf), v.Epoch
			} else {
				handle := i.state.Members[crypto.HashToken(v.Author)]
				return fmt.Sprintf("%v proposed the creation of %v board on behalf of %v", handle, v.Name, v.OnBehalfOf), v.Epoch
			}
		}
		return "", 0
	case *actions.UpdateBoard:
		hash := crypto.Hasher([]byte(v.Board))
		if board, ok := i.state.Boards[hash]; ok {
			if status {
				return fmt.Sprintf("%v board updated on behalf of %v", board.Name, board.Collective.Name), v.Epoch
			} else {
				handle := i.state.Members[crypto.HashToken(v.Author)]
				return fmt.Sprintf("%v proposed update of %v board on behalf of %v", handle, board.Name, board.Collective.Name), v.Epoch
			}
		}
		return "", 0
	case *actions.Pin:
		hash := crypto.Hasher([]byte(v.Board))
		if board, ok := i.state.Boards[hash]; ok {
			draftitle := i.state.Drafts[v.Draft].Title
			pinaction := []string{"unpinned from", "unpin"}
			if v.Pin {
				pinaction = []string{"pinned on", "pin"}
			}
			if status {
				return fmt.Sprintf("%v draft %v %v board on behalf of %v", draftitle, pinaction[0], board.Name, board.Collective.Name), v.Epoch
			} else {
				handle := i.state.Members[crypto.HashToken(v.Author)]
				return fmt.Sprintf("%v proposed %v of %v draft on %v board on behalf of %v", handle, pinaction[1], draftitle, board.Name, board.Collective.Name), v.Epoch
			}
		}
		return "", 0
	case *actions.BoardEditor:
		hash := crypto.Hasher([]byte(v.Board))
		if board, ok := i.state.Boards[hash]; ok {
			editor := i.state.Members[crypto.HashToken(v.Editor)]
			editorship := []string{"removed from", "removal of"}
			if v.Insert {
				editorship = []string{"included for", "inclusion of"}
			}
			if status {
				return fmt.Sprintf("%v editor %v editorship on %v board on behalf of %v", editor, editorship[0], board.Name, board.Collective.Name), v.Epoch
			} else {
				handle := i.state.Members[crypto.HashToken(v.Author)]
				return fmt.Sprintf("%v proposed %v of %v editor on %v board on behalf of %v", handle, editorship[1], editor, board.Name, board.Collective.Name), v.Epoch
			}
		}
		return "", 0
	case *actions.Draft:
		if draft, ok := i.state.Drafts[v.ContentHash]; ok {
			if draft.Authors.CollectiveName() != "" {
				if status {
					return fmt.Sprintf("%v draft was created on behalf of %v", draft.Title, draft.Authors.CollectiveName()), v.Epoch
				} else {
					handle := i.state.Members[crypto.HashToken(v.Author)]
					return fmt.Sprintf("%v proposed %v draft on behalf of %v", handle, draft.Title, draft.Authors.CollectiveName()), v.Epoch
				}
			}
			return "", 0
		}
		return "", 0
	case *actions.ReleaseDraft:
		if draft, ok := i.state.Drafts[v.ContentHash]; ok {
			if status {
				return fmt.Sprintf("%v draft was released on behalf of %v", draft.Title, draft.Authors.CollectiveName()), v.Epoch
			} else {
				handle := i.state.Members[crypto.HashToken(v.Author)]
				return fmt.Sprintf("%v proposed %v draft release on behalf of %v", handle, draft.Title, draft.Authors.CollectiveName()), v.Epoch
			}
		}
		return "", 0
	case *actions.Edit:
		if edit, ok := i.state.Edits[v.ContentHash]; ok {
			draft := i.state.Drafts[v.EditedDraft]
			if edit.Authors.CollectiveName() != "" {
				if status {
					return fmt.Sprintf("edit to %v draft was suggested on behalf of %v", draft.Title, edit.Authors.CollectiveName()), v.Epoch
				} else {
					handle := i.state.Members[crypto.HashToken(v.Author)]
					return fmt.Sprintf("%v proposed edit suggestion to %v draft on behalf of %v", handle, draft.Title, edit.Authors.CollectiveName()), v.Epoch
				}
			}
			if status {
				return fmt.Sprintf("edit to %v draft was suggested", draft.Title), v.Epoch
			} else {
				handle := i.state.Members[crypto.HashToken(v.Author)]
				return fmt.Sprintf("%v proposed edit suggestion to %v draft", handle, draft.Title), v.Epoch
			}
		}
		return "", 0
	case *actions.React:
		// reacthash := v.Hashed()
		// if ok := i.state.Reactions[reacthash]; ok {
		// if status {
		// 	return fmt.Sprintf("%v draft was released on behalf of %v", draft.Title, draft.Authors.CollectiveName()), v.Epoch
		// } else {
		// 	handle := i.state.Members[crypto.HashToken(v.Author)]
		// 	return fmt.Sprintf("%v proposed %v draft release on behalf of %v", handle, draft.Title, draft.Authors.CollectiveName()), v.Epoch
		// }
		// }
		return "", 0
	case *actions.CreateCollective:
		collectivehash := crypto.Hasher([]byte(v.Name))
		if collective, ok := i.state.Collectives[collectivehash]; ok {
			if status {
				return fmt.Sprintf("collective %v was created", collective.Name), v.Epoch
			} else {
				handle := i.state.Members[crypto.HashToken(v.Author)]
				return fmt.Sprintf("%v proposed %v collective creation", handle, collective.Name), v.Epoch
			}
		}
		return "", 0
	case *actions.UpdateCollective:
		collectivehash := crypto.Hasher([]byte(v.OnBehalfOf))
		if collective, ok := i.state.Collectives[collectivehash]; ok {
			if status {
				return fmt.Sprintf("%v collective update", collective.Name), v.Epoch
			} else {
				handle := i.state.Members[crypto.HashToken(v.Author)]
				return fmt.Sprintf("%v proposed update for %v collective", handle, collective.Name), v.Epoch
			}
		}
		return "", 0
	case *actions.RequestMembership:
		collectivehash := crypto.Hasher([]byte(v.Collective))
		if collective, ok := i.state.Collectives[collectivehash]; ok {
			handle := i.state.Members[crypto.HashToken(v.Author)]
			if v.Include {
				if status {
					return fmt.Sprintf("%v became a member of %v collective", handle, collective.Name), v.Epoch
				} else {
					return fmt.Sprintf("%v requested membership to %v collective", handle, collective.Name), v.Epoch
				}
			}
			if status {
				return fmt.Sprintf("%v left %v collective", handle, collective.Name), v.Epoch
			}
		}
		return "", 0
	case *actions.RemoveMember:
		collectivehash := crypto.Hasher([]byte(v.OnBehalfOf))
		if collective, ok := i.state.Collectives[collectivehash]; ok {
			member := i.state.Members[crypto.HashToken(v.Member)]
			if status {
				return fmt.Sprintf("%v was removed from %v collective", member, collective.Name), v.Epoch
			} else {
				handle := i.state.Members[crypto.HashToken(v.Author)]
				return fmt.Sprintf("%v requested removal of %v from %v collective", handle, member, collective.Name), v.Epoch
			}
		}
		return "", 0
	case *actions.Signin:
		authorhash := crypto.HashToken(v.Author)
		if _, ok := i.state.Members[authorhash]; ok {
			if status {
				return fmt.Sprintf("%v member joined Synergy", v.Handle), v.Epoch
			}
		}
		return "", 0
	}
	return "", 0
}
