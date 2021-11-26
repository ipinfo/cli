$VSN = "1.2.0"

# build the filename for the Zip archive and exe file
$FileName ="range2cidr_$($VSN)_windows_amd64"
$ZipFileName = "$($FileName).zip"

# download and extract zip
Invoke-WebRequest -Uri "https://github.com/range2cidr/cli/releases/download/range2cidr-$VSN/$FileName.zip" -OutFile ./$ZipFileName
Unblock-File ./$ZipFileName
Expand-Archive -Path ./$ZipFileName  -DestinationPath $env:LOCALAPPDATA\range2cidr -Force

# delete if already exists
if (Test-Path "$env:LOCALAPPDATA\range2cidr\range2cidr.exe") {
    Remove-Item "$env:LOCALAPPDATA\range2cidr\range2cidr.exe"
}
Rename-Item -Path "$env:LOCALAPPDATA\range2cidr\$FileName.exe" -NewName "range2cidr.exe"

# setting up env. and cleaning files 
[System.Environment]::SetEnvironmentVariable("PATH", $Env:Path + ";$env:LOCALAPPDATA\range2cidr", "Machine")
$env:PATH="$env:PATH;$env:LOCALAPPDATA\range2cidr"
Remove-Item -Path ./$ZipFileName
"You can use range2cidr now"