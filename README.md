# kindle-notebook-html-to-md

Convert Kindle **shared / exported notebook HTML** files into clean, structured **Markdown** files for Obsidian or any Markdown-based note system.

This tool is designed for users who export highlights and notes from the Kindle app using **“Export Notebook”** and want to manage them in tools like Obsidian, without relying on Kindle Cloud Reader or third-party sync plugins.

---

## Features

- 📘 Parse Kindle *Notebook Export* HTML files
- 📝 Convert highlights and notes to clean Markdown
- 🧱 Preserve chapter / section structure
- 📍 Keep location metadata
- 🎨 Preserve highlight colors (optional)
- 🗂 Obsidian-friendly output
- ⚡ Simple CLI interface

## Usage

```
$ knh2md input.html -o output.md
```

### Example for Obsidian vault

```
$ knh2md input.html -o ~/Obsidian/Kindle/engineering-leader.md
```


