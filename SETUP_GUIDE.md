# RESTerX_Packages - Setup & Upload Guide

## 📦 What's in this folder?

This `RESTerX_Packages` folder contains everything you need to upload to a separate public GitHub repository for hosting CLI binaries.

## 📁 Folder Contents

```
RESTerX_Packages/
├── .github/
│   └── workflows/
│       └── build-cli.yml          # GitHub Actions workflow for automated builds
├── cmd/
│   └── main.go                    # CLI entry point
├── pkg/                           # All CLI packages (23 files)
│   ├── auth.go
│   ├── codegen.go
│   ├── collections.go
│   ├── database.go
│   ├── delete.go
│   ├── get.go
│   ├── head.go
│   ├── http_client.go
│   ├── mockserver.go
│   ├── monitoring.go
│   ├── patch.go
│   ├── payment.go
│   ├── post.go
│   ├── pretty.go
│   ├── put.go
│   ├── room.go
│   ├── storage.go
│   ├── subscription.go
│   ├── testing.go
│   ├── types.go
│   ├── variables.go
│   └── workspace.go
├── web/                           # Web server code
│   ├── api/
│   ├── static/
│   └── server.go
├── .gitignore                     # Git ignore rules
├── CLI_README.md                  # CLI documentation
├── go.mod                         # Go dependencies
├── go.sum                         # Go dependencies checksums
├── install.sh                     # Unix/Linux/macOS installer
├── install.ps1                    # Windows PowerShell installer
├── Makefile                       # Build automation
├── README.md                      # Repository README
└── SETUP_GUIDE.md                 # This file
```

## 🚀 Quick Setup Steps

### Step 1: Create the GitHub Repository

1. Go to: https://github.com/new
2. Repository name: `RESTerX_Packages`
3. Description: `Pre-built binaries and installers for RESTerX CLI`
4. **Important:** Make it **PUBLIC** ✅
5. Do NOT initialize with README (we already have one)
6. Click "Create repository"

### Step 2: Upload Files

**Option A: Using GitHub CLI (Recommended)**

```bash
cd RESTerX_Packages

# Initialize git
git init
git add .
git commit -m "Initial commit: CLI packages repository"

# Add remote (replace YOUR_USERNAME with your GitHub username)
git remote add origin https://github.com/AkshatNaruka/RESTerX_Packages.git

# Push to GitHub
git branch -M main
git push -u origin main
```

**Option B: Using GitHub Desktop**

1. Open GitHub Desktop
2. File → Add Local Repository
3. Select the `RESTerX_Packages` folder
4. Click "Publish repository"
5. Make sure "Keep this code private" is **UNCHECKED** ✅
6. Click "Publish Repository"

**Option C: Manual Upload via GitHub Web**

1. Go to your new repository: https://github.com/AkshatNaruka/RESTerX_Packages
2. Click "uploading an existing file"
3. Drag and drop the entire `RESTerX_Packages` folder contents
4. Commit changes

### Step 3: Create Your First Release

After uploading the files:

```bash
cd RESTerX_Packages

# Create and push a version tag
git tag v1.0.0
git push origin v1.0.0
```

This will trigger the GitHub Actions workflow which will:
- ✅ Build binaries for all platforms (Windows, macOS, Linux)
- ✅ Generate SHA256 checksums
- ✅ Create a GitHub release automatically
- ✅ Upload all binaries as release assets

### Step 4: Verify the Release

1. Go to: https://github.com/AkshatNaruka/RESTerX_Packages/releases
2. You should see a new release `v1.0.0`
3. Check that all binaries are attached:
   - ✅ `resterx-cli-windows-amd64.exe`
   - ✅ `resterx-cli-darwin-amd64`
   - ✅ `resterx-cli-darwin-arm64`
   - ✅ `resterx-cli-linux-amd64`
   - ✅ `checksums.txt`

### Step 5: Test the Installers

**macOS/Linux:**
```bash
curl -fsSL https://raw.githubusercontent.com/AkshatNaruka/RESTerX_Packages/main/install.sh | bash
```

**Windows:**
```powershell
irm https://raw.githubusercontent.com/AkshatNaruka/RESTerX_Packages/main/install.ps1 | iex
```

## 🔄 How to Update the CLI

When you make changes to the CLI in the main repository:

1. **Copy updated files** from main repo to `RESTerX_Packages`
2. **Commit changes:**
   ```bash
   cd RESTerX_Packages
   git add .
   git commit -m "Update CLI to version X.Y.Z"
   git push
   ```
3. **Create a new tag:**
   ```bash
   git tag v1.0.1
   git push origin v1.0.1
   ```
4. GitHub Actions will automatically build and release

## 📊 Repository Settings

After creating the repository, configure these settings:

### General Settings
- ✅ Wikis: Disabled
- ✅ Issues: Enabled
- ✅ Projects: Disabled
- ✅ Discussions: Optional

### Actions Permissions
Go to: Settings → Actions → General
- ✅ Allow all actions and reusable workflows
- ✅ Read and write permissions (for creating releases)

### Branch Protection (Optional)
Go to: Settings → Branches
- Add rule for `main` branch
- ✅ Require pull request reviews before merging

## 🔐 Security

The repository includes:
- ✅ SHA256 checksums for all binaries
- ✅ Automated builds from source (no manual uploads)
- ✅ Proper `.gitignore` to exclude build artifacts
- ✅ Version tags for tracking releases

## 📝 Important Notes

1. **Keep it Public:** The repository MUST be public for installer scripts to work
2. **Version Tags:** Always use semantic versioning (v1.0.0, v1.1.0, v2.0.0)
3. **Test Locally:** Always test builds locally before creating a release tag
4. **Main Repo:** Keep the main RESTerX repository private if needed
5. **Sync Changes:** Update this repo whenever you change CLI code

## 🆘 Troubleshooting

### GitHub Actions Fails
- Check the workflow logs: Actions tab → Latest workflow run
- Ensure Go 1.21+ is available (it should be by default)
- Verify `go.mod` is correct

### Installer Script Fails
- Verify repository is public
- Check that release assets exist
- Test download URLs manually

### Build Errors
- Run `make build` locally first to test
- Check for compilation errors in Go code
- Ensure all dependencies are in `go.mod`

## 🎉 You're Done!

Your CLI packages repository is now set up and ready to:
- ✅ Automatically build binaries on release tags
- ✅ Host downloadable CLI binaries
- ✅ Provide one-command installers
- ✅ Keep track of version history

## 📚 Additional Resources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Go Build Documentation](https://golang.org/cmd/go/#hdr-Compile_packages_and_dependencies)
- [Semantic Versioning](https://semver.org/)

## 🔗 Related Links

- Main Repository: https://github.com/AkshatNaruka/RESTerX (private)
- Packages Repository: https://github.com/AkshatNaruka/RESTerX_Packages (public)
- Latest Release: https://github.com/AkshatNaruka/RESTerX_Packages/releases/latest

---

**Created:** November 13, 2025  
**For:** RESTerX CLI Distribution
