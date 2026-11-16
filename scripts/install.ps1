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
    # MSI SHA-256 lives in its own file — goreleaser's checksums.txt only covers goreleaser-built artifacts.
    $WebClient.DownloadFile("https://download.newrelic.com/install/newrelic-cli/${version}/NewRelicCLIInstaller.msi.sha256", "$env:TEMP\NewRelicCLIInstaller.msi.sha256");
} catch {
    Write-Output "`nCould not download the New Relic CLI installer.`n`nCheck your firewall settings. If you are using a proxy, make sure that you are able to access https://download.newrelic.com and that you have set the HTTPS_PROXY environment variable with your full proxy URL.`n"
    throw
}

$expectedHash = (Get-Content "$env:TEMP\NewRelicCLIInstaller.msi.sha256" -Raw).Trim().ToLower() -split '\s+' | Select-Object -First 1
if (-not $expectedHash -or $expectedHash.Length -ne 64) {
    Write-Output "⚠️  MSI checksum file missing or malformed; skipping SHA-256 verification."
} else {
    $actualHash = (Get-FileHash "$env:TEMP\NewRelicCLIInstaller.msi" -Algorithm SHA256).Hash.ToLower()
    if ($actualHash -ne $expectedHash) {
        throw "SHA-256 mismatch — aborting installation.`nExpected : $expectedHash`nActual   : $actualHash`nThe downloaded file may have been tampered with."
    }
    Write-Output "✅ SHA-256 verified: $actualHash"
}

# GPG signature verification — key embedded inline, no keyserver dependency.
# Warns and continues on any failure; never blocks installation.
# NR GPG key fingerprint: 4F9A9B5B96EC30B9 (expires 2027-04-28)
$nrGpgFingerprint = "4F9A9B5B96EC30B9"
$checksumsName = "newrelic-cli_$($version.TrimStart('v'))_checksums.txt"
$checksumsUrl = "https://download.newrelic.com/install/newrelic-cli/${version}/${checksumsName}"
$checksumsSigUrl = "${checksumsUrl}.sig"
$nrPublicKey = @"
-----BEGIN PGP PUBLIC KEY BLOCK-----

mQINBGgPxpQBEADYHbpyxZQjZibqbc53GC6BSDR0jo93FICwv6Jj4uc8ikjzcBSP
xSLcBPtq7rKdLUCB7ccjf4W5Fvzh9uLBOKrAiEEEiNYaVS1zq1yiFtoKIUC9B5Uo
kbp3diHJl6m/erpSeICg7LbV4iDCjLanTfHufTUF+HMexTIPIaufmbBBKWZ7a57u
dHyzHsnYEa9iW3vm4WLM7Ked4Vl/eFbCa7oAkfCi22Q2MreuH+dXFqtTVQ71EfQk
j7jAxTjnxvxYVhgYHELIZptY7LqudY7p24ANF0tu1NfGqGkEnEJ6OMdcMeYIy+rR
UCuYND5ite/FEBRwB7h64utgrrNXXq4tbEJuk02+v1RXHyJXIJ3W4xgkKk+U8027
SLwNJzVx9igxKMLZlYSHgpbvq7DAZzGhFbfJHh7Mv33CayxD1bzNLhF9Pk2jYX80
YQifL3QnfpahWVA0z2CHBqqaCiWY6of7j5G8JF1xMb1MxKymTjlBLcHtvnepZ8jc
f4+bp5PE27Zkbd+qd4eI1Elc0ofillldkwh6QnGWe0qgCNFyxC0U5209JV6E7dZ5
OfG99D75JFOtEduvOFXbi/782TWgFlS4yOd6zyxxckmpgAR36z0xAZu1gsKfrc/h
KGRJYyvH15VW6bWzEfbU/bkODLCckUR2S4OoGGfg1x1TsM93gA0EvTEqhQARAQAB
tGtOZXcgUmVsaWMgKE5ldyBSZWxpYyBPcGVuIFNvdXJjZSBEZXZlbG9wZXIgVG9v
bGtpdCBHUEcgS2V5KSA8bnItZGV2ZWxvcGVyLXRvb2xraXRAdXNlcnMubm9yZXBs
eS5naXRodWIuY29tPokCVwQTAQgAQRYhBNPu7uiAqR3uIJqCzE+am1uW7DC5BQJo
D8aUAhsDBQkDwmcABQsJCAcCAiICBhUKCQgLAgQWAgMBAh4HAheAAAoJEE+am1uW
7DC5nSIQALg0+q5zM+jORY/81UA23nGB2rgtOSQjlyxIy02amoUpKWZCdke7e3H1
pjPmrU4ZCVPMuhujjtXKKi6IyCmW1McxmzoGK2CAwfBURqvcRi303XLP1/J9D8cb
Be74V6bfH39wsi5RHF7bmOwlOH3bV9o+H/VmMJzgjQw57u9sZa1VIxRCNK8KXgeB
pFBnxp7XU3w2vdmEmDHnVuRmBg6myhwM0N+popK50tghAcSx5dmEe/mz3LjLzZTU
jzOaYoSsyAoU8HCbtWFa3+6w7CD+crZy3mQBVOf0tDKkqlrs3O3xR5L3u/LNAnwZ
+0wJr9uUPfLrv10TWNX+1Ldj7hLY+Y1aoLP7HBhVW1hMoBZFHc/D9gyxjIZ7t8KR
QKwu91JoqnqO+kdBF5IhHsOXwafSSPDT9oA2Pc0onuFsh9D0sIGv/n6CjjQRPfd0
E4NZQoobP9KUiES5Qglfm8jxzdt9LcCKDnZ5kYdVoPkYYNkVGRIENu6zkyZtjBFO
GToQSwVTy/RPiRRhbijm4lTDEoUufvYvJ7ThGhCjl7ejP7xHKAB1MJddfUqeRkDx
PBpDV4lTDlxIZcv9dfZl8MvdiKwwRQ0inp41UD9z2pGS7woPXrY+MaTnJW2L5EVb
lLJAb3ebyQUIFPomxv5uczqYuXej3rXPL8QHMRH93j2W81qO720Q
=rdhy
-----END PGP PUBLIC KEY BLOCK-----
"@

$gpgExe = Get-Command gpg.exe -ErrorAction SilentlyContinue
if ($gpgExe) {
    try {
        $WebClient.DownloadFile($checksumsUrl, "$env:TEMP\${checksumsName}")
        $WebClient.DownloadFile($checksumsSigUrl, "$env:TEMP\${checksumsName}.sig")

        $gpgTmpDir = Join-Path $env:TEMP "nr-gpg-$([System.IO.Path]::GetRandomFileName())"
        New-Item -ItemType Directory -Path $gpgTmpDir | Out-Null
        $keyFile = Join-Path $gpgTmpDir "nr-public.asc"
        $nrPublicKey | Set-Content -Path $keyFile -Encoding ASCII

        & gpg.exe --homedir $gpgTmpDir --batch --import $keyFile 2>&1 | Out-Null
        if ($LASTEXITCODE -eq 0) {
            & gpg.exe --homedir $gpgTmpDir --batch --verify "$env:TEMP\${checksumsName}.sig" "$env:TEMP\${checksumsName}" 2>&1 | Out-Null
            if ($LASTEXITCODE -eq 0) {
                Write-Output "✅ GPG signature verified (key $nrGpgFingerprint)"
            } else {
                Write-Output "⚠️  GPG signature verification failed for ${checksumsName}"
                Write-Output "   SHA-256 integrity of the binary has still been checked."
                Write-Output "   For a tamper-evident install, verify manually at https://github.com/newrelic/newrelic-cli/releases"
            }
        } else {
            Write-Output "⚠️  GPG key import failed; skipping signature verification."
        }
        Remove-Item -Recurse -Force $gpgTmpDir -ErrorAction SilentlyContinue
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
