#!/bin/bash
#
# RESTerX CLI Installer
# This script installs the RESTerX CLI tool on Unix-like systems (macOS, Linux)
#

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
REPO="AkshatNaruka/RESTerX_Packages"
BINARY_NAME="resterx-cli"
INSTALL_DIR="/usr/local/bin"
GITHUB_API="https://api.github.com/repos/${REPO}/releases/latest"
GITHUB_DOWNLOAD="https://github.com/${REPO}/releases/latest/download"

# Print with color
print_info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

print_header() {
    echo ""
    echo -e "${BLUE}╔══════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║${NC}      ${GREEN}RESTerX CLI Installer${NC}              ${BLUE}║${NC}"
    echo -e "${BLUE}╔══════════════════════════════════════════╗${NC}"
    echo ""
}

# Detect OS and Architecture
detect_platform() {
    local os=$(uname -s | tr '[:upper:]' '[:lower:]')
    local arch=$(uname -m)
    
    case "$os" in
        darwin)
            OS="darwin"
            ;;
        linux)
            OS="linux"
            ;;
        *)
            print_error "Unsupported operating system: $os"
            exit 1
            ;;
    esac
    
    case "$arch" in
        x86_64|amd64)
            ARCH="amd64"
            ;;
        arm64|aarch64)
            ARCH="arm64"
            ;;
        *)
            print_error "Unsupported architecture: $arch"
            exit 1
            ;;
    esac
    
    PLATFORM="${OS}-${ARCH}"
}

# Check prerequisites
check_prerequisites() {
    print_info "Checking prerequisites..."
    
    # Check if curl or wget is available
    if command -v curl &> /dev/null; then
        DOWNLOADER="curl -fsSL"
    elif command -v wget &> /dev/null; then
        DOWNLOADER="wget -qO-"
    else
        print_error "Neither curl nor wget found. Please install one of them."
        exit 1
    fi
    
    print_success "Prerequisites OK"
}

# Get latest version
get_latest_version() {
    print_info "Fetching latest version..."
    
    if command -v curl &> /dev/null; then
        VERSION=$(curl -s "${GITHUB_API}" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')
    else
        VERSION=$(wget -qO- "${GITHUB_API}" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')
    fi
    
    if [ -z "$VERSION" ]; then
        print_warning "Could not fetch version from GitHub API, using 'latest'"
        VERSION="latest"
    else
        print_success "Latest version: ${VERSION}"
    fi
}

# Download binary
download_binary() {
    print_info "Downloading RESTerX CLI for ${OS}/${ARCH}..."
    
    local binary_name="${BINARY_NAME}-${PLATFORM}"
    local download_url="${GITHUB_DOWNLOAD}/${binary_name}"
    local tmp_file="/tmp/${binary_name}"
    
    if command -v curl &> /dev/null; then
        if ! curl -fsSL -o "${tmp_file}" "${download_url}"; then
            print_error "Failed to download binary from ${download_url}"
            exit 1
        fi
    else
        if ! wget -qO "${tmp_file}" "${download_url}"; then
            print_error "Failed to download binary from ${download_url}"
            exit 1
        fi
    fi
    
    print_success "Downloaded successfully"
    echo "$tmp_file"
}

# Install binary
install_binary() {
    local tmp_file=$1
    
    print_info "Installing to ${INSTALL_DIR}/${BINARY_NAME}..."
    
    # Make binary executable
    chmod +x "${tmp_file}"
    
    # Check if we need sudo
    if [ -w "${INSTALL_DIR}" ]; then
        mv "${tmp_file}" "${INSTALL_DIR}/${BINARY_NAME}"
    else
        print_warning "Need sudo permissions to install to ${INSTALL_DIR}"
        sudo mv "${tmp_file}" "${INSTALL_DIR}/${BINARY_NAME}"
    fi
    
    print_success "Installed successfully"
}

# Show usage instructions
show_usage() {
    echo ""
    echo -e "${GREEN}╔══════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}║${NC}     Installation Complete! 🎉           ${GREEN}║${NC}"
    echo -e "${GREEN}╚══════════════════════════════════════════╝${NC}"
    echo ""
    echo -e "${BLUE}Quick Start:${NC}"
    echo ""
    echo -e "  ${GREEN}${BINARY_NAME}${NC}              - Start interactive mode"
    echo -e "  ${GREEN}${BINARY_NAME} web${NC}          - Start web interface (http://localhost:8080)"
    echo -e "  ${GREEN}${BINARY_NAME} web -p 3000${NC} - Start web interface on custom port"
    echo -e "  ${GREEN}${BINARY_NAME} version${NC}      - Show version information"
    echo -e "  ${GREEN}${BINARY_NAME} help${NC}         - Show help"
    echo ""
    echo -e "${BLUE}Example:${NC}"
    echo -e "  $ ${BINARY_NAME}"
    echo -e "  Select an HTTP method: GET"
    echo -e "  Enter URL: https://api.github.com"
    echo ""
    echo -e "${BLUE}Documentation:${NC}"
    echo -e "  https://github.com/${REPO}"
    echo ""
    echo -e "${YELLOW}Run '${BINARY_NAME}' now to get started!${NC}"
    echo ""
}

# Verify installation
verify_installation() {
    print_info "Verifying installation..."
    
    if command -v ${BINARY_NAME} &> /dev/null; then
        local version=$(${BINARY_NAME} version 2>/dev/null | grep "Version:" | awk '{print $2}')
        print_success "RESTerX CLI is ready! (${version})"
        return 0
    else
        print_error "Installation verification failed"
        print_warning "You may need to add ${INSTALL_DIR} to your PATH"
        print_info "Add this to your ~/.bashrc or ~/.zshrc:"
        echo "    export PATH=\"\$PATH:${INSTALL_DIR}\""
        return 1
    fi
}

# Main installation flow
main() {
    print_header
    
    detect_platform
    print_info "Detected platform: ${OS}/${ARCH}"
    
    check_prerequisites
    get_latest_version
    
    local tmp_file=$(download_binary)
    install_binary "${tmp_file}"
    
    if verify_installation; then
        show_usage
    fi
}

# Run main function
main
