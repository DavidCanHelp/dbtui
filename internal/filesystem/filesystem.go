package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// Entry represents a file or directory
type Entry struct {
	Name    string
	Path    string
	IsDir   bool
	Size    int64
	ModTime time.Time
	Mode    os.FileMode
}

// ReadDir reads the contents of a directory and returns sorted entries
// Directories are listed first, then files, both alphabetically
func ReadDir(path string) ([]Entry, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", path, err)
	}

	var result []Entry
	for _, entry := range entries {
		// Skip hidden files (starting with .)
		if len(entry.Name()) > 0 && entry.Name()[0] == '.' {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			// Skip entries we can't stat
			continue
		}

		result = append(result, Entry{
			Name:    entry.Name(),
			Path:    filepath.Join(path, entry.Name()),
			IsDir:   entry.IsDir(),
			Size:    info.Size(),
			ModTime: info.ModTime(),
			Mode:    info.Mode(),
		})
	}

	// Sort: directories first, then files, both alphabetically
	sort.Slice(result, func(i, j int) bool {
		if result[i].IsDir != result[j].IsDir {
			return result[i].IsDir
		}
		return result[i].Name < result[j].Name
	})

	return result, nil
}

// GetEntry returns information about a specific file or directory
func GetEntry(path string) (*Entry, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to stat %s: %w", path, err)
	}

	return &Entry{
		Name:    filepath.Base(path),
		Path:    path,
		IsDir:   info.IsDir(),
		Size:    info.Size(),
		ModTime: info.ModTime(),
		Mode:    info.Mode(),
	}, nil
}

// IsDir checks if the given path is a directory
func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// FormatSize formats a file size in human-readable format
func FormatSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}
