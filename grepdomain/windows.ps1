$ErrorActionPreference = "Stop"

$VSN = "3.3.3"

# build the filename for the Zip archive
$FileName = "grepdomain_$($VSN)_windows_amd64"
$ZipFileName = "$($FileName).zip"

# download and extract zip; -Force overwrites any previously installed exe
Invoke-WebRequest -Uri "https://github.com/ipinfo/cli/releases/download/grepdomain-$VSN/$FileName.zip" -OutFile ./$ZipFileName
Unblock-File ./$ZipFileName
Expand-Archive -Path ./$ZipFileName  -DestinationPath $env:LOCALAPPDATA\ipinfo -Force

# setting up env. 
$PathContent = [Environment]::GetEnvironmentVariable('path', 'Machine')
$IPinfoPath = "$env:LOCALAPPDATA\ipinfo"

# if Path already exists
if ($PathContent -ne $null) {
    if (-Not($PathContent -split ';' -contains $IPinfoPath)) {
        [System.Environment]::SetEnvironmentVariable("PATH", $Env:Path + ";$env:LOCALAPPDATA\ipinfo", "Machine")
    }
}
else {
    [System.Environment]::SetEnvironmentVariable("PATH", $Env:Path + ";$env:LOCALAPPDATA\ipinfo", "Machine")
}

# cleaning files
Remove-Item -Path ./$ZipFileName
"You can use grepdomain now."
