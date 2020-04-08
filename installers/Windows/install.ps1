#Requires -Version 5

# remote install:
#   Invoke-Expression (New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/newrelic/newrelic-cli/master/install.ps1')

# Exit on error
$old_erroractionpreference = $erroractionpreference
$erroractionpreference = 'stop'

if (($PSVersionTable.PSVersion.Major) -lt 5) {
    Write-Output "PowerShell 5 or later is required to install the New Relic CLI."
    Write-Output "Upgrade PowerShell: https://docs.microsoft.com/en-us/powershell/scripting/setup/installing-windows-powershell"
    break
}

# Make sure execution policy is set appropriately
$allowed = @('Unrestricted', 'RemoteSigned', 'ByPass')
if ((Get-ExecutionPolicy).ToString() -notin $allowed) {
    Write-Output "PowerShell requires an execution policy in [$($allowed -join ", ")] to install the New Relic CLI."
    Write-Output "For example, to set the execution policy to 'RemoteSigned' please run:"
    Write-Output "'Set-ExecutionPolicy RemoteSigned -scope CurrentUser'"
    break
}

function Ensure-Path($dir) { if(!(Test-Path $dir)) { mkdir $dir > $null }; Resolve-Path $dir }

function Get-FullPath($path) {
    $executionContext.sessionState.path.getUnresolvedProviderPathFromPSPath($path)
}

function Get-EnvironmentVar($name) {
    [environment]::getEnvironmentVariable($name,'User')
}

function Set-EnvironmentVar($name,$val) {
    [environment]::setEnvironmentVariable($name,$val,'User')
}

function Get-File($url,$to) {
    $wc = New-Object Net.Webclient
    $wc.Headers.Add('User-Agent', (Get-UserAgent))
    $wc.downloadFile($url,$to)
}

function Get-Useragent() {
    return "newrelic-cli-ps1-installer/1.0 (+https://raw.githubusercontent.com/newrelic/newrelic-cli/master/install.ps1/) PowerShell/$($PSVersionTable.PSVersion.Major).$($PSVersionTable.PSVersion.Minor) (Windows NT $([System.Environment]::OSVersion.Version.Major).$([System.Environment]::OSVersion.Version.Minor); $(if($env:PROCESSOR_ARCHITECTURE -eq 'AMD64'){'Win64; x64; '})$(if($env:PROCESSOR_ARCHITEW6432 -eq 'AMD64'){'WOW64; '})$PSEdition)"
}

function Write-Success($msg) { Write-Host $msg -f darkgreen }

$releaseuri = "https://github.com/newrelic/newrelic-cli/releases/latest/"
$req = Invoke-WebRequest -UseBasicParsing -Uri $releaseuri
try {
        $zipurl = ($req.Links | Where-Object href -Match "Windows_x86_64").href
        $zipurl = 'https://github.com' + $zipurl
    }
catch
    {
        Write-Warning "Error parsing release URL"
        Break
    }

# Set up the PATH variable
$installdir = $env:NEWRELIC_HOME, "$env:USERPROFILE\newrelic" | Where-Object { -not [String]::IsNullOrEmpty($_) } | Select-Object -First 1
$path = Get-Environmentvar 'PATH'
$dir = Get-Fullpath $installdir
Ensure-Path $dir
if($path -notmatch [regex]::escape($dir)) {
    Write-Output "Adding $dir to your path."

    # Set PATH for future sessions
    Set-EnvironmentVar 'PATH' "$dir;$path"

    # Set PATH for this session
    $env:PATH = "$dir;$env:PATH"
}

# Download the New Relic CLI zip
Write-Output 'Downloading New Relic CLI...'
$zipfile = "$dir\newrelic-cli.zip"
Get-File $zipurl $zipfile

Write-Output 'Extracting...'
Add-Type -Assembly "System.IO.Compression.FileSystem"
[IO.Compression.ZipFile]::ExtractToDirectory($zipfile, "$dir\_tmp")
Copy-Item "$dir\_tmp\*.exe" $dir -Recurse -Force
Remove-Item "$dir\_tmp", $zipfile -Recurse -Force

Write-Success 'New Relic CLI was installed successfully!'
Write-Output "Type 'newrelic help' for instructions."

$erroractionpreference = $old_erroractionpreference