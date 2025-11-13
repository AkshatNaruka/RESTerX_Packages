# RESTerX CLI Installer for Windows
# PowerShell script to install RESTerX CLI on Windows

param(
    [string]$InstallDir = "$env:LOCALAPPDATA\Programs\RESTerX"
)

$ErrorActionPreference = "Stop"

# Configuration
$Repo = "AkshatNaruka/RESTerX_Packages"
$BinaryName = "resterx-cli"
$GitHubAPI = "https://api.github.com/repos/$Repo/releases/latest"
$GitHubDownload = "https://github.com/$Repo/releases/latest/download"

# Colors
function Write-ColorOutput($ForegroundColor, $Message) {
    $fc = $host.UI.RawUI.ForegroundColor
    $host.UI.RawUI.ForegroundColor = $ForegroundColor
    Write-Output $Message
    $host.UI.RawUI.ForegroundColor = $fc
}

function Write-Info($Message) {
    Write-ColorOutput Blue "ℹ $Message"
}

function Write-Success($Message) {
    Write-ColorOutput Green "✓ $Message"
}

function Write-Error($Message) {
    Write-ColorOutput Red "✗ $Message"
}

function Write-Warning($Message) {
    Write-ColorOutput Yellow "⚠ $Message"
}

function Write-Header {
    Write-Output ""
    Write-ColorOutput Blue "╔══════════════════════════════════════════╗"
    Write-ColorOutput Blue "║      RESTerX CLI Installer              ║"
    Write-ColorOutput Blue "╚══════════════════════════════════════════╝"
    Write-Output ""
}

# Detect architecture
function Get-Architecture {
    $arch = $env:PROCESSOR_ARCHITECTURE
    if ($arch -eq "AMD64" -or $arch -eq "x86_64") {
        return "amd64"
    }
    elseif ($arch -eq "ARM64") {
        return "arm64"
    }
    else {
        throw "Unsupported architecture: $arch"
    }
}

# Get latest version
function Get-LatestVersion {
    Write-Info "Fetching latest version..."
    
    try {
        $response = Invoke-RestMethod -Uri $GitHubAPI -UseBasicParsing
        $version = $response.tag_name
        Write-Success "Latest version: $version"
        return $version
    }
    catch {
        Write-Warning "Could not fetch version from GitHub API, using 'latest'"
        return "latest"
    }
}

# Download binary
function Download-Binary {
    param($Arch)
    
    Write-Info "Downloading RESTerX CLI for Windows/$Arch..."
    
    $binaryName = "$BinaryName-windows-$Arch.exe"
    $downloadUrl = "$GitHubDownload/$binaryName"
    $tmpFile = "$env:TEMP\$binaryName"
    
    try {
        Invoke-WebRequest -Uri $downloadUrl -OutFile $tmpFile -UseBasicParsing
        Write-Success "Downloaded successfully"
        return $tmpFile
    }
    catch {
        Write-Error "Failed to download binary from $downloadUrl"
        Write-Error $_.Exception.Message
        exit 1
    }
}

# Install binary
function Install-Binary {
    param($TmpFile, $InstallDir)
    
    Write-Info "Installing to $InstallDir..."
    
    # Create install directory if it doesn't exist
    if (-not (Test-Path $InstallDir)) {
        New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    }
    
    # Copy binary
    $destFile = "$InstallDir\$BinaryName.exe"
    Copy-Item -Path $TmpFile -Destination $destFile -Force
    
    # Clean up temp file
    Remove-Item -Path $TmpFile -Force
    
    Write-Success "Installed successfully"
    return $destFile
}

# Add to PATH
function Add-ToPath {
    param($InstallDir)
    
    $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
    
    if ($userPath -notlike "*$InstallDir*") {
        Write-Info "Adding $InstallDir to PATH..."
        $newPath = "$userPath;$InstallDir"
        [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
        $env:Path = "$env:Path;$InstallDir"
        Write-Success "Added to PATH"
        Write-Warning "You may need to restart your terminal for PATH changes to take effect"
    }
    else {
        Write-Success "$InstallDir is already in PATH"
    }
}

# Show usage instructions
function Show-Usage {
    Write-Output ""
    Write-ColorOutput Green "╔══════════════════════════════════════════╗"
    Write-ColorOutput Green "║     Installation Complete! 🎉           ║"
    Write-ColorOutput Green "╚══════════════════════════════════════════╝"
    Write-Output ""
    Write-ColorOutput Blue "Quick Start:"
    Write-Output ""
    Write-Output "  $BinaryName              - Start interactive mode"
    Write-Output "  $BinaryName web          - Start web interface (http://localhost:8080)"
    Write-Output "  $BinaryName web -p 3000  - Start web interface on custom port"
    Write-Output "  $BinaryName version      - Show version information"
    Write-Output "  $BinaryName help         - Show help"
    Write-Output ""
    Write-ColorOutput Blue "Example:"
    Write-Output "  > $BinaryName"
    Write-Output "  Select an HTTP method: GET"
    Write-Output "  Enter URL: https://api.github.com"
    Write-Output ""
    Write-ColorOutput Blue "Documentation:"
    Write-Output "  https://github.com/$Repo"
    Write-Output ""
    Write-ColorOutput Yellow "Run '$BinaryName' now to get started!"
    Write-Output ""
    Write-Warning "Note: If command not found, please restart your terminal or PowerShell session"
    Write-Output ""
}

# Verify installation
function Test-Installation {
    param($InstallDir)
    
    Write-Info "Verifying installation..."
    
    $exePath = "$InstallDir\$BinaryName.exe"
    if (Test-Path $exePath) {
        try {
            $version = & $exePath version 2>$null | Select-String "Version:" | ForEach-Object { $_.Line.Split()[1] }
            Write-Success "RESTerX CLI is ready! ($version)"
            return $true
        }
        catch {
            Write-Warning "Binary exists but could not verify version"
            return $true
        }
    }
    else {
        Write-Error "Installation verification failed"
        return $false
    }
}

# Main installation flow
function Main {
    try {
        Write-Header
        
        $arch = Get-Architecture
        Write-Info "Detected platform: Windows/$arch"
        
        $version = Get-LatestVersion
        $tmpFile = Download-Binary -Arch $arch
        $installedFile = Install-Binary -TmpFile $tmpFile -InstallDir $InstallDir
        
        Add-ToPath -InstallDir $InstallDir
        
        if (Test-Installation -InstallDir $InstallDir) {
            Show-Usage
        }
    }
    catch {
        Write-Error "Installation failed: $($_.Exception.Message)"
        exit 1
    }
}

# Run main function
Main
