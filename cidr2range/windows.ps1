$VSN = "1.2.0"

# build the filename for the Zip archive and exe file
$FileName ="cidr2range_$($VSN)_windows_amd64"
$ZipFileName = "$($FileName).zip"

# download and extract zip
Invoke-WebRequest -Uri "https://github.com/cidr2range/cli/releases/download/cidr2range-$VSN/$FileName.zip" -OutFile ./$ZipFileName
Unblock-File ./$ZipFileName
Expand-Archive -Path ./$ZipFileName  -DestinationPath $env:LOCALAPPDATA\cidr2range -Force

# delete if already exists
if (Test-Path "$env:LOCALAPPDATA\cidr2range\cidr2range.exe") {
    Remove-Item "$env:LOCALAPPDATA\cidr2range\cidr2range.exe"
}
Rename-Item -Path "$env:LOCALAPPDATA\cidr2range\$FileName.exe" -NewName "cidr2range.exe"

# setting up env. and cleaning files 
[System.Environment]::SetEnvironmentVariable("PATH", $Env:Path + ";$env:LOCALAPPDATA\cidr2range", "Machine")
$env:PATH="$env:PATH;$env:LOCALAPPDATA\cidr2range"
Remove-Item -Path ./$ZipFileName
"You can use cidr2range now"