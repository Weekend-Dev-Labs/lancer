$ServiceName = "Lancer"
$DisplayName = "Lancer"
$Description = "Upload Service"
$ExecutablePath = "C:\Program Files\Lancer\lancer.exe"

# Create the service
New-Service -Name $ServiceName 
            -DisplayName $DisplayName 
            -Description $Description 
            -BinaryPathName $ExecutablePath 
            -StartupType Automatic

# Start the service
Start-Service -Name $ServiceName