$VSN = "2.6.1"

# build the filename for the Zip archive and exe file
$FileName ="ipinfo_$($VSN)_windows_amd64"
$ZipFileName = "$($FileName).zip"

# download and extract zip
Invoke-WebRequest -Uri "https://github.com/ipinfo/cli/releases/download/ipinfo-$VSN/$FileName.zip" -OutFile ./$ZipFileName
Unblock-File ./$ZipFileName
Expand-Archive -Path ./$ZipFileName  -DestinationPath $env:LOCALAPPDATA\ipinfo -Force

# delete if already exists
if (Test-Path "$env:LOCALAPPDATA\ipinfo\ipinfo.exe") {
    Remove-Item "$env:LOCALAPPDATA\ipinfo\ipinfo.exe"
}
Rename-Item -Path "$env:LOCALAPPDATA\ipinfo\$FileName.exe" -NewName "ipinfo.exe"

# setting up env. and cleaning files 
[System.Environment]::SetEnvironmentVariable("PATH", $Env:Path + ";$env:LOCALAPPDATA\ipinfo", "Machine")
$env:PATH="$env:PATH;$env:LOCALAPPDATA\ipinfo"
Remove-Item -Path ./$ZipFileName
"You can use ipinfo now"