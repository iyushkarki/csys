# csys - System Monitoring CLI

> Beautiful, developer-friendly system monitoring & git shortcuts for Mac & Linux

**Also available as `cs` for faster typing!**

A lightweight CLI tool that gives you instant, beautiful insights into your system's health. No more cryptic `df` output or hunting through Activity Monitor - just clean, readable information about your disk, memory, CPU, and more.

## Preview

<img width="501" height="288" alt="image" src="https://github.com/user-attachments/assets/2165c2a3-b31a-428f-b4c9-183c216d0918" />

## Features

**System Monitoring (Phase 1)**

- Beautiful system overview at a glance
- Disk usage for main mount
- Memory breakdown (used / total)
- CPU usage percentage
- Top 5 processes by memory
- Color-coded metrics (green / yellow / red based on usage)
- Live monitoring mode (updates every 2s)

**Port Management (Phase 2)**

- List all listening ports with process name, PID, and memory usage
- Kill processes on specific ports with confirmation
- Kill multiple ports at once (space-separated)
- Force kill option (--force flag for non-interactive mode)
- Color-coded port types (system ports, common dev ports, ephemeral)

**Disk Analysis (Phase 3)**

- Directory Scan with file type breakdown
- Visual storage usage for top consumers
- Disk Partition Scan with smart categorization (Primary vs System)
- Cross-platform support (Mac/Linux)

**Git Shortcuts (Phase 4)**

- Quick sync - Reset to origin/main in one command
- Branch cleanup - Delete all branches except current
- Fast commits - Add + commit (or + push) in one command
- Undo/amend - Soft reset and amend helpers
- WIP commits - Quick work-in-progress saves
- Branch creation and renaming shortcuts
- Tab completion for all commands

## Quick Start

### Installation

#### Option 1: One-liner (fastest install)

```bash
curl -fsSL https://raw.githubusercontent.com/iyushkarki/csys/main/install.sh | bash
```

This installs both `csys` and `cs` (shorter alias) with shell completions.

#### Option 2: Go install (requires Go 1.19+)

```bash
go install github.com/iyushkarki/csys@latest
```

#### Option 3: Build from source

```bash
git clone https://github.com/iyushkarki/csys
cd csys
go build -o csys .
sudo mv csys /usr/local/bin/
```

### Shell Completions

After installing, set up tab completion for your shell:

```bash
# Bash
csys completion bash > /etc/bash_completion.d/csys

# Zsh (add to ~/.zshrc: fpath=(~/.zsh/completions $fpath))
mkdir -p ~/.zsh/completions
csys completion zsh > ~/.zsh/completions/_csys

# Fish
csys completion fish > ~/.config/fish/completions/csys.fish
```

The install script sets up completions automatically.

### Usage

**System Monitoring:**

```bash
# Snapshot view (one-time system check)
csys

# Live monitoring (updates every 2 seconds)
csys --live

# Help
csys --help
```

**Port Management:**

```bash
# List all listening ports
csys ports

# Kill process on port 3000
csys ports kill 3000

# Kill multiple ports
csys ports kill 3000 8080 5432

# Force kill without confirmation
csys ports kill 3000 --force

# Help
csys ports --help
csys ports kill --help
```

**Disk Analysis:**

```bash
# Scan current directory
csys scan

# Scan specific path
csys scan --path ~/Downloads

# Scan all disk partitions
csys scan disk
```

**Git Shortcuts:**

```bash
# All git commands are two words: cs g<command>

cs gsync             # Fetch + hard reset to origin/main
cs gsync develop     # Reset to origin/develop
cs gclean            # Delete all branches except current
cs gac "message"     # Add all + commit
cs gco main          # Switch to existing branch
cs gcb feature/auth  # Create and switch to new branch
cs gbrn new-name     # Rename current branch
cs gpush             # Push to remote
cs gpull             # Pull latest changes
cs gfp               # Force push (with lease)
cs gsoft 2           # Soft reset HEAD~2
cs gundo             # Undo last commit (keep staged)
cs gwip              # Quick WIP commit
cs gamend            # Amend last commit (keep message)
cs gamend "new msg"  # Amend last commit with new message
cs grb               # Rebase current branch onto origin/main
cs grb develop       # Rebase current branch onto origin/develop
cs grb main feat/x   # Checkout feat/x, rebase onto origin/main
cs glog              # Show last 10 commits (graph)
cs glog 20           # Show last 20 commits
cs gst               # Quick git status
```

## Tech Stack

- **Cobra** - CLI framework
- **Lipgloss** - Terminal styling
- **gopsutil** - Cross-platform system info & network connections
- **go-humanize** - Human-readable formatting
- **syscall** - Cross-platform process signaling (SIGTERM/SIGKILL)

## Roadmap

- **Phase 1** - Core system monitor (snapshot + live modes)
- **Phase 2** - Port management (list + kill + force kill)
- **Phase 3** - Disk analysis and directory scanning
- **Phase 4** - Git shortcuts (sync, clean, commit helpers)
- **Phase 5** - Cache detection (npm, docker, etc)
- **Phase 6** - Interactive cleanup wizard
- **Phase 7** - Advanced monitoring (network, temps, battery)

## Supported Platforms

- macOS (Intel & Apple Silicon)
- Linux (x86-64 & ARM64)

## Development

### Build

```bash
go build -o csys .
```

### Test

```bash
go test ./...
```

### Run

```bash
./csys
./csys --live
```

## License

MIT

## Contributing

Contributions welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

---

**Built with love**
