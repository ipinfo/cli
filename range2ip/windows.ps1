$VSN = "1.0.0"

# build the filename for the Zip archive and exe file
$FileName ="range2ip_$($VSN)_windows_amd64"
$ZipFileName = "$($FileName).zip"

# download and extract zip
Invoke-WebRequest -Uri "https://github.com/range2ip/cli/releases/download/range2ip-$VSN/$FileName.zip" -OutFile ./$ZipFileName
Unblock-File ./$ZipFileName
Expand-Archive -Path ./$ZipFileName  -DestinationPath $env:LOCALAPPDATA\range2ip -Force

# delete if already exists
if (Test-Path "$env:LOCALAPPDATA\range2ip\range2ip.exe") {
    Remove-Item "$env:LOCALAPPDATA\range2ip\range2ip.exe"
}
Rename-Item -Path "$env:LOCALAPPDATA\range2ip\$FileName.exe" -NewName "range2ip.exe"

# setting up env. and cleaning files 
[System.Environment]::SetEnvironmentVariable("PATH", $Env:Path + ";$env:LOCALAPPDATA\range2ip", "Machine")
$env:PATH="$env:PATH;$env:LOCALAPPDATA\range2ip"
Remove-Item -Path ./$ZipFileName
"You can use range2ip now"