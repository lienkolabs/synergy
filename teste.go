package main

import (
	"os"
	"text/template"
)

const (
	a = `{{define "a"}} Ruben {{.}} {{end}} `
	b = `{{templae "a" Damião}}`
)

func main() {
	t := template.New("teste")
	t.Parse(a)
	t.Parse(b)
	t.Execute(os.Stdout, "")
}
