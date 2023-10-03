package main

import (
	"fmt"
	"os"
	"text/template"
)

const templLocal = `
{{.Title}}
=============================
{{range .Story}}
{{.}}

{{end}}
{{if .Options}}
-----------------------------
{{range $i, $elem := .Options}}
> {{$i | inc}}. {{ .Text }}
{{end}}
{{end}}
`

func local(data map[string]StoryArc) {
	tpl := template.Must(template.New("page").
		Funcs(template.FuncMap{"inc": inc}).
		Parse(templLocal))
	key := "intro"
	for {
		fmt.Printf("\033[2J\033[1;1H")
		tpl.Execute(os.Stdout, data[key])
		if len(data[key].Options) > 0 {
			var option int
			n, _ := fmt.Scanf("%d\n", &option)
			if n != 1 {
				continue
			}
			if len(data[key].Options) < option {
				continue
			}
			key = data[key].Options[option-1].Arc
		} else {
			break
		}
	}
}

func inc(i int) int {
	return i + 1
}
