# RESTerX CLI - Terminal Client

The RESTerX CLI is a powerful command-line interface for testing APIs directly from your terminal.

## 🚀 Quick Install (Recommended)

### One-Command Installation

The easiest way to install RESTerX CLI is using our automatic installer scripts:

#### macOS / Linux
```bash
curl -fsSL https://raw.githubusercontent.com/AkshatNaruka/RESTerX_Packages/main/install.sh | bash
```

**What it does:**
- Automatically detects your OS and architecture
- Downloads the correct binary from GitHub releases
- Makes it executable
- Installs to `/usr/local/bin` (may require sudo)
- Adds to your PATH
- Shows usage instructions after installation

#### Windows (PowerShell)
```powershell
irm https://raw.githubusercontent.com/AkshatNaruka/RESTerX_Packages/main/install.ps1 | iex
```

**What it does:**
- Detects your Windows architecture
- Downloads the correct binary
- Installs to `%LOCALAPPDATA%\Programs\RESTerX`
- Adds to your user PATH automatically
- Shows usage instructions after installation

> **Note:** You may need to restart your terminal or PowerShell session for PATH changes to take effect.

---

## 📥 Manual Download

### Latest Release

If you prefer manual installation, download the precompiled binaries for your platform:

| Platform | Architecture | Download |
|----------|-------------|----------|
| 🪟 Windows | x64 | [resterx-cli-windows-amd64.exe](https://github.com/AkshatNaruka/RESTerX_Packages/releases/latest/download/resterx-cli-windows-amd64.exe) |
| 🍎 macOS | Intel (amd64) | [resterx-cli-darwin-amd64](https://github.com/AkshatNaruka/RESTerX_Packages/releases/latest/download/resterx-cli-darwin-amd64) |
| 🍎 macOS | Apple Silicon (arm64) | [resterx-cli-darwin-arm64](https://github.com/AkshatNaruka/RESTerX_Packages/releases/latest/download/resterx-cli-darwin-arm64) |
| 🐧 Linux | x64 | [resterx-cli-linux-amd64](https://github.com/AkshatNaruka/RESTerX_Packages/releases/latest/download/resterx-cli-linux-amd64) |

### Checksums

Verify your download integrity:
- [checksums.txt](https://github.com/AkshatNaruka/RESTerX_Packages/releases/latest/download/checksums.txt)

To verify:
```bash
# Download the checksums file
curl -LO https://github.com/AkshatNaruka/RESTerX_Packages/releases/latest/download/checksums.txt

# Verify your downloaded binary (example for Linux)
sha256sum -c checksums.txt --ignore-missing
```

## 📦 Installation

### Windows

1. Download `resterx-cli-windows-amd64.exe`
2. Rename to `resterx-cli.exe` (optional)
3. Add to your PATH or run directly:
   ```cmd
   .\resterx-cli.exe
   ```

### macOS

1. Download the appropriate binary for your Mac:
   - Intel Macs: `resterx-cli-darwin-amd64`
   - Apple Silicon (M1/M2/M3): `resterx-cli-darwin-arm64`

2. Make it executable:
   ```bash
   chmod +x resterx-cli-darwin-*
   ```

3. Move to your PATH (optional):
   ```bash
   sudo mv resterx-cli-darwin-* /usr/local/bin/resterx-cli
   ```

4. Run:
   ```bash
   ./resterx-cli-darwin-arm64
   # or if moved to PATH:
   resterx-cli
   ```

**Note**: On macOS, you may need to allow the app in System Preferences > Security & Privacy if you get a security warning.

### Linux

1. Download `resterx-cli-linux-amd64`

2. Make it executable:
   ```bash
   chmod +x resterx-cli-linux-amd64
   ```

3. Move to your PATH (optional):
   ```bash
   sudo mv resterx-cli-linux-amd64 /usr/local/bin/resterx-cli
   ```

4. Run:
   ```bash
   ./resterx-cli-linux-amd64
   # or if moved to PATH:
   resterx-cli
   ```

## 🎯 Usage

### Interactive Mode

Start the interactive menu by running the CLI without arguments:

```bash
resterx-cli
```

This will present you with a menu to select HTTP methods (GET, POST, PUT, PATCH, DELETE, HEAD).

### Web Interface Mode

Start the web server interface:

```bash
resterx-cli web
```

Or specify a custom port:

```bash
resterx-cli web --port 3000
```

Then open your browser to `http://localhost:8080` (or your specified port).

### Version Information

Check the version and build information:

```bash
resterx-cli version
```

Or use the standard version flag:

```bash
resterx-cli --version
```

## 🔧 Available Commands

| Command | Description |
|---------|-------------|
| `resterx-cli` | Start interactive menu for API testing |
| `resterx-cli web` | Start web server interface |
| `resterx-cli version` | Show version information |
| `resterx-cli help` | Show help information |

## 🌐 Features

- ✅ **Interactive Terminal UI**: User-friendly menu for selecting HTTP methods
- ✅ **Web Interface**: Full-featured web UI accessible via browser
- ✅ **Multiple HTTP Methods**: GET, POST, PUT, PATCH, DELETE, HEAD
- ✅ **Cross-Platform**: Works on Windows, macOS, and Linux
- ✅ **Standalone Binary**: No dependencies required
- ✅ **Lightweight**: Fast startup and low memory footprint

## 🛠️ Building from Source

If you prefer to build from source:

### Prerequisites
- Go 1.21 or later
- Git

### Build Steps

1. Clone the repository:
   ```bash
   git clone https://github.com/AkshatNaruka/RESTerX_Packages.git
   cd RESTerX
   ```

2. Build using Make:
   ```bash
   make build
   ```
   
   Or build for all platforms:
   ```bash
   make build-all
   ```

3. The binaries will be in the `dist/` directory.

### Build Targets

| Target | Description |
|--------|-------------|
| `make build` | Build for current platform |
| `make build-all` | Build for all platforms |
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
3. Alternatively, run: `xattr -d com.apple.quarantine resterx-cli-darwin-*`

### Permission Denied (Linux/macOS)

If you get a "permission denied" error:
```bash
chmod +x resterx-cli-*
```

### Command Not Found

If the binary is not in your PATH:
- Use the full path: `./resterx-cli-linux-amd64`
- Or add it to PATH: `export PATH=$PATH:$(pwd)`

## 📝 Support

For issues, feature requests, or contributions:
- GitHub Issues: https://github.com/AkshatNaruka/RESTerX_Packages/issues
- Documentation: See main [README.md](README.md)

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
