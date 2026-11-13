
@echo sign binaries
@echo ================

"C:\Program Files (x86)\Windows Kits\8.0\bin\x86\signtool.exe" sign /t http://timestamp.digicert.com /n "Ping Castle SAS" /a bin\release\PingCastle*.exe

"C:\Program Files (x86)\Windows Kits\8.0\bin\x86\signtool.exe" sign  /as /d PingCastle /tr http://timestamp.digicert.com /td sha256 /fd sha256 /n "Ping Castle SAS" /a bin\release\PingCastle*.exe

@echo Done signing
:exit
pause
