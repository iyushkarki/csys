# csys - System Monitoring CLI

> Beautiful, developer-friendly system monitoring tool for Mac & Linux

A lightweight CLI tool that gives you instant, beautiful insights into your system's health. No more cryptic `df` output or hunting through Activity Monitor - just clean, readable information about your disk, memory, CPU, and more.

## ✨ Features

**Current (Phase 1):**
- 📊 Beautiful system overview at a glance
- 💾 Disk usage for main mount
- 🧠 Memory breakdown (used/total)
- ⚡ CPU usage percentage
- 📈 Top 5 processes by memory
- 🎨 Color-coded metrics (green/yellow/red based on usage)
- 📡 Live monitoring mode (updates every 2s)

## 🚀 Quick Start

### Installation

#### Option 1: One-liner (requires GitHub release)
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

```bash
# Snapshot view (one-time system check)
csys

# Live monitoring (updates every 2 seconds)
csys --live

# Help
csys --help
```

## 📸 Screenshot

<img width="648" height="282" alt="image" src="https://github.com/user-attachments/assets/5ae49153-0e55-4583-9f93-ad491ce9443d" />

## 🛠️ Tech Stack

- **Cobra** - CLI framework
- **Lipgloss** - Terminal styling
- **gopsutil** - Cross-platform system info
- **go-humanize** - Human-readable formatting

## 📋 Roadmap

- **Phase 1** ✅ Core system monitor (snapshot + live modes)
- **Phase 2** 🔜 Disk analysis and directory scanning
- **Phase 3** 🔜 Cache detection (npm, docker, etc)
- **Phase 4** 🔜 Interactive cleanup wizard
- **Phase 5** 🔜 Developer tools (ports, git repos)
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
