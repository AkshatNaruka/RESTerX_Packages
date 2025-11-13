#!/bin/bash
#
# Quick Setup Script for RESTerX_Packages Repository
# This script helps you quickly push the packages to GitHub
#

set -e

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}╔══════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║   RESTerX_Packages - Quick Setup Script            ║${NC}"
echo -e "${BLUE}╚══════════════════════════════════════════════════════╝${NC}"
echo ""

# Check if we're in the right directory
if [ ! -f "Makefile" ] || [ ! -f "install.sh" ]; then
    echo -e "${RED}Error: Please run this script from the RESTerX_Packages directory${NC}"
    exit 1
fi

echo -e "${BLUE}Step 1: Initializing Git repository...${NC}"
git init

echo -e "${BLUE}Step 2: Adding all files...${NC}"
git add .

echo -e "${BLUE}Step 3: Creating initial commit...${NC}"
git commit -m "Initial commit: RESTerX CLI packages repository

- CLI source code (cmd, pkg, web)
- Build system (Makefile, GitHub Actions)
- Installer scripts (install.sh, install.ps1)
- Documentation (README.md, CLI_README.md)
- Go dependencies (go.mod, go.sum)"

echo ""
echo -e "${YELLOW}Important: Make sure you've created the repository on GitHub first!${NC}"
echo -e "${YELLOW}Repository name: RESTerX_Packages${NC}"
echo -e "${YELLOW}Visibility: PUBLIC ✅${NC}"
echo ""

read -p "Have you created the repository on GitHub? (y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}Please create the repository first at:${NC}"
    echo "https://github.com/new"
    echo ""
    echo -e "${YELLOW}Then run this script again.${NC}"
    exit 0
fi

echo ""
echo -e "${BLUE}Step 4: Adding remote origin...${NC}"
git remote add origin https://github.com/AkshatNaruka/RESTerX_Packages.git

echo -e "${BLUE}Step 5: Renaming branch to main...${NC}"
git branch -M main

echo -e "${BLUE}Step 6: Pushing to GitHub...${NC}"
git push -u origin main

echo ""
echo -e "${GREEN}✅ Repository successfully pushed to GitHub!${NC}"
echo ""

read -p "Do you want to create the first release tag (v1.0.0)? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo ""
    echo -e "${BLUE}Creating release tag v1.0.0...${NC}"
    git tag v1.0.0
    git push origin v1.0.0
    
    echo ""
    echo -e "${GREEN}✅ Release tag created!${NC}"
    echo ""
    echo -e "${BLUE}GitHub Actions is now building binaries...${NC}"
    echo "Check the progress at:"
    echo "https://github.com/AkshatNaruka/RESTerX_Packages/actions"
    echo ""
    echo "Once complete, your release will be available at:"
    echo "https://github.com/AkshatNaruka/RESTerX_Packages/releases/tag/v1.0.0"
fi

echo ""
echo -e "${GREEN}╔══════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║              Setup Complete! 🎉                      ║${NC}"
echo -e "${GREEN}╚══════════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "${BLUE}Next Steps:${NC}"
echo ""
echo "1. Wait for GitHub Actions to complete (2-3 minutes)"
echo "2. Check releases: https://github.com/AkshatNaruka/RESTerX_Packages/releases"
echo "3. Test the installer:"
echo "   curl -fsSL https://raw.githubusercontent.com/AkshatNaruka/RESTerX_Packages/main/install.sh | bash"
echo ""
echo -e "${BLUE}Repository URL:${NC}"
echo "https://github.com/AkshatNaruka/RESTerX_Packages"
echo ""
