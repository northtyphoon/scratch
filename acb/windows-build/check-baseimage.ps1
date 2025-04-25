$ErrorActionPreference = 'Stop'

# fetch the current version number from base image
$current=(Get-ItemProperty "HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion")
$currUBR=$current.UBR

[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls -bor 
[Net.SecurityProtocolType]::Tls11 -bor [Net.SecurityProtocolType]::Tls12

# fetch the maximum version number from MCR by filtering and sorting the JSON result
$prefix="$($current.CurrentMajorVersionNumber).$($current.CurrentMinorVersionNumber).$($current.CurrentBuildNumber)."

try
{
    $json=$(Invoke-WebRequest -UseBasicParsing https://mcr.microsoft.com/v2/windows/servercore/tags/list | ConvertFrom-Json)
}
catch [System.Net.WebException]
{
  $ex =  $PSItem.Exception
  Write-Host "Error is " $ex.Message
  
  foreach($key in $ex.Response.Headers.Keys)
  {
     Write-Host "Response Header " $key "::" $ex.Response.Headers[$key]     
  }

  Write-Host "Waiting 3 seconds and try again"
  Start-Sleep -Seconds 3

  $json=$(Invoke-WebRequest -UseBasicParsing https://mcr.microsoft.com/v2/windows/servercore/tags/list | ConvertFrom-Json)
}

$hubUBR=($json.tags | Where-Object -FilterScript { $_.StartsWith($prefix) -and $_ -Match "^\d+\.\d+\.\d+\.\d+$" } |%{[System.Version]$_}|sort)[-1].Revision

Write-Output "Base image Update Build Revision $currUBR, Hub Update Build Revision $hubUBR"