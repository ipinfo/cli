$VSN = "1.0.0"

# build the filename for the Zip archive and exe file
$FileName ="cidr2ip_$($VSN)_windows_amd64"
$ZipFileName = "$($FileName).zip"

# download and extract zip
Invoke-WebRequest -Uri "https://github.com/cidr2ip/cli/releases/download/cidr2ip-$VSN/$FileName.zip" -OutFile ./$ZipFileName
Unblock-File ./$ZipFileName
Expand-Archive -Path ./$ZipFileName  -DestinationPath $env:LOCALAPPDATA\cidr2ip -Force

# delete if already exists
if (Test-Path "$env:LOCALAPPDATA\cidr2ip\cidr2ip.exe") {
    Remove-Item "$env:LOCALAPPDATA\cidr2ip\cidr2ip.exe"
}
Rename-Item -Path "$env:LOCALAPPDATA\cidr2ip\$FileName.exe" -NewName "cidr2ip.exe"

# setting up env. and cleaning files 
[System.Environment]::SetEnvironmentVariable("PATH", $Env:Path + ";$env:LOCALAPPDATA\cidr2ip", "Machine")
$env:PATH="$env:PATH;$env:LOCALAPPDATA\cidr2ip"
Remove-Item -Path ./$ZipFileName
"You can use cidr2ip now"