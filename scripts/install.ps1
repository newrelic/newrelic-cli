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
    $checksumsName = "newrelic-cli_$($version.TrimStart('v'))_checksums.txt"
    $WebClient.DownloadFile("https://download.newrelic.com/install/newrelic-cli/${version}/NewRelicCLIInstaller.msi", "$env:TEMP\NewRelicCLIInstaller.msi");
    $WebClient.DownloadFile("https://download.newrelic.com/install/newrelic-cli/${version}/${checksumsName}", "$env:TEMP\${checksumsName}");
} catch {
    Write-Output "`nCould not download the New Relic CLI installer.`n`nCheck your firewall settings. If you are using a proxy, make sure that you are able to access https://download.newrelic.com and that you have set the HTTPS_PROXY environment variable with your full proxy URL.`n"
    throw
}

$checksumContent = Get-Content "$env:TEMP\${checksumsName}"
$expectedHash = $null
foreach ($line in $checksumContent) {
    if ($line -match "NewRelicCLIInstaller\.msi") {
        $expectedHash = ($line -split '\s+')[0].Trim()
        break
    }
}

if (-not $expectedHash) {
    throw "Could not find checksum for NewRelicCLIInstaller.msi in ${checksumsName}. The release may be incomplete."
}

$actualHash = (Get-FileHash "$env:TEMP\NewRelicCLIInstaller.msi" -Algorithm SHA256).Hash.ToLower()
if ($actualHash -ne $expectedHash.ToLower()) {
    throw "Checksum mismatch - aborting installation.`nExpected : $expectedHash`nActual   : $actualHash`nThe downloaded file may have been tampered with."
}

Write-Output "✅ SHA-256 verified: $actualHash"

# GPG signature verification of the checksums file.
# Warns and continues on any failure — never blocks installation.
# NR GPG key fingerprint: 4F9A9B5B96EC30B9 (expires 2027-04-27)
# Published on keys.openpgp.org — fetch via: gpg --keyserver hkps://keys.openpgp.org --recv-keys 4F9A9B5B96EC30B9
$nrGpgFingerprint = "4F9A9B5B96EC30B9"
$nrGpgKeyserver = "hkps://keys.openpgp.org"
$checksumsUrl = "https://download.newrelic.com/install/newrelic-cli/${version}/${checksumsName}"
$checksumsSigUrl = "${checksumsUrl}.sig"

$gpgExe = Get-Command gpg.exe -ErrorAction SilentlyContinue
if ($gpgExe) {
    try {
        $WebClient.DownloadFile($checksumsSigUrl, "$env:TEMP\${checksumsName}.sig")

        # Import the NR public key from keys.openpgp.org if not already in the keyring
        $keyCheck = & gpg.exe --list-keys $nrGpgFingerprint 2>&1
        if ($LASTEXITCODE -ne 0) {
            & gpg.exe --keyserver $nrGpgKeyserver --recv-keys $nrGpgFingerprint 2>&1 | Out-Null
        }

        $verifyOutput = & gpg.exe --verify "$env:TEMP\${checksumsName}.sig" "$env:TEMP\${checksumsName}" 2>&1
        if ($LASTEXITCODE -eq 0) {
            Write-Output "✅ GPG signature verified (key $nrGpgFingerprint)"
        } else {
            Write-Output "⚠️  GPG signature verification failed for ${checksumsName}"
            Write-Output "   The checksums file could not be confirmed as genuine."
            Write-Output "   SHA-256 integrity of the binary has still been checked."
            Write-Output "   For a tamper-evident install, verify manually at https://github.com/newrelic/newrelic-cli/releases"
        }
    } catch {
        Write-Output "⚠️  GPG signature verification skipped: $($_.Exception.Message)"
    }
} else {
    Write-Output "⚠️  gpg.exe not found; skipping signature verification."
    Write-Output "   Install Gpg4win for stronger supply-chain assurance."
}

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
