# Manta - Agent Development Guide

This guide provides coding agents with essential information about the Manta pomodoro timer project.

## Project Overview

Manta is a lightweight, blazing-fast Pomodoro-style timer built with Go, using the Bubble Tea TUI framework.

**Tech Stack:**
- Language: Go 1.25.6
- TUI Framework: Charmbracelet Bubble Tea
- Audio: ebitengine/oto + hajimehoshi/go-mp3
- Module: `github.com/ihorbryk/manta`

**Project Structure:**
```
manta/
├── cmd/manta/          # Main entry point
├── internal/           # Internal packages (not exported)
│   ├── model.go       # Bubble Tea model & UI logic
│   ├── player.go      # Audio playback
│   ├── notify.go      # System notifications
│   └── tick.go        # Timer tick logic
└── assets/            # Static assets (audio files)
```

## Build & Run Commands

### Build
```bash
go build -o manta ./cmd/manta
```

### Run
```bash
go run ./cmd/manta
```

### Install
```bash
go install github.com/ihorbryk/manta/cmd/manta
```

### Format Code
```bash
gofmt -w .
# Or with simplification
gofmt -s -w .
```

### Lint (if golangci-lint is installed)
```bash
golangci-lint run
```

### Test
```bash
# Run all tests
go test ./...

# Run tests in a specific package
go test ./internal

# Run a single test
go test ./internal -run TestName

# Run with verbose output
go test -v ./...

# Run with coverage
go test -cover ./...
```

### Mod Management
```bash
# Tidy dependencies
go mod tidy

# Verify dependencies
go mod verify

# Download dependencies
go mod download
```

## Code Style Guidelines

### General Principles
- Follow standard Go conventions and idioms
- Use `gofmt` for formatting (tabs for indentation)
- Keep code simple and readable
- Prefer composition over inheritance

### Imports
- Group imports in this order:
  1. Standard library
  2. External dependencies
  3. Internal packages
- Separate groups with blank lines
- Use named imports only when necessary to avoid conflicts (e.g., `tea "github.com/charmbracelet/bubbletea"`)

**Example:**
```go
import (
    "fmt"
    "time"

    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"

    "github.com/ihorbryk/manta/internal"
)
```

### Naming Conventions
- **Packages:** Short, lowercase, single-word names (e.g., `internal`)
- **Exported identifiers:** PascalCase (e.g., `NewModel`, `PlayNotification`)
- **Unexported identifiers:** camelCase (e.g., `tickCmd`, `helpStyle`)
- **Constants:** ALL_CAPS for package-level constants (e.g., `WORKTIME`, `RESTTIME`)
- **Acronyms:** Keep consistent case (e.g., `ID`, `API`, `URL`)

### Types
- Use descriptive type names
- Prefer struct types over type aliases unless semantically meaningful
- Document exported types

**Example:**
```go
type model struct {
    progress progress.Model
    timeLeft int
    timeType string
    cursor   int
    choice   string
    pause    bool
    endTime  time.Time
}
```

### Functions
- Keep functions small and focused
- Document exported functions
- Return errors rather than panicking (except in truly unrecoverable situations)
- Use named return values sparingly (only when it improves clarity)

**Example:**
```go
func NewModel() model {
    return model{
        progress: progress.New(progress.WithDefaultGradient()),
        timeLeft: 0,
        timeType: WORKTIME,
    }
}
```

### Error Handling
- Always check errors; don't ignore them
- Return errors to the caller when possible
- Use `panic()` only for truly unrecoverable errors (e.g., initialization failures)
- Wrap errors with context when appropriate

**Example:**
```go
if _, err := tea.NewProgram(m).Run(); err != nil {
    fmt.Println("Oh no!", err)
    os.Exit(1)
}
```

### Comments
- Write comments for exported functions, types, and constants
- Use `//` for single-line comments
- Start comments with the identifier name for package-level declarations
- Keep comments concise and up-to-date

**Example:**
```go
// tickMsg represents a timer tick event
type tickMsg time.Time

// tickCmd returns a command that sends a tick message every second
func tickCmd() tea.Cmd {
    return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
        return tickMsg(t)
    })
}
```

### Bubble Tea Patterns
- Follow the Elm Architecture: Model, Update, View
- Use type switches in `Update()` for message handling
- Return appropriate commands from `Update()`
- Keep `View()` pure (no side effects)
- Use custom message types for domain events

### Constants
- Define constants at package level
- Group related constants together
- Use `const` blocks for multiple constants

**Example:**
```go
const (
    work     = 25 * 60
    rest     = 5 * 60
    WORKTIME = "work"
    RESTTIME = "rest"
)
```

## Testing Guidelines
- Place test files alongside source files with `_test.go` suffix
- Use table-driven tests for multiple test cases
- Follow the Arrange-Act-Assert pattern
- Use meaningful test names (e.g., `TestModelUpdate_KeyPress`)

## Key Development Notes
- The app uses Bubble Tea's Elm Architecture (Model-Update-View)
- Timer state is managed through the `model` struct
- Audio playback is synchronous (blocks until completion)
- System notifications use `terminal-notifier` (macOS specific)
- Main business logic is in `internal/` package

## Common Tasks
- **Adding a new timer mode:** Update `mapping`, `choices`, and handle in `Update()`
- **Changing timer durations:** Modify `work` and `rest` constants
- **Customizing UI:** Edit `View()` and lipgloss styles
- **Adding keyboard shortcuts:** Add cases in the `tea.KeyMsg` switch
