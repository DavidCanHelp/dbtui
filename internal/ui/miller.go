package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/davidcanhelp/dbtui/internal/filesystem"
)

const (
	columnWidth  = 30
	maxColumns   = 3 // Number of visible columns
	visibleItems = 20
)

var (
	columnStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, true, false, false).
			BorderForeground(lipgloss.Color("240")).
			Width(columnWidth).
			Padding(0, 1)

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("0")).
				Background(lipgloss.Color("12")).
				Bold(true)

	normalItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	directoryItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("12")).
				Bold(true)

	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("12")).
			Bold(true).
			Underline(true)
)

// Column represents a single column in the Miller columns view
type Column struct {
	Path    string
	Entries []filesystem.Entry
	Cursor  int // Selected item index
	Scroll  int // Scroll offset
}

// RenderMillerColumns renders the Miller columns view
func RenderMillerColumns(columns []Column, activeColumnIdx int, height int) string {
	// Determine which columns to show
	startCol := 0
	endCol := len(columns)

	// If we have more columns than maxColumns, show only the relevant ones
	if len(columns) > maxColumns {
		// Show columns around the active column
		startCol = activeColumnIdx - 1
		if startCol < 0 {
			startCol = 0
		}
		endCol = startCol + maxColumns
		if endCol > len(columns) {
			endCol = len(columns)
			startCol = endCol - maxColumns
			if startCol < 0 {
				startCol = 0
			}
		}
	}

	var renderedColumns []string

	for i := startCol; i < endCol; i++ {
		col := columns[i]
		isActive := i == activeColumnIdx

		rendered := renderColumn(col, isActive, height)
		renderedColumns = append(renderedColumns, rendered)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, renderedColumns...)
}

func renderColumn(col Column, isActive bool, height int) string {
	var lines []string

	// Calculate visible range
	start := col.Scroll
	end := start + visibleItems
	if end > len(col.Entries) {
		end = len(col.Entries)
	}

	// Render items
	for i := start; i < end; i++ {
		entry := col.Entries[i]
		isSelected := i == col.Cursor && isActive

		// Build the item text
		itemText := entry.Name
		if entry.IsDir {
			itemText += "/"
		}

		// Truncate if too long
		maxLen := columnWidth - 4
		if len(itemText) > maxLen {
			itemText = itemText[:maxLen-3] + "..."
		}

		// Apply style
		var styledItem string
		if isSelected {
			styledItem = selectedItemStyle.Render(itemText)
		} else if entry.IsDir {
			styledItem = directoryItemStyle.Render(itemText)
		} else {
			styledItem = normalItemStyle.Render(itemText)
		}

		lines = append(lines, styledItem)
	}

	// Add scroll indicators
	if col.Scroll > 0 {
		lines = append([]string{normalItemStyle.Render("↑ more...")}, lines...)
	}
	if end < len(col.Entries) {
		lines = append(lines, normalItemStyle.Render("↓ more..."))
	}

	// Pad to height
	for len(lines) < height-2 {
		lines = append(lines, "")
	}
	if len(lines) > height-2 {
		lines = lines[:height-2]
	}

	content := strings.Join(lines, "\n")

	// Add border
	style := columnStyle.Height(height)
	if isActive {
		style = style.BorderForeground(lipgloss.Color("12"))
	}

	return style.Render(content)
}

// RenderBreadcrumb renders the breadcrumb navigation showing the current path
func RenderBreadcrumb(pathParts []string, width int) string {
	if len(pathParts) == 0 {
		return headerStyle.Width(width).Render("/")
	}

	// Build breadcrumb
	breadcrumb := "/" + strings.Join(pathParts, " > ")

	// Truncate if too long
	if len(breadcrumb) > width-4 {
		breadcrumb = "..." + breadcrumb[len(breadcrumb)-(width-7):]
	}

	return headerStyle.Width(width).Render(breadcrumb)
}
