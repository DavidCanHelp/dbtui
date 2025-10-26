package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/davidcanhelp/dbtui/internal/filesystem"
)

var (
	infoBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1).
			Width(40)

	infoLabelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("12")).
			Bold(true)

	infoValueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))
)

// RenderFileInfo renders the file information panel for the given entry
func RenderFileInfo(entry *filesystem.Entry, height int) string {
	if entry == nil {
		return infoBoxStyle.Height(height).Render("No selection")
	}

	var lines []string

	// Name
	lines = append(lines, fmt.Sprintf("%s %s",
		infoLabelStyle.Render("Name:"),
		infoValueStyle.Render(entry.Name)))

	// Type
	fileType := "File"
	if entry.IsDir {
		fileType = "Directory"
	}
	lines = append(lines, fmt.Sprintf("%s %s",
		infoLabelStyle.Render("Type:"),
		infoValueStyle.Render(fileType)))

	// Size (only for files)
	if !entry.IsDir {
		lines = append(lines, fmt.Sprintf("%s %s",
			infoLabelStyle.Render("Size:"),
			infoValueStyle.Render(filesystem.FormatSize(entry.Size))))
	}

	// Modified time
	lines = append(lines, fmt.Sprintf("%s %s",
		infoLabelStyle.Render("Modified:"),
		infoValueStyle.Render(entry.ModTime.Format("Jan 02, 2006 15:04"))))

	// Permissions
	lines = append(lines, fmt.Sprintf("%s %s",
		infoLabelStyle.Render("Permissions:"),
		infoValueStyle.Render(entry.Mode.String())))

	// Path
	lines = append(lines, "")
	lines = append(lines, infoLabelStyle.Render("Path:"))
	lines = append(lines, infoValueStyle.Render(entry.Path))

	content := strings.Join(lines, "\n")
	return infoBoxStyle.Height(height).Render(content)
}
