$ErrorActionPreference = 'Stop';
 
$packageName = 'ipinfo'
$toolsDir = Split-Path -Parent $MyInvocation.MyCommand.Definition

$packageArgs = @{
    packageName    = $packageName
    url64          = 'https://github.com/ipinfo/cli/releases/download/ipinfo-3.3.0/ipinfo_3.3.0_windows_amd64.zip'
    checksum64     = '2ab00f6a289d308b9ac04be9153aa566e93acff3da8a62dd283022a0cb6ee34d'
    checksumType64 = 'sha256'
    url            = 'https://github.com/ipinfo/cli/releases/download/ipinfo-3.3.0/ipinfo_3.3.0_windows_386.zip'
    checksumType   = 'sha256'
    checksum       = 'ddc0d4f17ca972cfc4bc76fc4775f67514e67bc654db4c9695993600a905525b'
    destination    = $toolsDir
}

$targetPath = Join-Path -Path $toolsDir -ChildPath "ipinfo.exe"
if (Test-Path $targetPath) {
    Remove-Item $targetPath
}
Install-ChocolateyZipPackage @packageArgs
Get-ChildItem -Path $toolsDir -Filter "ipinfo*.exe" | Rename-Item  -NewName "ipinfo.exe"
