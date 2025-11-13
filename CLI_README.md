# RESTerX CLI - Complete Terminal Guide

The RESTerX CLI is a powerful, lightweight command-line tool for testing REST APIs directly from your terminal. This guide covers everything you need to know to use RESTerX CLI effectively.

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

### Interactive Mode (Primary Interface)

The main way to use RESTerX CLI is through its interactive terminal interface. Simply run:

```bash
resterx-cli
```

This will launch an interactive menu where you can:

1. **Select an HTTP method** from the menu:
   - GET - Retrieve data from an API
   - POST - Send data to create resources
   - PUT - Update existing resources
   - PATCH - Partially update resources
   - DELETE - Remove resources
   - HEAD - Get headers without body
   - Exit - Close the application

2. **Enter the URL** you want to test

3. **View the formatted response** with syntax-highlighted JSON

#### Example Session

```bash
$ resterx-cli

? Select an HTTP method: 
  GET
  POST
  PUT
  PATCH
  HEAD
  DELETE
  Exit

# After selecting GET:
Selected GET method
? Enter URL:: https://api.github.com/users/octocat

# The response will be displayed with formatted JSON:
{
  "login": "octocat",
  "id": 1,
  "name": "The Octocat",
  "company": "@github",
  ...
}

# Then you're back to the method selection menu
? Select an HTTP method: 
```

### Making Different Types of Requests

#### GET Requests
GET requests are used to retrieve data from an API. Simply select GET and enter the URL:

```bash
? Select an HTTP method: GET
Selected GET method
? Enter URL:: https://jsonplaceholder.typicode.com/posts/1
```

#### POST Requests
POST requests send data to create new resources. After selecting POST and entering the URL, you'll be prompted for the request body:

```bash
? Select an HTTP method: POST
Selected POST method
? Enter URL:: https://jsonplaceholder.typicode.com/posts
? Enter request body (JSON): {"title": "Test", "body": "Test post", "userId": 1}
```

#### PUT Requests
PUT requests update existing resources completely:

```bash
? Select an HTTP method: PUT
Selected PUT method
? Enter URL:: https://jsonplaceholder.typicode.com/posts/1
? Enter request body (JSON): {"id": 1, "title": "Updated", "body": "Updated content", "userId": 1}
```

#### PATCH Requests
PATCH requests partially update resources:

```bash
? Select an HTTP method: PATCH
Selected PATCH method
? Enter URL:: https://jsonplaceholder.typicode.com/posts/1
? Enter request body (JSON): {"title": "Partially Updated"}
```

#### DELETE Requests
DELETE requests remove resources:

```bash
? Select an HTTP method: DELETE
Selected DELETE method
? Enter URL:: https://jsonplaceholder.typicode.com/posts/1
```

#### HEAD Requests
HEAD requests retrieve only headers without the response body:

```bash
? Select an HTTP method: HEAD
Selected HEAD method
? Enter URL:: https://api.github.com
```

### Version Information

Check the version and build information:

```bash
resterx-cli version
```

Output:
```
RESTerX CLI
Version:    v1.0.1
Commit:     abc1234
Build Date: 2025-11-13T12:00:00Z
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

Or just:

```bash
resterx-cli help
```

## 🔧 Available Commands

| Command | Description | Example |
|---------|-------------|---------|
| `resterx-cli` | Start interactive menu for API testing | `resterx-cli` |
| `resterx-cli version` | Show version and build information | `resterx-cli version` |
| `resterx-cli help` | Show help information | `resterx-cli --help` |

## 🌐 Features

- ✅ **Interactive Terminal UI**: User-friendly menu with keyboard navigation
- ✅ **Multiple HTTP Methods**: GET, POST, PUT, PATCH, DELETE, HEAD
- ✅ **Pretty-Printed Responses**: JSON responses with syntax highlighting
- ✅ **Response Headers**: View HTTP headers in the response
- ✅ **Status Codes**: Clear display of HTTP status codes
- ✅ **Cross-Platform**: Works on Windows, macOS, and Linux
- ✅ **Standalone Binary**: No dependencies or runtime required
- ✅ **Lightweight**: Fast startup and minimal resource usage
- ✅ **Continuous Testing**: Stay in the loop - make multiple requests without restarting

## 💡 Tips & Best Practices

### Testing Public APIs

RESTerX CLI works great with public APIs. Here are some popular ones to try:

```bash
# GitHub API
https://api.github.com/users/octocat

# JSONPlaceholder (Fake REST API for testing)
https://jsonplaceholder.typicode.com/posts
https://jsonplaceholder.typicode.com/users/1

# HTTP status codes testing
https://httpstat.us/200
https://httpstat.us/404
```

### Working with JSON

When making POST, PUT, or PATCH requests, ensure your JSON is properly formatted:

**Good:**
```json
{"title": "Test", "body": "Content", "userId": 1}
```

**Also Good (multiline):**
```json
{
  "title": "Test",
  "body": "Content",
  "userId": 1
}
```

### Keyboard Navigation

- Use **↑** and **↓** arrow keys to navigate the menu
- Press **Enter** to select an option
- Type your URL and press **Enter**
- Select "Exit" to quit the application

### Continuous Testing

The CLI stays running after each request, so you can:
1. Test multiple endpoints without restarting
2. Try different HTTP methods on the same API
3. Compare responses quickly

## 🛠️ Building from Source

If you prefer to build from source:

### Prerequisites
- Go 1.21 or later
- Git
- Make (optional, but recommended)

### Build Steps

1. Clone the repository:
   ```bash
   git clone https://github.com/AkshatNaruka/RESTerX_Packages.git
   cd RESTerX_Packages
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Build using Make:
   ```bash
   make build
   ```
   
   Or build manually:
   ```bash
   go build -o dist/resterx-cli ./cmd
   ```
   
   Or build for all platforms:
   ```bash
   make build-all
   ```

4. The binaries will be in the `dist/` directory.

### Build Targets

| Target | Description |
|--------|-------------|
| `make build` | Build for current platform |
| `make build-all` | Build for all platforms (Windows, macOS Intel, macOS ARM, Linux) |
| `make checksums` | Generate SHA256 checksums for binaries |
| `make test` | Run tests |
| `make clean` | Remove build artifacts (dist/ and build/ directories) |
| `make install` | Install binary to local Go bin directory |
| `make deps` | Download and tidy Go dependencies |
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

## 🔍 Understanding the Output

### Response Components

When you make a request, RESTerX CLI shows:

1. **Status Line**: HTTP version and status code
   ```
   HTTP/1.1 200 OK
   ```

2. **Headers**: Response headers
   ```
   Content-Type: application/json; charset=utf-8
   Content-Length: 123
   ```

3. **Body**: Pretty-printed JSON response
   ```json
   {
     "id": 1,
     "title": "Example"
   }
   ```

### HTTP Status Codes

Common status codes you'll see:

| Code | Meaning | Description |
|------|---------|-------------|
| 200 | OK | Request succeeded |
| 201 | Created | Resource created successfully |
| 204 | No Content | Success but no content to return |
| 400 | Bad Request | Invalid request format |
| 401 | Unauthorized | Authentication required |
| 403 | Forbidden | Access denied |
| 404 | Not Found | Resource doesn't exist |
| 500 | Server Error | Server-side error |

## 🎓 Example Workflows

### Example 1: Testing a Public API

```bash
$ resterx-cli
? Select an HTTP method: GET
Selected GET method
? Enter URL:: https://api.github.com/zen

# Returns GitHub's Zen quote
```

### Example 2: Creating a Resource

```bash
$ resterx-cli
? Select an HTTP method: POST
Selected POST method
? Enter URL:: https://jsonplaceholder.typicode.com/posts
? Enter request body (JSON): {"title": "My Post", "body": "Hello World", "userId": 1}

# Returns created resource with ID
```

### Example 3: Updating a Resource

```bash
$ resterx-cli
? Select an HTTP method: PUT
Selected PUT method
? Enter URL:: https://jsonplaceholder.typicode.com/posts/1
? Enter request body (JSON): {"id": 1, "title": "Updated", "body": "New content", "userId": 1}

# Returns updated resource
```

## 📊 Comparison with Other Tools

| Feature | RESTerX CLI | curl | Postman | HTTPie |
|---------|-------------|------|---------|--------|
| Interactive UI | ✅ | ❌ | ✅ | ❌ |
| Terminal-based | ✅ | ✅ | ❌ | ✅ |
| No installation needed | ✅ (single binary) | ✅ (built-in) | ❌ | ❌ |
| Pretty JSON output | ✅ | ❌ | ✅ | ✅ |
| Easy for beginners | ✅ | ❌ | ✅ | ⚠️ |
| Cross-platform | ✅ | ✅ | ✅ | ✅ |

## ❓ Frequently Asked Questions

### Q: Do I need to install anything else?
**A:** No! RESTerX CLI is a standalone binary with no dependencies.

### Q: Can I use it with authenticated APIs?
**A:** Currently, the CLI supports basic API testing. Authentication features may be added in future versions.

### Q: Does it save request history?
**A:** The current version focuses on real-time testing. History features are planned for future releases.

### Q: Can I use custom headers?
**A:** Custom headers support is planned for future versions.

### Q: Is it open source?
**A:** Yes! The source code is available in this repository.

### Q: How do I report bugs or request features?
**A:** Please open an issue on [GitHub Issues](https://github.com/AkshatNaruka/RESTerX_Packages/issues).

## 📝 Support

For issues, feature requests, or contributions:
- **GitHub Issues**: [Report bugs or request features](https://github.com/AkshatNaruka/RESTerX_Packages/issues)
- **Discussions**: [Ask questions or share ideas](https://github.com/AkshatNaruka/RESTerX_Packages/discussions)
- **Documentation**: See main [README.md](README.md)

## 🤝 Contributing

We welcome contributions! Here's how you can help:

1. **Report Bugs**: Found a bug? [Open an issue](https://github.com/AkshatNaruka/RESTerX_Packages/issues)
2. **Suggest Features**: Have an idea? [Start a discussion](https://github.com/AkshatNaruka/RESTerX_Packages/discussions)
3. **Submit PRs**: Fork, code, and submit a pull request
4. **Improve Docs**: Help us make the documentation better

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

**Happy API Testing! 🚀**
