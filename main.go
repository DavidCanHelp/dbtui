package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/davidcanhelp/dbtui/internal/dropbox"
	"github.com/davidcanhelp/dbtui/internal/model"
)

const version = "1.0.0"

func main() {
	// Handle version flag
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Printf("dbtui version %s\n", version)
		os.Exit(0)
	}

	// Detect Dropbox path
	dropboxPath, err := dropbox.DetectDropboxPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintln(os.Stderr, "\nPlease ensure Dropbox is installed and set up on your system.")
		os.Exit(1)
	}

	// Create the model
	m := model.NewModel(dropboxPath)

	// Create the Bubble Tea program
	p := tea.NewProgram(m, tea.WithAltScreen())

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
