package api

import (
	"os"

	"github.com/lienkolabs/swell/crypto"
)

const maxFileSize = 10000

type TruncatedFile struct {
	Hash  crypto.Hash
	Parts [][]byte
}

func splitBytes(bytes []byte) *TruncatedFile {
	truncated := TruncatedFile{
		Hash:  crypto.Hasher(bytes),
		Parts: make([][]byte, len(bytes)/maxFileSize+1),
	}
	for n := 0; n < len(truncated.Parts); n++ {
		if (n+1)*maxFileSize >= len(bytes) {
			truncated.Parts[n] = bytes[n*maxFileSize:]
		} else {
			truncated.Parts[n] = bytes[n*maxFileSize : (n+1)*maxFileSize]
		}

	}
	return &truncated
}

func loadFile(filename string) (*TruncatedFile, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	truncated := TruncatedFile{
		Hash:  crypto.Hasher(bytes),
		Parts: make([][]byte, len(bytes)/maxFileSize+1),
	}
	for n := 0; n < len(truncated.Parts); n++ {
		if (n+1)*maxFileSize >= len(bytes) {
			truncated.Parts[n] = bytes[n*maxFileSize:]
		} else {
			truncated.Parts[n] = bytes[n*maxFileSize : (n+1)*maxFileSize]
		}

	}
	return &truncated, nil
}
