package knh2md

import (
	"io"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// Notebook represents the entire Kindle notebook export
type Notebook struct {
	Title    string
	Author   string
	Sections []Section
}

// Section represents a chapter or section in the book
type Section struct {
	Heading string
	Notes   []Note
}

// Note represents a highlight or note
type Note struct {
	Type     string // "highlight" or "note"
	Color    string // yellow, blue, pink, orange, green, aqua
	Location string
	Text     string
}

// Parse parses the Kindle Notebook HTML and returns a Notebook struct
func Parse(r io.Reader) (*Notebook, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	notebook := &Notebook{}
	var currentSection *Section
	var pendingNote *Note

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "div" {
			class := getAttr(n, "class")
			switch class {
			case "bookTitle":
				notebook.Title = strings.TrimSpace(getTextContent(n))
			case "authors":
				notebook.Author = strings.TrimSpace(getTextContent(n))
			case "sectionHeading":
				// Save pending note to current section before starting new section
				if pendingNote != nil && currentSection != nil {
					currentSection.Notes = append(currentSection.Notes, *pendingNote)
					pendingNote = nil
				}
				// Start a new section
				if currentSection != nil {
					notebook.Sections = append(notebook.Sections, *currentSection)
				}
				currentSection = &Section{
					Heading: strings.TrimSpace(getTextContent(n)),
					Notes:   []Note{},
				}
			case "noteHeading":
				// Save previous pending note
				if pendingNote != nil && currentSection != nil {
					currentSection.Notes = append(currentSection.Notes, *pendingNote)
				}
				// Parse the note heading
				pendingNote = parseNoteHeading(n)
			case "noteText":
				if pendingNote != nil {
					pendingNote.Text = strings.TrimSpace(getTextContent(n))
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(doc)

	// Save the last pending note and section
	if pendingNote != nil && currentSection != nil {
		currentSection.Notes = append(currentSection.Notes, *pendingNote)
	}
	if currentSection != nil {
		notebook.Sections = append(notebook.Sections, *currentSection)
	}

	return notebook, nil
}

// getAttr returns the value of the specified attribute
func getAttr(n *html.Node, key string) string {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

// getTextContent returns the text content of a node and its children
func getTextContent(n *html.Node) string {
	var sb strings.Builder
	var extractText func(*html.Node)
	extractText = func(n *html.Node) {
		if n.Type == html.TextNode {
			sb.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractText(c)
		}
	}
	extractText(n)
	return sb.String()
}

// parseNoteHeading parses the noteHeading div and extracts type, color, and location
func parseNoteHeading(n *html.Node) *Note {
	note := &Note{}
	text := strings.TrimSpace(getTextContent(n))

	// Check if it's a Note or Highlight
	if strings.HasPrefix(text, "Note") {
		note.Type = "note"
	} else if strings.HasPrefix(text, "Highlight") {
		note.Type = "highlight"
		// Extract color from span class
		note.Color = extractColor(n)
	}

	// Extract location
	locRegex := regexp.MustCompile(`Location\s+(\d+)`)
	if matches := locRegex.FindStringSubmatch(text); len(matches) > 1 {
		note.Location = matches[1]
	}

	return note
}

// extractColor finds the highlight_* class in child span elements
func extractColor(n *html.Node) string {
	var color string
	var findColor func(*html.Node)
	findColor = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "span" {
			class := getAttr(n, "class")
			if strings.HasPrefix(class, "highlight_") {
				color = strings.TrimPrefix(class, "highlight_")
				return
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findColor(c)
		}
	}
	findColor(n)
	return color
}
