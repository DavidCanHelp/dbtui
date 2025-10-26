package model

import (
	"fmt"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/davidcanhelp/dbtui/internal/filesystem"
	"github.com/davidcanhelp/dbtui/internal/ui"
)

type Model struct {
	rootPath  string      // Dropbox root path
	columns   []ui.Column // Column hierarchy
	activeCol int         // Currently active column index
	width     int         // Terminal width
	height    int         // Terminal height
	err       error       // Error state
	pathParts []string    // Path components for breadcrumb
	quitting  bool        // Whether the app is quitting
}

// NewModel creates a new model with the given Dropbox root path
func NewModel(rootPath string) Model {
	m := Model{
		rootPath:  rootPath,
		columns:   make([]ui.Column, 0),
		activeCol: 0,
		pathParts: make([]string, 0),
	}

	// Initialize with the root directory
	m.loadColumn(rootPath, 0)

	return m
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.quitting = true
			return m, tea.Quit

		case "up", "k":
			m.moveCursorUp()

		case "down", "j":
			m.moveCursorDown()

		case "left", "h":
			m.navigateUp()

		case "right", "l", "enter":
			m.navigateDown()
		}
	}

	return m, nil
}

// View renders the UI
func (m Model) View() string {
	if m.quitting {
		return "Goodbye!\n"
	}

	if m.err != nil {
		return fmt.Sprintf("Error: %v\nPress q to quit.\n", m.err)
	}

	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	// Calculate layout dimensions
	contentHeight := m.height - 3 // Leave room for breadcrumb and footer

	// Render breadcrumb
	breadcrumb := ui.RenderBreadcrumb(m.pathParts, m.width)

	// Render Miller columns
	millerColumns := ui.RenderMillerColumns(m.columns, m.activeCol, contentHeight)

	// Get selected entry for info panel
	var selectedEntry *filesystem.Entry
	if m.activeCol < len(m.columns) && m.columns[m.activeCol].Cursor < len(m.columns[m.activeCol].Entries) {
		selectedEntry = &m.columns[m.activeCol].Entries[m.columns[m.activeCol].Cursor]
	}

	// Render file info panel
	fileInfo := ui.RenderFileInfo(selectedEntry, contentHeight)

	// Combine columns and info panel
	mainContent := lipgloss.JoinHorizontal(lipgloss.Top, millerColumns, fileInfo)

	// Footer with help
	footer := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("↑↓: navigate  ←: back  →: forward  q: quit")

	// Combine all parts
	return lipgloss.JoinVertical(lipgloss.Left, breadcrumb, mainContent, footer)
}

// loadColumn loads directory entries for a given path at the specified column index
func (m *Model) loadColumn(path string, colIdx int) {
	entries, err := filesystem.ReadDir(path)
	if err != nil {
		m.err = err
		return
	}

	// Create the column
	col := ui.Column{
		Path:    path,
		Entries: entries,
		Cursor:  0,
		Scroll:  0,
	}

	// Update columns slice
	if colIdx >= len(m.columns) {
		m.columns = append(m.columns, col)
	} else {
		// Replace this column and remove everything after it
		m.columns = append(m.columns[:colIdx], col)
	}
}

// moveCursorUp moves the cursor up in the active column
func (m *Model) moveCursorUp() {
	if m.activeCol >= len(m.columns) {
		return
	}

	col := &m.columns[m.activeCol]
	if col.Cursor > 0 {
		col.Cursor--

		// Adjust scroll if needed
		if col.Cursor < col.Scroll {
			col.Scroll = col.Cursor
		}

		// Load preview if the new selection is a directory
		m.updatePreview()
	}
}

// moveCursorDown moves the cursor down in the active column
func (m *Model) moveCursorDown() {
	if m.activeCol >= len(m.columns) {
		return
	}

	col := &m.columns[m.activeCol]
	if col.Cursor < len(col.Entries)-1 {
		col.Cursor++

		// Adjust scroll if needed
		visibleItems := 20
		if col.Cursor >= col.Scroll+visibleItems {
			col.Scroll = col.Cursor - visibleItems + 1
		}

		// Load preview if the new selection is a directory
		m.updatePreview()
	}
}

// navigateUp moves to the parent directory
func (m *Model) navigateUp() {
	if m.activeCol > 0 {
		// Move to the previous column
		m.activeCol--

		// Remove columns after the active one
		m.columns = m.columns[:m.activeCol+1]

		// Update path parts
		if len(m.pathParts) > 0 {
			m.pathParts = m.pathParts[:len(m.pathParts)-1]
		}

		// Update preview
		m.updatePreview()
	}
}

// navigateDown enters the selected directory
func (m *Model) navigateDown() {
	if m.activeCol >= len(m.columns) {
		return
	}

	col := &m.columns[m.activeCol]
	if col.Cursor >= len(col.Entries) {
		return
	}

	entry := col.Entries[col.Cursor]

	// Only navigate if it's a directory
	if !entry.IsDir {
		return
	}

	// Check if preview column already exists
	if m.activeCol+1 < len(m.columns) {
		// Just move to it
		m.activeCol++
	} else {
		// Load the directory
		m.loadColumn(entry.Path, m.activeCol+1)
		m.activeCol++
	}

	// Update path parts
	m.pathParts = append(m.pathParts, entry.Name)

	// Load preview of the new selection
	m.updatePreview()
}

// updatePreview loads a preview column for the currently selected directory
func (m *Model) updatePreview() {
	if m.activeCol >= len(m.columns) {
		return
	}

	col := &m.columns[m.activeCol]
	if col.Cursor >= len(col.Entries) {
		return
	}

	entry := col.Entries[col.Cursor]

	// Only preview directories
	if !entry.IsDir {
		// Remove any preview column
		if m.activeCol+1 < len(m.columns) {
			m.columns = m.columns[:m.activeCol+1]
		}
		return
	}

	// Load preview column
	m.loadColumn(entry.Path, m.activeCol+1)
}

// GetCurrentPath returns the current full path
func (m *Model) GetCurrentPath() string {
	if len(m.pathParts) == 0 {
		return m.rootPath
	}
	return filepath.Join(m.rootPath, strings.Join(m.pathParts, string(filepath.Separator)))
}
