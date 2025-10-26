package dropbox

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Info represents the structure of Dropbox's info.json file
type Info struct {
	Personal *Account `json:"personal,omitempty"`
	Business *Account `json:"business,omitempty"`
}

// Account represents a Dropbox account configuration
type Account struct {
	Path             string `json:"path"`
	Host             int64  `json:"host"`
	IsTeam           bool   `json:"is_team"`
	SubscriptionType string `json:"subscription_type"`
}

// DetectDropboxPath finds the Dropbox folder path by reading ~/.dropbox/info.json
// It prioritizes personal account over business account if both exist
func DetectDropboxPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	infoPath := filepath.Join(homeDir, ".dropbox", "info.json")

	data, err := os.ReadFile(infoPath)
	if err != nil {
		return "", fmt.Errorf("failed to read Dropbox info file at %s: %w (is Dropbox installed?)", infoPath, err)
	}

	var info Info
	if err := json.Unmarshal(data, &info); err != nil {
		return "", fmt.Errorf("failed to parse Dropbox info file: %w", err)
	}

	// Prefer personal account, fall back to business
	if info.Personal != nil && info.Personal.Path != "" {
		return info.Personal.Path, nil
	}

	if info.Business != nil && info.Business.Path != "" {
		return info.Business.Path, nil
	}

	return "", fmt.Errorf("no Dropbox path found in info file")
}
