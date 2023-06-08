package main

import (
	"fmt"
	"log"
	"os"

	"github.com/lienkolabs/swell/crypto"
)

type Epoch int

func (e Epoch) Epoch() uint64 {
	return 0
}

func (e Epoch) Token() crypto.Token {
	return crypto.ZeroToken
}

func main() {
	jsonFile := os.Args[1]
	json, err := os.ReadFile(jsonFile)
	if err != nil {
		log.Fatal(err)
	}
	draft := ParseCreateDraft(json)
	instruct, _ := InstructDraft(draft, Epoch(0))
	fmt.Printf("%#v\n", DraftInstructionToJSON(instruct))
}
