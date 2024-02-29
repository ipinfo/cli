$ErrorActionPreference = 'Stop';
 
$packageName = 'ipinfo'
$toolsDir = Split-Path -Parent $MyInvocation.MyCommand.Definition

$packageArgs = @{
    packageName    = $packageName
    url64          = ''
    checksum64     = ''
    checksumType64 = ''
    url            = ''
    checksumType   = ''
    checksum       = ''
    destination    = $toolsDir
}

$targetPath = Join-Path -Path $toolsDir -ChildPath "ipinfo.exe"
if (Test-Path $targetPath) {
    Remove-Item $targetPath
}
Install-ChocolateyZipPackage @packageArgs
Get-ChildItem -Path $toolsDir -Filter "ipinfo*.exe" | Rename-Item  -NewName "ipinfo.exe"
