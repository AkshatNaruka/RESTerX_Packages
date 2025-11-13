# RESTerX CLI Packages

This repository contains pre-built binaries and installer scripts for the [RESTerX](https://github.com/AkshatNaruka/RESTerX) CLI tool.

## 🚀 Quick Install

### macOS / Linux
```bash
curl -fsSL https://raw.githubusercontent.com/AkshatNaruka/RESTerX_Packages/main/install.sh | bash
```

### Windows (PowerShell)
```powershell
irm https://raw.githubusercontent.com/AkshatNaruka/RESTerX_Packages/main/install.ps1 | iex
```

## 📥 Manual Download

Download the latest pre-built binaries from the [Releases](https://github.com/AkshatNaruka/RESTerX_Packages/releases/latest) page:

| Platform | Architecture | Download Link |
|----------|-------------|---------------|
| 🪟 Windows | x64 (amd64) | [Download](https://github.com/AkshatNaruka/RESTerX_Packages/releases/latest/download/resterx-cli-windows-amd64.exe) |
| 🍎 macOS | Intel (amd64) | [Download](https://github.com/AkshatNaruka/RESTerX_Packages/releases/latest/download/resterx-cli-darwin-amd64) |
| 🍎 macOS | Apple Silicon (arm64) | [Download](https://github.com/AkshatNaruka/RESTerX_Packages/releases/latest/download/resterx-cli-darwin-arm64) |
| 🐧 Linux | x64 (amd64) | [Download](https://github.com/AkshatNaruka/RESTerX_Packages/releases/latest/download/resterx-cli-linux-amd64) |

### Verify Downloads

Download the [checksums.txt](https://github.com/AkshatNaruka/RESTerX_Packages/releases/latest/download/checksums.txt) file and verify:

```bash
sha256sum -c checksums.txt --ignore-missing
```

## 📖 Documentation

For complete documentation, visit:
- [CLI Documentation](./CLI_README.md)
- [Main Repository](https://github.com/AkshatNaruka/RESTerX)

## 🔧 What is RESTerX?

RESTerX is a powerful API testing tool with:
- ✨ Interactive CLI interface
- 🌐 Web-based UI
- 📝 Request history and collections
- 🔐 Authentication support
- 🧪 Testing and mock server capabilities
- 💾 Workspace management

## 🏗️ Building from Source

If you want to build the CLI yourself:

```bash
# Clone this repository
git clone https://github.com/AkshatNaruka/RESTerX_Packages.git
cd RESTerX_Packages

# Install dependencies
go mod download

# Build for your platform
make build

# Or build for all platforms
make build-all
```

## 📦 Release Process

This repository uses GitHub Actions to automatically build binaries when a new tag is pushed:

```bash
git tag v1.0.0
git push origin v1.0.0
```

The CI/CD pipeline will:
1. Build binaries for all supported platforms
2. Generate SHA256 checksums
3. Create a GitHub release
4. Upload all artifacts

## 🤝 Contributing

Contributions are welcome! Please submit issues and pull requests to the [main RESTerX repository](https://github.com/AkshatNaruka/RESTerX).

## 📄 License

This project is licensed under the MIT License.

## 🔗 Links

- [Main Repository](https://github.com/AkshatNaruka/RESTerX)
- [Latest Release](https://github.com/AkshatNaruka/RESTerX_Packages/releases/latest)
- [Report Issues](https://github.com/AkshatNaruka/RESTerX_Packages/issues)

---

**Note:** This is a separate repository created specifically for hosting pre-built binaries. The source code is maintained in the [main RESTerX repository](https://github.com/AkshatNaruka/RESTerX).
