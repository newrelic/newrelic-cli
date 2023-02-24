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
}
catch {
    Write-Output "`nCould not download the New Relic CLI installer.`n`nCheck your firewall settings. If you are using a proxy, make sure that you are able to access https://download.newrelic.com and that you have set the HTTPS_PROXY environment variable with your full proxy URL.`n"
    throw
}
msiexec.exe /qn /i $env:TEMP\NewRelicCLIInstaller.msi | Out-Null;