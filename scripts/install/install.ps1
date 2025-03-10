# Configuration Variables
$GitHubRepo = "weekend-dev-labs/lancer"
$ReleaseTag = "v4.0.0"
$InstallDir = "$env:USERPROFILE\AppData\Local\Lancer"
$ReleaseVersion = "4.0.0"

# Detect System Architecture
$Arch = if ($env:PROCESSOR_ARCHITECTURE -eq "AMD64") { "amd64" } elseif ($env:PROCESSOR_ARCHITECTURE -eq "x86") { "386" } elseif ($env:PROCESSOR_ARCHITECTURE -like "*ARM64*") { "arm64" } else { Write-Error "Unsupported architecture"; exit 1 }
$OS = "windows"

$AssetName = "lancer_${ReleaseVersion}_${OS}_${Arch}.tar.gz"
$ChecksumFile = "lancer_${ReleaseVersion}_checksums.txt"

# Log Function
function Log($Message) {
    Write-Host "[ $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss') ] $Message"
}

# Step 1: Download Release Asset and Checksum File
$DownloadUrl = "https://github.com/$GitHubRepo/releases/download/$ReleaseTag/$AssetName"
$ChecksumUrl = "https://github.com/$GitHubRepo/releases/download/$ReleaseTag/$ChecksumFile"

Log "Download URL: $DownloadUrl"
Log "Checksum URL: $ChecksumUrl"

$AssetPath = Join-Path $env:TEMP $AssetName
$ChecksumPath = Join-Path $env:TEMP $ChecksumFile

try {
    Log "Fetching release asset and checksum file..."
    Invoke-WebRequest -Uri $DownloadUrl -OutFile $AssetPath -ErrorAction Stop
    Invoke-WebRequest -Uri $ChecksumUrl -OutFile $ChecksumPath -ErrorAction Stop
} catch {
    Write-Error "Failed to download release files: $_"
    exit 1
}

# Step 2: Verify Checksum
Log "Verifying checksum..."
$ExpectedChecksum = Select-String -Path $ChecksumPath -Pattern "$AssetName" | ForEach-Object { $_.Line -replace "$AssetName", "" -replace "\s", "" }
$ActualChecksum = Get-FileHash -Path $AssetPath -Algorithm SHA256 | Select-Object -ExpandProperty Hash

if ($null -eq $ExpectedChecksum -or $ExpectedChecksum -ne $ActualChecksum) {
    Write-Error "Checksum verification failed for $AssetName (Expected: $ExpectedChecksum, Actual: $ActualChecksum)"
    exit 1
} else {
    Log "Checksum verification passed."
}

# Step 3: Extract and Install
Log "Extracting and installing..."
try {
    # Ensure destination directory exists
    if (!(Test-Path $InstallDir)) {
        New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    }

    # Use tar to extract the .tar.gz file
    tar -xzf $AssetPath -C $InstallDir

    Log "Installed successfully to $InstallDir."
} catch {
    Write-Error "Failed to extract and install: $_"
    exit 1
}

# Add Lancer path to the user's environment PATH
[System.Environment]::SetEnvironmentVariable("PATH", $env:PATH + ";$InstallDir", [System.EnvironmentVariableTarget]::User)

Write-Host "Lancer path added to environment variables: $LancerPath"