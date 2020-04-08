
makecert.exe -r -pe -ss MY -sky exchange -n CN="New Relic CLI" CodeSign.cer
signtool sign /v /s MY /n "New Relic CLI" /t http://timestamp.verisign.com/scripts/timstamp.dll bin\x64\Release\NewRelicCLIInstaller.msi