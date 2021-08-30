[Net.ServicePointManager]::SecurityProtocol = 'tls12, tls';
$version = (New-Object System.Net.WebClient).DownloadString("https://download.newrelic.com/install/newrelic-cli/currentVersion.txt").Trim();
(New-Object System.Net.WebClient).DownloadFile("https://download.newrelic.com/install/newrelic-cli/${version}/NewRelicCLIInstaller.msi", "$env:TEMP\NewRelicCLIInstaller.msi");
msiexec.exe /qn /i $env:TEMP\NewRelicCLIInstaller.msi | Out-Null;