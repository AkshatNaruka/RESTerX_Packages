# RESTerX CLI - Terminal API Testing Tool

RESTerX CLI is a powerful, lightweight command-line tool for testing REST APIs directly from your terminal. Built with Go, it provides an interactive interface for making HTTP requests and viewing responses in a clean, formatted way.

## ✨ Features

- **Interactive Terminal Interface**: User-friendly menu for selecting HTTP methods
- **Multiple HTTP Methods**: Support for GET, POST, PUT, PATCH, DELETE, and HEAD requests
- **Response Formatting**: Pretty-printed JSON responses with syntax highlighting
- **Cross-Platform**: Works seamlessly on Windows, macOS, and Linux
- **Lightweight & Fast**: Single binary with no dependencies
- **Easy Installation**: One-command installation scripts for all platforms

## 🚀 Quick Install Guide

### macOS / Linux
```bash
curl -fsSL https://raw.githubusercontent.com/AkshatNaruka/RESTerX_Packages/main/install.sh | bash
```

**What it does:**
- Automatically detects your OS and architecture
- Downloads the correct binary from GitHub releases
- Installs to `/usr/local/bin` (may require sudo)
- Makes it executable and adds to your PATH

### Windows (PowerShell)
```powershell
irm https://raw.githubusercontent.com/AkshatNaruka/RESTerX_Packages/main/install.ps1 | iex
```

**What it does:**
- Detects your Windows architecture
- Downloads the correct binary
- Installs to `%LOCALAPPDATA%\Programs\RESTerX`
- Automatically adds to your user PATH

> **Note:** You may need to restart your terminal or PowerShell session for PATH changes to take effect.

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

### Installation After Manual Download

**Windows:**
1. Download `resterx-cli-windows-amd64.exe`
2. Rename to `resterx-cli.exe` (optional)
3. Add to your PATH or run directly

**macOS:**
1. Download the appropriate binary for your Mac
2. Make it executable: `chmod +x resterx-cli-darwin-*`
3. Move to PATH: `sudo mv resterx-cli-darwin-* /usr/local/bin/resterx-cli`
4. Run: `resterx-cli`

**Note:** On macOS, you may need to allow the app in System Preferences > Security & Privacy.

**Linux:**
1. Download `resterx-cli-linux-amd64`
2. Make it executable: `chmod +x resterx-cli-linux-amd64`
3. Move to PATH: `sudo mv resterx-cli-linux-amd64 /usr/local/bin/resterx-cli`
4. Run: `resterx-cli`

## 🎯 Usage

### Interactive Mode

Start the interactive menu by running the CLI without arguments:

```bash
resterx-cli
```

This will present you with an interactive menu where you can:
1. Select an HTTP method (GET, POST, PUT, PATCH, DELETE, HEAD)
2. Enter the URL you want to test
3. View the formatted response

**Example workflow:**
```bash
$ resterx-cli
? Select an HTTP method: GET
Selected GET method
? Enter URL:: https://api.github.com/users/octocat
# Response will be displayed with formatted JSON
```

### Version Information

Check the version and build information:

```bash
resterx-cli version
```

Or use the standard version flag:

```bash
resterx-cli --version
```

### Getting Help

View available commands and options:

```bash
resterx-cli --help
```

## 📚 Detailed Documentation

For more detailed information about installation and usage, see [CLI_README.md](./CLI_README.md).

## 🔧 What is RESTerX?

RESTerX is a powerful API testing tool designed for developers who prefer working in the terminal. Key features include:

- **Interactive CLI interface**: Easy-to-use terminal UI for API testing
- **Multiple HTTP methods**: GET, POST, PUT, PATCH, DELETE, HEAD
- **Response formatting**: Pretty-printed, colorized JSON output
- **Cross-platform support**: Works on Windows, macOS, and Linux
- **No dependencies**: Single standalone binary
- **Lightweight**: Fast startup and minimal resource usage

## 🛠️ Building from Source

If you prefer to build the CLI yourself:

### Prerequisites
- Go 1.21 or later
- Git

### Build Steps

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

The binaries will be available in the `dist/` directory.

### Available Make Targets

| Target | Description |
|--------|-------------|
| `make build` | Build for current platform |
| `make build-all` | Build for all platforms (Windows, macOS, Linux) |
| `make checksums` | Generate SHA256 checksums |
| `make test` | Run tests |
| `make clean` | Remove build artifacts |
| `make install` | Install binary locally |
| `make help` | Show all available targets |

## 🐛 Troubleshooting

### macOS Security Warning

If you see "cannot be opened because the developer cannot be verified":
1. Go to System Preferences > Security & Privacy
2. Click "Open Anyway" for the blocked application
3. Or run: `xattr -d com.apple.quarantine resterx-cli-darwin-*`

### Permission Denied (Linux/macOS)

If you get a "permission denied" error:
```bash
chmod +x resterx-cli-*
```

### Command Not Found

If the binary is not in your PATH after installation:
- Restart your terminal/shell session
- Or use the full path to run the binary
- Or manually add it to PATH: `export PATH=$PATH:/usr/local/bin`

### Windows Installation Issues

If the installer fails on Windows:
- Run PowerShell as Administrator
- Ensure your execution policy allows scripts: `Set-ExecutionPolicy RemoteSigned -Scope CurrentUser`
- Restart PowerShell after installation

## 📦 Release Process

This repository uses GitHub Actions to automatically build binaries when a new tag is pushed:

```bash
git tag v1.0.2
git push origin v1.0.2
```

The CI/CD pipeline automatically:
1. ✅ Builds binaries for all supported platforms
2. ✅ Generates SHA256 checksums
3. ✅ Creates a GitHub release
4. ✅ Uploads all artifacts

## 🤝 Contributing

Contributions are welcome! To contribute:

1. Fork this repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

For bug reports and feature requests, please [open an issue](https://github.com/AkshatNaruka/RESTerX_Packages/issues).

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🔗 Links

- [Latest Release](https://github.com/AkshatNaruka/RESTerX_Packages/releases/latest)
- [Report Issues](https://github.com/AkshatNaruka/RESTerX_Packages/issues)
- [Changelog](https://github.com/AkshatNaruka/RESTerX_Packages/releases)

## 💬 Support

- **Issues**: [GitHub Issues](https://github.com/AkshatNaruka/RESTerX_Packages/issues)
- **Discussions**: [GitHub Discussions](https://github.com/AkshatNaruka/RESTerX_Packages/discussions)

---

**Made with ❤️ for developers who love the terminal**
