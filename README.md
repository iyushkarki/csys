# csys - System Monitoring CLI

> Beautiful, developer-friendly system monitoring tool for Mac & Linux

A lightweight CLI tool that gives you instant, beautiful insights into your system's health. No more cryptic `df` output or hunting through Activity Monitor - just clean, readable information about your disk, memory, CPU, and more.

## 📸 Preview

<img width="501" height="288" alt="image" src="https://github.com/user-attachments/assets/2165c2a3-b31a-428f-b4c9-183c216d0918" />

## ✨ Features

**System Monitoring (Phase 1)**

- 🧭 **Beautiful system overview at a glance**
- 💽 **Disk usage for main mount**
- 🧠 **Memory breakdown (used / total)**
- ⚙️ **CPU usage percentage**
- 📊 **Top 5 processes by memory**
- 🎨 **Color-coded metrics (green / yellow / red based on usage)**
- 🔄 **Live monitoring mode (updates every 2s)**

**Port Management (Phase 2)**

- 🔌 **List all listening ports** with process name, PID, and memory usage
- 🛑 **Kill processes on specific ports** with confirmation
- ⚡ **Kill multiple ports at once** (space-separated)
- ⚙️ **Force kill option** (--force flag for non-interactive mode)
- 🎨 **Color-coded port types** (system ports, common dev ports, ephemeral)

**Disk Analysis (Phase 3)**

- 📂 **Directory Scan** with file type breakdown
- 📊 **Visual storage usage** for top consumers
- 💾 **Disk Partition Scan** with smart categorization (Primary vs System)
- 🖥️ **Cross-platform support** (Mac/Linux)

## 🚀 Quick Start

### Installation

#### Option 1: One-liner (fastest install)

```bash
curl -fsSL https://raw.githubusercontent.com/iyushkarki/csys/main/install.sh | bash
```

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

## 🛠️ Tech Stack

- **Cobra** - CLI framework
- **Lipgloss** - Terminal styling
- **gopsutil** - Cross-platform system info & network connections
- **go-humanize** - Human-readable formatting
- **syscall** - Cross-platform process signaling (SIGTERM/SIGKILL)

## 📋 Roadmap

- **Phase 1** ✅ Core system monitor (snapshot + live modes)
- **Phase 2** ✅ Port management (list + kill + force kill)
- **Phase 3** ✅ Disk analysis and directory scanning
- **Phase 4** 🔜 Cache detection (npm, docker, etc)
- **Phase 5** 🔜 Interactive cleanup wizard
- **Phase 6** 🔜 Advanced monitoring (network, temps, battery)

## 💻 Supported Platforms

- macOS (Intel & Apple Silicon)
- Linux (x86-64 & ARM64)

## 📝 Development

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

## 📄 License

MIT

## 🤝 Contributing

This is a personal project, but feedback and ideas are welcome!

---

**Built with ❤️**
