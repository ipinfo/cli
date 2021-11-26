$VSN = "1.0.0"

# build the filename for the Zip archive and exe file
$FileName ="prips_$($VSN)_windows_amd64"
$ZipFileName = "$($FileName).zip"

# download and extract zip
Invoke-WebRequest -Uri "https://github.com/prips/cli/releases/download/prips-$VSN/$FileName.zip" -OutFile ./$ZipFileName
Unblock-File ./$ZipFileName
Expand-Archive -Path ./$ZipFileName  -DestinationPath $env:LOCALAPPDATA\prips -Force

# delete if already exists
if (Test-Path "$env:LOCALAPPDATA\prips\prips.exe") {
    Remove-Item "$env:LOCALAPPDATA\prips\prips.exe"
}
Rename-Item -Path "$env:LOCALAPPDATA\prips\$FileName.exe" -NewName "prips.exe"

# setting up env. and cleaning files 
[System.Environment]::SetEnvironmentVariable("PATH", $Env:Path + ";$env:LOCALAPPDATA\prips", "Machine")
$env:PATH="$env:PATH;$env:LOCALAPPDATA\prips"
Remove-Item -Path ./$ZipFileName
"You can use prips now"