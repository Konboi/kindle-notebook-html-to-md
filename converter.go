package knh2md

import (
	"bytes"
	"text/template"
)

// ColorEmoji maps highlight colors to emoji
var ColorEmoji = map[string]string{
	"yellow": "🟨",
	"pink":   "🟥",
	"orange": "🟧",
	"green":  "🟩",
	"aqua":   "🟦",
}

// NoteEmoji is the emoji used for notes
const NoteEmoji = "📝"

// DefaultTemplate is the default Markdown template
const DefaultTemplate = `# {{.Title}}

## Metadata
* Author: [[{{.Author}}]]

## Highlights
{{range .Sections}}
### {{.Heading}}
{{range .Notes}}
> {{emoji .}} {{.Text}}

- Location: {{.Location}}

---
{{end}}{{end}}`

// templateFuncs provides helper functions for templates
var templateFuncs = template.FuncMap{
	"emoji": func(note Note) string {
		if note.Type == "note" {
			return NoteEmoji
		}
		if e, ok := ColorEmoji[note.Color]; ok {
			return e
		}
		return "🟨" // default to yellow
	},
}

// Convert converts a Notebook to Markdown using the default template
func Convert(nb *Notebook) string {
	result, _ := ConvertWithTemplate(nb, DefaultTemplate)
	return result
}

// ConvertWithTemplate converts a Notebook to Markdown using a custom template
func ConvertWithTemplate(nb *Notebook, tmplStr string) (string, error) {
	tmpl, err := template.New("notebook").Funcs(templateFuncs).Parse(tmplStr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, nb); err != nil {
		return "", err
	}

	return buf.String(), nil
}
