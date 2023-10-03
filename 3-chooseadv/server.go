package main

import (
	"html/template"
	"log"
	"net/http"
)

const templ = `<!DOCTYPE html>
<html>
<head><title>Choose your own adventure</title></head>
<body>
<div class="content">
<h1>{{.Title}}</h1>
{{range .Story}}
<p>{{.}}</p>
{{end}}
{{if .Options}}
<hr>
{{range .Options}}
<p><a href="/{{.Arc}}">{{.Text}}</a></p>
{{end}}
{{end}}
</div>
<style>
.content {
	max-width: 600px;
	margin: auto;
  }
</style>
</body>
</html>
`

type pageHandler struct {
	data map[string]StoryArc
	tpl  *template.Template
}

func (h *pageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/" {
		h.tpl.Execute(w, h.data["intro"])
		return
	} else {
		s := r.RequestURI[1:]
		if _, ok := h.data[s]; ok {
			h.tpl.Execute(w, h.data[s])
			return
		}
	}
	http.NotFound(w, r)
}

func startServer(data map[string]StoryArc) {
	pageHandler := new(pageHandler)
	pageHandler.data = data
	pageHandler.tpl = template.Must(template.New("page").Parse(templ))
	http.Handle("/", pageHandler)
	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}
