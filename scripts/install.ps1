$p = New-Object System.Security.Principal.WindowsPrincipal([System.Security.Principal.WindowsIdentity]::GetCurrent())

if (!$p.IsInRole([System.Security.Principal.WindowsBuiltInRole]::Administrator)) {
    throw 'This script requires admin privileges to run and the current Windows PowerShell session is not running as Administrator. Start Windows PowerShell by using the Run as Administrator option, and then try running the script again.'
}

[Net.ServicePointManager]::SecurityProtocol = 'tls12, tls';

$WebClient = New-Object System.Net.WebClient

if ($env:HTTPS_PROXY) {
    $WebClient.Proxy = New-Object System.Net.WebProxy($env:HTTPS_PROXY, $true)

}

$version = $null

try {
    $version = $WebClient.DownloadString("https://download.newrelic.com/install/newrelic-cli/currentVersion.txt").Trim();
    $WebClient.DownloadFile("https://download.newrelic.com/install/newrelic-cli/${version}/NewRelicCLIInstaller.msi", "$env:TEMP\NewRelicCLIInstaller.msi");
    $WebClient.DownloadFile("https://download.newrelic.com/install/newrelic-cli/${version}/checksums_sha256.txt", "$env:TEMP\checksums_sha256.txt");
} catch {
    Write-Output "`nCould not download the New Relic CLI installer.`n`nCheck your firewall settings. If you are using a proxy, make sure that you are able to access https://download.newrelic.com and that you have set the HTTPS_PROXY environment variable with your full proxy URL.`n"
    throw
}

$checksumContent = Get-Content "$env:TEMP\checksums_sha256.txt"
$expectedHash = $null
foreach ($line in $checksumContent) {
    if ($line -match "NewRelicCLIInstaller\.msi") {
        $expectedHash = ($line -split '\s+')[0].Trim()
        break
    }
}

if (-not $expectedHash) {
    throw "Could not find checksum for NewRelicCLIInstaller.msi in checksums_sha256.txt. The release may be incomplete."
}

$actualHash = (Get-FileHash "$env:TEMP\NewRelicCLIInstaller.msi" -Algorithm SHA256).Hash.ToLower()
if ($actualHash -ne $expectedHash.ToLower()) {
    throw "Checksum mismatch - aborting installation.`nExpected : $expectedHash`nActual   : $actualHash`nThe downloaded file may have been tampered with."
}

Write-Output "Checksum verified: $actualHash"

try {
function Find-UninstallGuids {
    param (
        [Parameter(Mandatory)]
        [string]$Match
    )

    $baseKeys = Get-ChildItem -Path HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall `
    | % { $_.Name.TrimStart("HKEY_LOCAL_MACHINE\") }

    $wowKeys = Get-ChildItem -Path HKLM:\SOFTWARE\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall `
    | % { $_.Name.TrimStart("HKEY_LOCAL_MACHINE\") }

    $allKeys = $baseKeys + $wowKeys

    $uninstallIds = New-Object System.Collections.ArrayList
    foreach ($key in $allKeys) {
        $keyData = Get-Item -LiteralPath "HKLM:\$key" -ErrorAction SilentlyContinue
        if ($keyData) {
            $name = $keyData.GetValue("DisplayName")
            if ($name -and $name -match $Match) {
                $keyId = Split-Path $key -Leaf
                $uninstallIds.Add($keyId) | Out-Null
            }
        }
    }

    if ($uninstallIds.Count -eq 0) {
        return @()
    }

    return $uninstallIds
}
  $uninstallIds = Find-UninstallGuids -Match "New Relic CLI"

  foreach ($uninstallId in $uninstallIds) {
    $arguments = "/x $uninstallId /qn"

    try {
      Start-Process msiexec.exe -ArgumentList $arguments
    } catch {
      throw $_.Exception
    }
  }
} catch {
  throw $_.Exception
}

msiexec.exe /qn /i $env:TEMP\NewRelicCLIInstaller.msi | Out-Null;