#!/bin/bash

# Configurable variables
GITHUB_REPO="weekend-dev-labs/lancer"  # Replace with the GitHub repo (e.g., user/project)
RELEASE_TAG="v1.0.0"            # Replace with the release tag (e.g., v2.0.4)
INSTALL_DIR="/usr/local/bin"    # Installation directory
RELEASE_VERSION="1.0.0"

# Detect system architecture and OS
OS="$(uname | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  i386|i686) ARCH="386" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

case "$OS" in
  darwin|linux|windows) : ;;
  *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

ASSET_NAME="lancer_${RELEASE_VERSION}_${OS}_${ARCH}.tar.gz"
CHECKSUM_FILE="lancer_${RELEASE_VERSION}_checksums.txt"

echo "$ASSET_NAME"
echo "$CHECKSUM_FILE"

# Functions
log() {
  echo "[$(date +'%Y-%m-%d %H:%M:%S')] $1"
}

error_exit() {
  echo "Error: $1" >&2
  exit 1
}

# Check for PAT argument
if [ "$#" -ne 1 ]; then
  error_exit "Usage: $0 <personal-access-token>"
fi
PAT_TOKEN="$1"

# Step 1: Download release asset and checksum file
log "Fetching release asset and checksum file..."
DOWNLOAD_URL="https://github.com/$GITHUB_REPO/releases/download/$RELEASE_TAG/$ASSET_NAME"
CHECKSUM_URL="https://github.com/$GITHUB_REPO/releases/download/$RELEASE_TAG/$CHECKSUM_FILE"

curl -H "Authorization: token $PAT_TOKEN" -L -o "$ASSET_NAME" "$DOWNLOAD_URL" || error_exit "Failed to download $ASSET_NAME"
curl -H "Authorization: token $PAT_TOKEN" -L -o "$CHECKSUM_FILE" "$CHECKSUM_URL" || error_exit "Failed to download $CHECKSUM_FILE"

# Step 2: Verify checksum
log "Verifying checksum..."
log "Asset name being searched: $ASSET_NAME in $CHECKSUM_FILE file"

EXPECTED_CHECKSUM=$(grep "$ASSET_NAME" "$CHECKSUM_FILE" | awk '{print $1}')
ACTUAL_CHECKSUM=$(sha256sum "$ASSET_NAME" | awk '{print $1}')

# Debugging outputs
log "Expected checksum (from file): $EXPECTED_CHECKSUM"
log "Actual checksum (calculated): $ACTUAL_CHECKSUM"

if [ -z "$EXPECTED_CHECKSUM" ]; then
  error_exit "No checksum found for $ASSET_NAME in $CHECKSUM_FILE. Please verify the checksum file and asset name."
fi

if [ "$EXPECTED_CHECKSUM" != "$ACTUAL_CHECKSUM" ]; then
  error_exit "Checksum verification failed for $ASSET_NAME (Expected: $EXPECTED_CHECKSUM, Actual: $ACTUAL_CHECKSUM)"
else
  log "Checksum verification passed."
fi


# Step 3: Extract and install
log "Extracting and installing..."
if [[ "$ASSET_NAME" == *.tar.gz ]]; then
  tar -xzf "$ASSET_NAME" || error_exit "Failed to extract $ASSET_NAME"
  
  # List the files in the archive and identify the binary file
  INSTALLABLE_FILES=$(tar -tf "$ASSET_NAME")

  # Log the extracted files
  log "Extracted files: $INSTALLABLE_FILES"

  # Loop through the files and look for the binary file (assuming binary has a known name or extension)
  for FILE in $INSTALLABLE_FILES; do
    # Check if the file is an executable binary (adjust the condition as necessary for your binary)
    if [[ -x "$FILE" && ! -d "$FILE" ]]; then
      # Move the binary file to the installation directory
      mv "$FILE" "$INSTALL_DIR" || error_exit "Failed to move $FILE to $INSTALL_DIR"
      log "Moved $FILE to $INSTALL_DIR"
      break  # Stop once the binary is found and moved
    fi
  done

elif [[ "$ASSET_NAME" == *.deb ]]; then
  sudo dpkg -i "$ASSET_NAME" || error_exit "Failed to install $ASSET_NAME"
else
  error_exit "Unsupported file format: $ASSET_NAME"
fi

# Cleanup step: Remove all downloaded content
log "Cleaning up downloaded files..."
rm -f "$ASSET_NAME" "$CHECKSUM_FILE" || error_exit "Failed to remove downloaded files"
log "Cleanup completed successfully."

log "Installation completed successfully."


