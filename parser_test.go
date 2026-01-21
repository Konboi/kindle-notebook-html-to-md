package knh2md

import (
	"os"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	f, err := os.Open("testdata/sample-notebook.html")
	if err != nil {
		t.Fatalf("Failed to open test file: %v", err)
	}
	defer f.Close()

	notebook, err := Parse(f)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Test title
	if notebook.Title != "Test Book Title" {
		t.Errorf("Title = %q, want %q", notebook.Title, "Test Book Title")
	}

	// Test author
	if notebook.Author != "Test Author" {
		t.Errorf("Author = %q, want %q", notebook.Author, "Test Author")
	}

	// Test sections count
	if len(notebook.Sections) != 2 {
		t.Fatalf("Sections count = %d, want 2", len(notebook.Sections))
	}

	// Test first section
	section1 := notebook.Sections[0]
	if section1.Heading != "Chapter 1" {
		t.Errorf("Section1.Heading = %q, want %q", section1.Heading, "Chapter 1")
	}
	if len(section1.Notes) != 2 {
		t.Fatalf("Section1.Notes count = %d, want 2", len(section1.Notes))
	}

	// Test yellow highlight
	note1 := section1.Notes[0]
	if note1.Type != "highlight" {
		t.Errorf("Note1.Type = %q, want %q", note1.Type, "highlight")
	}
	if note1.Color != "yellow" {
		t.Errorf("Note1.Color = %q, want %q", note1.Color, "yellow")
	}
	if note1.Location != "100" {
		t.Errorf("Note1.Location = %q, want %q", note1.Location, "100")
	}
	if note1.Text != "This is a yellow highlight." {
		t.Errorf("Note1.Text = %q, want %q", note1.Text, "This is a yellow highlight.")
	}

	// Test note
	note2 := section1.Notes[1]
	if note2.Type != "note" {
		t.Errorf("Note2.Type = %q, want %q", note2.Type, "note")
	}
	if note2.Location != "150" {
		t.Errorf("Note2.Location = %q, want %q", note2.Location, "150")
	}
	if note2.Text != "This is a note." {
		t.Errorf("Note2.Text = %q, want %q", note2.Text, "This is a note.")
	}

	// Test second section
	section2 := notebook.Sections[1]
	if section2.Heading != "Chapter 2" {
		t.Errorf("Section2.Heading = %q, want %q", section2.Heading, "Chapter 2")
	}
	if len(section2.Notes) != 1 {
		t.Fatalf("Section2.Notes count = %d, want 1", len(section2.Notes))
	}

	// Test aqua highlight
	note3 := section2.Notes[0]
	if note3.Type != "highlight" {
		t.Errorf("Note3.Type = %q, want %q", note3.Type, "highlight")
	}
	if note3.Color != "aqua" {
		t.Errorf("Note3.Color = %q, want %q", note3.Color, "aqua")
	}
}

func TestParseAllColors(t *testing.T) {
	colors := []string{"yellow", "pink", "orange", "green", "aqua"}

	for _, color := range colors {
		html := `
<html>
<body>
<div class="bookTitle">Book</div>
<div class="authors">Author</div>
<div class="sectionHeading">Section</div>
<div class="noteHeading">
	Highlight(<span class="highlight_` + color + `">` + color + `</span>) - Location 100
</div>
<div class="noteText">Text</div>
</body>
</html>
`
		notebook, err := Parse(strings.NewReader(html))
		if err != nil {
			t.Fatalf("Parse failed for color %s: %v", color, err)
		}

		if len(notebook.Sections) != 1 || len(notebook.Sections[0].Notes) != 1 {
			t.Fatalf("Expected 1 section with 1 note for color %s", color)
		}

		note := notebook.Sections[0].Notes[0]
		if note.Color != color {
			t.Errorf("Color = %q, want %q", note.Color, color)
		}
	}
}

func TestParseEmptyHTML(t *testing.T) {
	html := `<html><body></body></html>`

	notebook, err := Parse(strings.NewReader(html))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if notebook.Title != "" {
		t.Errorf("Title = %q, want empty", notebook.Title)
	}
	if notebook.Author != "" {
		t.Errorf("Author = %q, want empty", notebook.Author)
	}
	if len(notebook.Sections) != 0 {
		t.Errorf("Sections count = %d, want 0", len(notebook.Sections))
	}
}
