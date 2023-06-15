package main

import (
	"log"
	"os"
	"text/template"

	"github.com/lienkolabs/swell/crypto"
)

type Members struct {
	Members []Member
}

type Member struct {
	Hash   string
	Handle string
}

func main() {

	users := []string{"Ruben", "Larissa", "Odulfo"}
	members := Members{
		Members: make([]Member, 3),
	}
	for n, user := range users {
		hash := crypto.Hasher([]byte(user))
		text, _ := hash.MarshalText()
		members.Members[n] = Member{
			Handle: user,
			Hash:   string(text),
		}
	}

	t, err := template.ParseFiles("members.html")
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(os.Stdout, members)
}
