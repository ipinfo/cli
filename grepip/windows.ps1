$VSN = "1.2.1"

# build the filename for the Zip archive and exe file
$FileName ="grepip_$($VSN)_windows_amd64"
$ZipFileName = "$($FileName).zip"

# download and extract zip
Invoke-WebRequest -Uri "https://github.com/grepip/cli/releases/download/grepip-$VSN/$FileName.zip" -OutFile ./$ZipFileName
Unblock-File ./$ZipFileName
Expand-Archive -Path ./$ZipFileName  -DestinationPath $env:LOCALAPPDATA\grepip -Force

# delete if already exists
if (Test-Path "$env:LOCALAPPDATA\grepip\grepip.exe") {
    Remove-Item "$env:LOCALAPPDATA\grepip\grepip.exe"
}
Rename-Item -Path "$env:LOCALAPPDATA\grepip\$FileName.exe" -NewName "grepip.exe"

# setting up env. and cleaning files 
[System.Environment]::SetEnvironmentVariable("PATH", $Env:Path + ";$env:LOCALAPPDATA\grepip", "Machine")
$env:PATH="$env:PATH;$env:LOCALAPPDATA\grepip"
Remove-Item -Path ./$ZipFileName
"You can use grepip now"