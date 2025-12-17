# Contributing to csys

Thanks for your interest in contributing! 🎉

## Prerequisites

- **Go 1.21+** ([download](https://go.dev/dl/))
- **Git**
- **macOS or Linux**

## Setup

```bash
git clone https://github.com/<your-username>/csys.git
cd csys
go mod download
go build -o csys .
./csys
```

## Development Workflow

```bash
go build -o csys .   # Build the binary
./csys               # Test your changes
go fmt ./...         # Auto-format code (fixes spacing, brackets)
go vet ./...         # Catch common bugs (unused vars, bad calls)
```

## PR Workflow

1. **Fork** & clone the repo
2. **Create a branch** (see naming below)
3. **Make changes** – use the commands above to build & verify
4. **Commit** (see message format below)
5. **Push** & open a Pull Request

### Branch Naming

| Prefix      | Use for       |
| ----------- | ------------- |
| `feature/`  | New features  |
| `fix/`      | Bug fixes     |
| `docs/`     | Documentation |
| `refactor/` | Code cleanup  |

Example: `feature/cache-detection`, `fix/memory-display`, `docs/readme-update`

### Commit Messages

Format: `type: short description`

| Type        | Use for       |
| ----------- | ------------- |
| `feat:`     | New feature   |
| `fix:`      | Bug fix       |
| `docs:`     | Documentation |
| `refactor:` | Code cleanup  |
| `chore:`    | Maintenance   |

Example: `feat: add docker cache detection`, `fix: correct memory calculation`

## What's Welcome?

- 🐛 **Bug fixes**
- ✨ **Features** – see [roadmap](README.md#-roadmap)
- 📝 **Docs** – typos, examples, clarity
- 🐧 **Platform support** – Linux edge cases

## Questions?

Open an issue – happy to help!

---

**Thanks!** 🙌
