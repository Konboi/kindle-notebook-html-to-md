package knh2md

import (
	"strings"
	"testing"
)

func TestConvert(t *testing.T) {
	notebook := &Notebook{
		Title:  "Test Book",
		Author: "Test Author",
		Sections: []Section{
			{
				Heading: "Chapter 1",
				Notes: []Note{
					{
						Type:     "highlight",
						Color:    "yellow",
						Location: "100",
						Text:     "Yellow highlight text",
					},
					{
						Type:     "note",
						Location: "150",
						Text:     "Note text",
					},
				},
			},
		},
	}

	result := Convert(notebook)

	// Check title
	if !strings.Contains(result, "# Test Book\n") {
		t.Error("Result should contain title as H1")
	}

	// Check metadata
	if !strings.Contains(result, "## Metadata\n") {
		t.Error("Result should contain Metadata heading")
	}
	if !strings.Contains(result, "* Author: [[Test Author]]") {
		t.Error("Result should contain author in wikilink format")
	}

	// Check highlights heading
	if !strings.Contains(result, "## Highlights\n") {
		t.Error("Result should contain Highlights heading")
	}

	// Check section heading
	if !strings.Contains(result, "### Chapter 1\n") {
		t.Error("Result should contain section as H3")
	}

	// Check yellow highlight with emoji
	if !strings.Contains(result, "> 🟨 Yellow highlight text") {
		t.Error("Result should contain yellow highlight with emoji")
	}

	// Check note with emoji
	if !strings.Contains(result, "> 📝 Note text") {
		t.Error("Result should contain note with emoji")
	}

	// Check location
	if !strings.Contains(result, "- Location: 100") {
		t.Error("Result should contain location")
	}

	// Check separator
	if !strings.Contains(result, "---") {
		t.Error("Result should contain separator")
	}
}

func TestConvertAllColors(t *testing.T) {
	tests := []struct {
		color string
		emoji string
	}{
		{"yellow", "🟨"},
		{"pink", "🟥"},
		{"orange", "🟧"},
		{"green", "🟩"},
		{"aqua", "🟦"},
	}

	for _, tt := range tests {
		notebook := &Notebook{
			Title:  "Book",
			Author: "Author",
			Sections: []Section{
				{
					Heading: "Section",
					Notes: []Note{
						{
							Type:     "highlight",
							Color:    tt.color,
							Location: "100",
							Text:     "Text",
						},
					},
				},
			},
		}

		result := Convert(notebook)

		if !strings.Contains(result, "> "+tt.emoji+" Text") {
			t.Errorf("Color %s should produce emoji %s", tt.color, tt.emoji)
		}
	}
}

func TestConvertUnknownColorDefaultsToYellow(t *testing.T) {
	notebook := &Notebook{
		Title:  "Book",
		Author: "Author",
		Sections: []Section{
			{
				Heading: "Section",
				Notes: []Note{
					{
						Type:     "highlight",
						Color:    "unknown",
						Location: "100",
						Text:     "Text",
					},
				},
			},
		},
	}

	result := Convert(notebook)

	if !strings.Contains(result, "> 🟨 Text") {
		t.Error("Unknown color should default to yellow emoji")
	}
}

func TestConvertEmptyNotebook(t *testing.T) {
	notebook := &Notebook{
		Title:    "Empty Book",
		Author:   "Author",
		Sections: []Section{},
	}

	result := Convert(notebook)

	if !strings.Contains(result, "# Empty Book") {
		t.Error("Result should contain title")
	}
	if !strings.Contains(result, "## Highlights") {
		t.Error("Result should contain Highlights heading even if empty")
	}
}

func TestColorEmojiMap(t *testing.T) {
	expected := map[string]string{
		"yellow": "🟨",
		"pink":   "🟥",
		"orange": "🟧",
		"green":  "🟩",
		"aqua":   "🟦",
	}

	for color, emoji := range expected {
		if ColorEmoji[color] != emoji {
			t.Errorf("ColorEmoji[%s] = %s, want %s", color, ColorEmoji[color], emoji)
		}
	}
}

func TestNoteEmoji(t *testing.T) {
	if NoteEmoji != "📝" {
		t.Errorf("NoteEmoji = %s, want 📝", NoteEmoji)
	}
}

func TestConvertWithTemplate(t *testing.T) {
	notebook := &Notebook{
		Title:  "Test Book",
		Author: "Test Author",
		Sections: []Section{
			{
				Heading: "Chapter 1",
				Notes: []Note{
					{
						Type:     "highlight",
						Color:    "yellow",
						Location: "100",
						Text:     "Highlight text",
					},
				},
			},
		},
	}

	customTemplate := `Title: {{.Title}}
Author: {{.Author}}
{{range .Sections}}## {{.Heading}}
{{range .Notes}}[{{emoji .}}] {{.Text}} ({{.Location}})
{{end}}{{end}}`

	result, err := ConvertWithTemplate(notebook, customTemplate)
	if err != nil {
		t.Fatalf("ConvertWithTemplate failed: %v", err)
	}

	if !strings.Contains(result, "Title: Test Book") {
		t.Error("Result should contain custom title format")
	}
	if !strings.Contains(result, "Author: Test Author") {
		t.Error("Result should contain custom author format")
	}
	if !strings.Contains(result, "[🟨] Highlight text (100)") {
		t.Error("Result should contain custom note format")
	}
}

func TestConvertWithTemplateError(t *testing.T) {
	notebook := &Notebook{Title: "Test"}

	_, err := ConvertWithTemplate(notebook, "{{.InvalidField}}")
	if err == nil {
		t.Error("Expected error for invalid template field")
	}
}
