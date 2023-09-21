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

func (i *Index) ActionToString(action actions.Action, status bool) (string, string, crypto.Token, uint64, string) {
	switch v := action.(type) {
	case *actions.ImprintStamp:
		////fmt.Println("stamp")
		if draft, ok := i.state.Drafts[v.Hash]; ok {
			if status {
				return fmt.Sprintf("%v stamped %v", v.OnBehalfOf, draft.Title), crypto.EncodeHash(draft.DraftHash), v.Author, v.Epoch, "stamp"
			} else {
				handle := i.state.Members[crypto.HashToken(v.Author)]
				return fmt.Sprintf("%v proposed %v stamp for %v", handle, v.OnBehalfOf, draft.Title), crypto.EncodeHash(draft.DraftHash), v.Author, v.Epoch, "stamp"
			}
		}
		return "", "", v.Author, 0, ""
	case *actions.CreateEvent:
		////fmt.Println("cevent")
		// hash do evento eh o hash da acao do evento
		eventhash := v.Hashed()
		if event, ok := i.state.Events[eventhash]; ok {
			if status {
				return fmt.Sprintf("%v event created on behalf of %v", v.OnBehalfOf, event.Collective.Name), crypto.EncodeHash(eventhash), v.Author, v.Epoch, "create event"
			} else {
				handle := i.state.Members[crypto.HashToken(v.Author)]
				return fmt.Sprintf("%v proposed a %v event on behalf of %v", handle, v.StartAt.Format("2006-01-02"), v.OnBehalfOf), crypto.EncodeHash(eventhash), v.Author, v.Epoch, "create event"
			}
		}
		return "", "", v.Author, 0, ""
	case *actions.CancelEvent:
		//fmt.Println("cancel event")
		// hash eh o hash do evento original
		if event, ok := i.state.Events[v.Hash]; ok {
			if status {
				return fmt.Sprintf("%v event cancelled on behalf of %v", event.Collective.Name, event.Collective.Name), crypto.EncodeHash(v.Hash), v.Author, v.Epoch, "cancel event"
			} else {
				handle := i.state.Members[crypto.HashToken(v.Author)]
				return fmt.Sprintf("%v proposed %v event cancellation on behalf of %v", handle, event.StartAt.Format("2006-01-02"), event.Collective.Name), crypto.EncodeHash(v.Hash), v.Author, v.Epoch, "cancel event"
			}
		}
		return "", "", v.Author, 0, ""
	case *actions.UpdateEvent:
		//fmt.Println("update event")
		// hash eh o hash do evento original
		if event, ok := i.state.Events[v.EventHash]; ok {
			if status {
				return fmt.Sprintf("%v event update on behalf of %v", event.StartAt.Format("2006-01-02"), event.Collective.Name), crypto.EncodeHash(v.EventHash), v.Author, v.Epoch, "update event"
			} else {
				handle := i.state.Members[crypto.HashToken(v.Author)]
				return fmt.Sprintf("%v proposed update for %v event on behalf of %v", handle, event.StartAt.Format("2006-01-02"), event.Collective.Name), crypto.EncodeHash(v.EventHash), v.Author, v.Epoch, "update event"
			}
		}
		return "", "", v.Author, 0, ""
	case *actions.CheckinEvent:
		//fmt.Println("checik")
		// hash eh o hash do evento
		if event, ok := i.state.Events[v.EventHash]; ok {
			if status {
				handle := i.state.Members[crypto.HashToken(v.Author)]
				return fmt.Sprintf("%v checkedin on %v event by %v ", handle, event.StartAt.Format("2006-01-02"), event.Collective.Name), crypto.EncodeHash(v.EventHash), v.Author, v.Epoch, "event checkin"
			}
		}
		return "", "", v.Author, 0, ""
	case *actions.GreetCheckinEvent:
		//fmt.Println("greet")
		if event, ok := i.state.Events[v.EventHash]; ok {
			if status {
				handle := i.state.Members[crypto.HashToken(v.Author)]
				return fmt.Sprintf("%v greeted checkin on %v event by %v ", handle, event.StartAt.Format("2006-01-02"), event.Collective.Name), crypto.EncodeHash(v.EventHash), v.Author, v.Epoch, "event greet"
			}
		}
		return "", "", v.Author, 0, ""
	case *actions.CreateBoard:
		//fmt.Println("cboard")
		// hash do board eh o hash do nome do board que esta sendo criado
		boardhash := v.Hashed()
		if board, ok := i.state.Boards[boardhash]; ok {
			if status {
				return fmt.Sprintf("%v created on behalf of %v", board.Name, v.OnBehalfOf), crypto.EncodeHash(boardhash), v.Author, v.Epoch, "create board"
			} else {
				handle := i.state.Members[crypto.HashToken(v.Author)]
				return fmt.Sprintf("%v proposed the creation of %v on behalf of %v", handle, v.Name, v.OnBehalfOf), crypto.EncodeHash(boardhash), v.Author, v.Epoch, "create board"
			}
		}
		return "", "", v.Author, 0, ""
	case *actions.UpdateBoard:
		//fmt.Println("upboard")
		hash := crypto.Hasher([]byte(v.Board))
		if board, ok := i.state.Boards[hash]; ok {
			if status {
				return fmt.Sprintf("%v updated on behalf of %v", board.Name, board.Collective.Name), crypto.EncodeHash(hash), v.Author, v.Epoch, "update board"
			} else {
				handle := i.state.Members[crypto.HashToken(v.Author)]
				return fmt.Sprintf("%v proposed update of %v on behalf of %v", handle, board.Name, board.Collective.Name), crypto.EncodeHash(hash), v.Author, v.Epoch, "update board"
			}
		}
		return "", "", v.Author, 0, ""
	case *actions.Pin:
		//fmt.Println("pin")
		hash := crypto.Hasher([]byte(v.Board))
		if board, ok := i.state.Boards[hash]; ok {
			if draft, ok := i.state.Drafts[v.Draft]; ok {
				pinaction := []string{"unpinned from", "unpin from"}
				if v.Pin {
					pinaction = []string{"pinned on", "pin on"}
				}
				if status {
					return fmt.Sprintf("%v %v %v on behalf of %v", draft.Title, pinaction[0], board.Name, board.Collective.Name), crypto.EncodeHash(draft.DraftHash), v.Author, v.Epoch, pinaction[1]
				} else {
					handle := i.state.Members[crypto.HashToken(v.Author)]
					return fmt.Sprintf("%v proposed %v of %v %v on behalf of %v", handle, pinaction[1], draft.Title, board.Name, board.Collective.Name), crypto.EncodeHash(draft.DraftHash), v.Author, v.Epoch, pinaction[1]
				}
			}
		}
		return "", "", v.Author, 0, ""
	case *actions.BoardEditor:
		//fmt.Println("beditor")
		hash := crypto.Hasher([]byte(v.Board))
		if board, ok := i.state.Boards[hash]; ok {
			editor := i.state.Members[crypto.HashToken(v.Editor)]
			editorship := []string{"removed from", "removal of", "from", "editor removal"}
			if v.Insert {
				editorship = []string{"included for", "inclusion of", "for", "editor inclusion"}
			}
			if status {
				return fmt.Sprintf("%v %v editorship of %v on behalf of %v", editor, editorship[0], board.Name, board.Collective.Name), crypto.EncodeHash(hash), v.Author, v.Epoch, editorship[3]
			} else {
				handle := i.state.Members[crypto.HashToken(v.Author)]
				return fmt.Sprintf("%v proposed %v %v %v %v editorship on behalf of %v", handle, editorship[1], editor, editorship[2], board.Name, board.Collective.Name), crypto.EncodeHash(hash), v.Author, v.Epoch, editorship[3]
			}
		}
		return "", "", v.Author, 0, ""
	case *actions.Draft:
		//fmt.Println("draft")
		if draft, ok := i.state.Drafts[v.ContentHash]; ok {
			if draft.Authors.CollectiveName() != "" {
				if status {
					return fmt.Sprintf("%v was created on behalf of %v", draft.Title, draft.Authors.CollectiveName()), crypto.EncodeHash(draft.DraftHash), v.Author, v.Epoch, "new draft"
				} else {
					handle := i.state.Members[crypto.HashToken(v.Author)]
					return fmt.Sprintf("%v proposed %v on behalf of %v", handle, draft.Title, draft.Authors.CollectiveName()), crypto.EncodeHash(draft.DraftHash), v.Author, v.Epoch, "new draft"
				}
			}
		}
		return "", "", v.Author, 0, ""
	case *actions.ReleaseDraft:
		//fmt.Println("release")
		if draft, ok := i.state.Drafts[v.ContentHash]; ok {
			if status {
				return fmt.Sprintf("%v was released on behalf of %v", draft.Title, draft.Authors.CollectiveName()), crypto.EncodeHash(draft.DraftHash), v.Author, v.Epoch, "release"
			} else {
				handle := i.state.Members[crypto.HashToken(v.Author)]
				return fmt.Sprintf("%v proposed %v release on behalf of %v", handle, draft.Title, draft.Authors.CollectiveName()), crypto.EncodeHash(draft.DraftHash), v.Author, v.Epoch, "release"
			}
		}
		return "", "", v.Author, 0, ""
	case *actions.Edit:
		//fmt.Println("edit")
		if edit, ok := i.state.Edits[v.ContentHash]; ok {
			draft := i.state.Drafts[v.EditedDraft]
			if edit.Authors.CollectiveName() != "" {
				if status {
					return fmt.Sprintf(" %v edit was suggested on behalf of %v", draft.Title, edit.Authors.CollectiveName()), crypto.EncodeHash(v.ContentHash), v.Author, v.Epoch, "edit"
				} else {
					handle := i.state.Members[crypto.HashToken(v.Author)]
					return fmt.Sprintf("%v proposed %v's edit on behalf of %v", handle, draft.Title, edit.Authors.CollectiveName()), crypto.EncodeHash(v.ContentHash), v.Author, v.Epoch, "edit"
				}
			}
			if status {
				return fmt.Sprintf("%v edit was suggested", draft.Title), crypto.EncodeHash(v.ContentHash), v.Author, v.Epoch, "edit"
			} else {
				handle := i.state.Members[crypto.HashToken(v.Author)]
				return fmt.Sprintf("%v proposed %v's edit", handle, draft.Title), crypto.EncodeHash(v.ContentHash), v.Author, v.Epoch, "edit"
			}
		}
		return "", "", v.Author, 0, ""
	case *actions.React:
		//fmt.Println("reatc")
		// reacthash := v.Hashed()
		// if ok := i.state.Reactions[reacthash]; ok {
		// if status {
		// 	return fmt.Sprintf("%v draft was released on behalf of %v", draft.Title, draft.Authors.CollectiveName()), v.Epoch
		// } else {
		// 	handle := i.state.Members[crypto.HashToken(v.Author)]
		// 	return fmt.Sprintf("%v proposed %v draft release on behalf of %v", handle, draft.Title, draft.Authors.CollectiveName()), v.Epoch
		// }
		// }
		return "", "", v.Author, 0, ""
	case *actions.CreateCollective:
		//fmt.Println("ccol")
		collectivehash := crypto.Hasher([]byte(v.Name))
		if collective, ok := i.state.Collectives[collectivehash]; ok {
			if status {
				return fmt.Sprintf("%v was created", collective.Name), crypto.EncodeHash(collectivehash), v.Author, v.Epoch, "create collective"
			} else {
				handle := i.state.Members[crypto.HashToken(v.Author)]
				return fmt.Sprintf("%v proposed %v creation", handle, collective.Name), crypto.EncodeHash(collectivehash), v.Author, v.Epoch, "create collective"
			}
		}
		return "", "", v.Author, 0, ""
	case *actions.UpdateCollective:
		//fmt.Println("ucol")
		collectivehash := crypto.Hasher([]byte(v.OnBehalfOf))
		if collective, ok := i.state.Collectives[collectivehash]; ok {
			if status {
				return fmt.Sprintf("%v update", collective.Name), crypto.EncodeHash(collectivehash), v.Author, v.Epoch, "update collective"
			} else {
				handle := i.state.Members[crypto.HashToken(v.Author)]
				return fmt.Sprintf("%v proposed update for %v", handle, collective.Name), crypto.EncodeHash(collectivehash), v.Author, v.Epoch, "update collective"
			}
		}
		return "", "", v.Author, 0, ""
	case *actions.RequestMembership:
		//fmt.Println("memebr")
		collectivehash := crypto.Hasher([]byte(v.Collective))
		if collective, ok := i.state.Collectives[collectivehash]; ok {
			handle := i.state.Members[crypto.HashToken(v.Author)]
			if v.Include {
				if status {
					return fmt.Sprintf("%v became a member of %v", handle, collective.Name), crypto.EncodeHash(collectivehash), v.Author, v.Epoch, "request membership"
				} else {
					return fmt.Sprintf("%v requested membership to %v", handle, collective.Name), crypto.EncodeHash(collectivehash), v.Author, v.Epoch, "request membership"
				}
			}
			if status {
				return fmt.Sprintf("%v left %v", handle, collective.Name), crypto.EncodeHash(collectivehash), v.Author, v.Epoch, "request membership"
			}
		}
		return "", "", v.Author, 0, ""
	case *actions.RemoveMember:
		//fmt.Println("rmember")
		collectivehash := crypto.Hasher([]byte(v.OnBehalfOf))
		if collective, ok := i.state.Collectives[collectivehash]; ok {
			member := i.state.Members[crypto.HashToken(v.Member)]
			if status {
				return fmt.Sprintf("%v was removed from %v", member, collective.Name), crypto.EncodeHash(collectivehash), v.Author, v.Epoch, "remove member"
			} else {
				handle := i.state.Members[crypto.HashToken(v.Author)]
				return fmt.Sprintf("%v requested removal of %v from %v", handle, member, collective.Name), crypto.EncodeHash(collectivehash), v.Author, v.Epoch, "remove member"
			}
		}
		return "", "", v.Author, 0, ""
	case *actions.Signin:
		//fmt.Println("sign")
		authorhash := crypto.HashToken(v.Author)
		if _, ok := i.state.Members[authorhash]; ok {
			if status {
				return fmt.Sprintf("%v joined Synergy", v.Handle), "", v.Author, v.Epoch, "sign in"
			}
		}
		return "", "", v.Author, 0, ""
	}
	return "", "", crypto.ZeroToken, 0, ""
}
