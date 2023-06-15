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
	Name          string
	Description   string
	Majority      int
	SuperMajority int
}

func main() {
	// Define a template.
	const collectiveHTML = `
	<!DOCTYPE html>
	<html>
		<head>
			<meta charset="UTF-8">
			<title>Collective</title>
		</head>
		<body>
			<ul>
{{range .Collectives}}
<li>
			<h1>{{.Name}} Collective</h1>
			<br/>
			<h2>Description</h2>
			<p>{{.Description}}</p>
			<br/>
			<h2>Majority Policy</h2>
			<p>{{.Majority}}</p>
			<h2>Super Majority Policy</h2>
			<p>{{.SuperMajority}}</p>
</li>
{{end}}
			</ul>
		</body>
	</html>
`

	collectives := Collectives{Collectives: []Collective{
		{"First Collective", "first description", 10, 50},
		{"Random Collective", "random description", 100, 100},
		{"Default Collective", "default description", 1, 20},
	}}

	// Create a new template and parse the letter into it.
	t := template.Must(template.New("letter").Parse(collectiveHTML))

	// Execute the template for each recipient.
	err := t.Execute(os.Stdout, collectives)
	if err != nil {
		log.Println("executing template:", err)
	}
}
