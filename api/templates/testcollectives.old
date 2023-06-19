package main

import (
	"log"
	"os"
	"text/template"
)

type Collectives struct {
	Collectives []Collective
}

type Collective struct {
	Name        string
	Description string
}

func main() {

	names := []string{"ColetivoA", "ColetivoB", "ColetivoC", "ColetivoD"}
	descriptions := []string{"primeiro coletivo", "segundo coletivo", "terceiro coletivo", "quarto coletivo"}

	collectives := Collectives{
		Collectives: make([]Collective, 4),
	}

	for n, nome := range names {
		collectives.Collectives[n] = Collective{
			Name:        nome,
			Description: descriptions[n],
		}
	}

	t, err := template.ParseFiles("collectives.html")
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(os.Stdout, collectives)
}
