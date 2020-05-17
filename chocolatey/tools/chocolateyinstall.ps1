
$ErrorActionPreference = 'Stop';
$toolsDir   = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$url64      = 'https://github.com/jftuga/timeit/releases/download/v1.2.1/timeit_1.2.1_windows_amd64.zip'

$packageArgs = @{
  packageName   = $env:ChocolateyPackageName
  unzipLocation = $toolsDir
  fileType      = 'exe'
  url           = $url
  url64bit      = $url64

  softwareName  = 'timeit*'
  checksum64    = 'bc8d55cdd5cdd8a83e0411792b74891758f390986eeaf442e7e50720d8f38009'
  checksumType64= 'sha256'

  validExitCodes= @(0)
}

Install-ChocolateyZipPackage @packageArgs
