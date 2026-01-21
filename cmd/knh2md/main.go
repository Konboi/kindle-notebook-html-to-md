package main

import (
	"fmt"
	"os"

	knh2md "github.com/Konboi/kindle-notebook-html-to-md"
)

func main() {
	if err := knh2md.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
