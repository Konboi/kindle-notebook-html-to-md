package knh2md

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
)

const (
	appName        = "knh2md"
	appDescription = "Convert Kindle Notebook HTML to Markdown"
)

// CLI represents the command-line interface
type CLI struct {
	Input    string `arg:"" help:"Input HTML file path" type:"existingfile"`
	Output   string `short:"o" help:"Output file path (default: stdout)"`
	Template string `short:"t" help:"Custom template file path" type:"existingfile"`
}

// Run executes the CLI command
func (c *CLI) Run() error {
	// Open input file
	f, err := os.Open(c.Input)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer f.Close()

	// Parse HTML
	notebook, err := Parse(f)
	if err != nil {
		return fmt.Errorf("failed to parse HTML: %w", err)
	}

	// Convert to Markdown
	var markdown string
	if c.Template != "" {
		tmplContent, err := os.ReadFile(c.Template)
		if err != nil {
			return fmt.Errorf("failed to read template file: %w", err)
		}
		markdown, err = ConvertWithTemplate(notebook, string(tmplContent))
		if err != nil {
			return fmt.Errorf("failed to convert with template: %w", err)
		}
	} else {
		markdown = Convert(notebook)
	}

	// Output
	if c.Output == "" {
		// Write to stdout
		fmt.Print(markdown)
	} else {
		// Write to file
		if err := os.WriteFile(c.Output, []byte(markdown), 0644); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
	}

	return nil
}

// Execute parses command-line arguments and runs the CLI
func Execute() error {
	var cli CLI
	kong.Parse(&cli,
		kong.Name(appName),
		kong.Description(appDescription),
		kong.UsageOnError(),
	)
	return cli.Run()
}
