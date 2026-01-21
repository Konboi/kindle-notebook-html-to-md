# knh2md

A CLI tool to convert Kindle Notebook HTML to Obsidian-compatible Markdown.

Converts HTML files exported from the Kindle app's "Export Notebook" feature into Markdown format for easy management in Obsidian and other note-taking tools.

## Installation

```bash
go install github.com/Konboi/kindle-notebook-html-to-md/cmd/knh2md@latest
```

Or clone and build:

```bash
git clone https://github.com/Konboi/kindle-notebook-html-to-md.git
cd kindle-notebook-html-to-md
go build -o knh2md ./cmd/knh2md
```

## Usage

```bash
# Output to stdout
knh2md input.html

# Output to file
knh2md input.html -o output.md

# Output directly to Obsidian vault
knh2md input.html -o ~/Obsidian/Kindle/book-name.md

# Use custom template
knh2md input.html -t custom.tmpl -o output.md
```

### Options

| Option | Short | Description |
|--------|-------|-------------|
| `--output` | `-o` | Output file path (default: stdout) |
| `--template` | `-t` | Custom template file path |

## How to Export Kindle Notebook HTML

1. Open a book in the Kindle app
2. Open the highlights/notes view
3. Select "Export" or "Share" → "Export Notebook"
4. Save as HTML format

## Output Format

```markdown
# Book Title

## Metadata
* Author: [[Author Name]]

## Highlights

### Chapter Heading

> 🟨 Highlight text

- Location: 165

---

> 📝 Note text

- Location: 200

---
```

## Highlight Colors and Emojis

| Color  | Emoji |
|--------|-------|
| yellow | 🟨    |
| pink   | 🟥    |
| orange | 🟧    |
| green  | 🟩    |
| aqua   | 🟦    |
| note   | 📝    |

## Custom Templates

You can customize the output format using Go's [text/template](https://pkg.go.dev/text/template) syntax.

### Template Variables

| Variable | Description |
|----------|-------------|
| `{{.Title}}` | Book title |
| `{{.Author}}` | Author name |
| `{{.Sections}}` | Array of sections |

#### Section Variables (within `{{range .Sections}}`)

| Variable | Description |
|----------|-------------|
| `{{.Heading}}` | Section/chapter heading |
| `{{.Notes}}` | Array of notes in this section |

#### Note Variables (within `{{range .Notes}}`)

| Variable | Description |
|----------|-------------|
| `{{.Type}}` | "highlight" or "note" |
| `{{.Color}}` | Highlight color (yellow, pink, orange, green, aqua) |
| `{{.Location}}` | Kindle location number |
| `{{.Text}}` | Highlight/note text content |
| `{{emoji .}}` | Get emoji for the note (helper function) |

### Example Custom Template

```
# {{.Title}}

Author: {{.Author}}
{{range .Sections}}
## {{.Heading}}
{{range .Notes}}
- {{emoji .}} {{.Text}} (Location: {{.Location}})
{{end}}{{end}}
```

### Default Template

```
# {{.Title}}

## Metadata
* Author: [[{{.Author}}]]

## Highlights
{{range .Sections}}
### {{.Heading}}
{{range .Notes}}
> {{emoji .}} {{.Text}}

- Location: {{.Location}}

---
{{end}}{{end}}
```

## License

MIT
