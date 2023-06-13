package state

import (
	"errors"

	"github.com/lienkolabs/swell/crypto"
)

type MultipartMedia struct {
	Hash crypto.Hash
	Part byte
	Of   byte
	Data []byte
}

type PendingMedia struct {
	Hash          crypto.Hash
	NumberOfParts byte
	Parts         []*MultipartMedia
}

func (p *PendingMedia) Append(m *MultipartMedia) ([]byte, error) {
	if m.Of != p.NumberOfParts || m.Part > m.Of {
		return nil, errors.New("incompatible number of parts")
	}
	p.Parts[m.Part] = m
	size := 0
	for _, part := range p.Parts {
		if part == nil {
			return nil, nil
		}
		size += len(m.Data)
	}
	concanate := make([]byte, 0, size)
	for _, part := range p.Parts {
		concanate = append(concanate, part.Data...)
	}
	if crypto.Hasher(concanate) != p.Hash {
		return nil, errors.New("incompatible hash")
	}
	return concanate, nil
}

func (s *State) MediaPart(m *MultipartMedia) error {
	if _, ok := s.Media[m.Hash]; ok {
		return errors.New("media already exists")
	}
	if m.Of == 1 {
		s.Media[m.Hash] = m.Data
		return nil
	}
	if pending, ok := s.PendingMedia[m.Hash]; ok {
		media, err := pending.Append(m)
		if err != nil {
			return err
		}
		if media != nil {
			delete(s.PendingMedia, m.Hash)
			s.Media[m.Hash] = media
		}
		return nil
	}
	pending := PendingMedia{
		Hash:          m.Hash,
		NumberOfParts: m.Of,
		Parts:         make([]*MultipartMedia, int(m.Of)),
	}
	pending.Parts[m.Part] = m
	s.PendingMedia[m.Hash] = &pending
	return nil
}
