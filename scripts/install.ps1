$p = New-Object System.Security.Principal.WindowsPrincipal([System.Security.Principal.WindowsIdentity]::GetCurrent())
if (!$p.IsInRole([System.Security.Principal.WindowsBuiltInRole]::Administrator)) {
    throw 'This script requires admin privileges to run and the current Windows PowerShell session is not running as Administrator. Start Windows PowerShell by using the Run as Administrator option, and then try running the script again.'
}
[Net.ServicePointManager]::SecurityProtocol = 'tls12, tls';
$version = (New-Object System.Net.WebClient).DownloadString("https://download.newrelic.com/install/newrelic-cli/currentVersion.txt").Trim();
(New-Object System.Net.WebClient).DownloadFile("https://download.newrelic.com/install/newrelic-cli/${version}/NewRelicCLIInstaller.msi", "$env:TEMP\NewRelicCLIInstaller.msi");
msiexec.exe /qn /i $env:TEMP\NewRelicCLIInstaller.msi | Out-Null;