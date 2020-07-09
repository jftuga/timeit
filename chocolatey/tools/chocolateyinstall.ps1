
$ErrorActionPreference = 'Stop';
$toolsDir   = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$url64      = 'https://github.com/jftuga/timeit/releases/download/v1.2.3/timeit_1.2.3_windows_amd64.zip'

$packageArgs = @{
  packageName   = $env:ChocolateyPackageName
  unzipLocation = $toolsDir
  fileType      = 'exe'
  url           = $url
  url64bit      = $url64

  softwareName  = 'timeit*'
  checksum64    = 'e1b0b0b608d0a87a8cae58e6bd4de6d43be47d9377cdab45148693db3d015206'
  checksumType64= 'sha256'

  validExitCodes= @(0)
}

Install-ChocolateyZipPackage @packageArgs
