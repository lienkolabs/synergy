package main

import (
	"log"
	"os"
	"text/template"

	"github.com/lienkolabs/swell/crypto"
)

type Drafts struct {
	Drafts []Draft
}

type Draft struct {
	Title       string
	Author      string
	CoAuthors   []string
	Description string
	Keywords    []string
	Hash        crypto.Hash
}

func main() {

	titles := []string{"DraftA", "DraftB", "DraftC", "DraftD"}
	authors := []string{"autorA", "autorB", "autorC", "autorD"}
	coauthors := [][]string{[]string{"aa", "ab"}, []string{"ba", "bb", "bc"}, []string{"ca"}, []string{"da,db"}}
	descriptions := []string{"descriptionA", "descriptionB", "descriptionC", "descriptionD"}
	keywords := [][]string{[]string{"palavra", "palavraa"}, []string{"paaaaliavra"}, []string{}, []string{}}

	drafts := Drafts{
		Drafts: make([]Draft, 4),
	}

	for n, title := range titles {
		drafts.Drafts[n] = Draft{
			Title:       title,
			Author:      authors[n],
			CoAuthors:   coauthors[n],
			Description: descriptions[n],
			Keywords:    keywords[n],
			Hash:        crypto.Hasher([]byte(title)),
		}
	}

	t, err := template.ParseFiles("drafts.html")
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(os.Stdout, drafts)
}
