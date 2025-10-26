# DBTUI - Dropbox Terminal User Interface

[![Version](https://img.shields.io/badge/version-1.0.0-blue.svg)](https://github.com/davidcanhelp/dbtui/releases)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

A macOS Finder-style Miller columns file navigator for Dropbox, built with Go and Bubble Tea.

## Features

- **Miller Columns Navigation**: Browse your Dropbox files using the familiar multi-column view from macOS Finder
- **Keyboard-Driven**: Navigate entirely with arrow keys or vim-style hjkl keys
- **File Information Panel**: View detailed metadata for selected files and directories
- **Auto-Detection**: Automatically finds your Dropbox folder location
- **Clean Terminal UI**: Beautiful, responsive interface built with Bubble Tea and Lipgloss

## Prerequisites

- Go 1.21 or later
- Dropbox installed and configured on your macOS system

## Installation

### Option 1: Install with Go

```bash
go install github.com/davidcanhelp/dbtui@latest
```

### Option 2: Build from Source

```bash
# Clone the repository
git clone https://github.com/davidcanhelp/dbtui.git
cd dbtui

# Build the application
go build -o dbtui

# Optional: Move to your PATH
sudo mv dbtui /usr/local/bin/
```

### Option 3: Download Binary

Download the latest release from the [releases page](https://github.com/davidcanhelp/dbtui/releases).

## Usage

Simply run the application:

```bash
dbtui
```

The app will automatically detect your Dropbox folder location from `~/.dropbox/info.json`.

To check the version:

```bash
dbtui --version
```

### Keyboard Controls

| Key | Action |
|-----|--------|
| `↑` or `k` | Move selection up |
| `↓` or `j` | Move selection down |
| `→` or `l` or `Enter` | Enter directory (navigate forward) |
| `←` or `h` | Go to parent directory (navigate back) |
| `q` or `Esc` or `Ctrl+C` | Quit application |

### Navigation

- Use **up/down arrows** to move the selection within the current column
- Use **right arrow** to enter a directory and move deeper in the hierarchy
- Use **left arrow** to go back to the parent directory
- The **file info panel** on the right shows details about the currently selected item

## How It Works

1. **Dropbox Detection**: Reads your Dropbox configuration from `~/.dropbox/info.json` to locate your Dropbox folder
2. **Miller Columns**: Displays up to 3 columns at a time, showing your current location, its contents, and a preview of the selected item
3. **Real-time Navigation**: As you navigate, the columns update to show the directory hierarchy
4. **File Information**: The right panel displays size, modification date, permissions, and full path

## Project Structure

```
dbtui/
├── main.go                      # Application entry point
├── internal/
│   ├── dropbox/
│   │   └── detector.go          # Dropbox path auto-detection
│   ├── model/
│   │   └── model.go             # Bubble Tea application model
│   ├── ui/
│   │   ├── miller.go            # Miller columns rendering
│   │   └── fileinfo.go          # File information panel
│   └── filesystem/
│       └── filesystem.go        # File system operations
└── README.md
```

## License

See [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
