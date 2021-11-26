$VSN = "1.1.0"

# build the filename for the Zip archive and exe file
$FileName ="randip_$($VSN)_windows_amd64"
$ZipFileName = "$($FileName).zip"

# download and extract zip
Invoke-WebRequest -Uri "https://github.com/randip/cli/releases/download/randip-$VSN/$FileName.zip" -OutFile ./$ZipFileName
Unblock-File ./$ZipFileName
Expand-Archive -Path ./$ZipFileName  -DestinationPath $env:LOCALAPPDATA\randip -Force

# delete if already exists
if (Test-Path "$env:LOCALAPPDATA\randip\randip.exe") {
    Remove-Item "$env:LOCALAPPDATA\randip\randip.exe"
}
Rename-Item -Path "$env:LOCALAPPDATA\randip\$FileName.exe" -NewName "randip.exe"

# setting up env. and cleaning files 
[System.Environment]::SetEnvironmentVariable("PATH", $Env:Path + ";$env:LOCALAPPDATA\randip", "Machine")
$env:PATH="$env:PATH;$env:LOCALAPPDATA\randip"
Remove-Item -Path ./$ZipFileName
"You can use randip now"