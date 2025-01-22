$ServiceName = "Lancer"
$DisplayName = "Lancer"
$Description = "Upload Service"  # Note: New-Service doesn't support a 'Description' parameter
$ExecutablePath = "$env:USERPROFILE\AppData\Local\Lancer\lancer.exe"

# Create the service
New-Service -Name $ServiceName `
            -DisplayName $DisplayName `
            -BinaryPathName $ExecutablePath `
            -StartupType Automatic

# Start the service
Start-Service -Name $ServiceName
