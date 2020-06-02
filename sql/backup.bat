@ECHO OFF
@setlocal enableextensions
@cd /d "%~dp0"
 
SET PGPATH="D:/Program Files/PostgreSQL/10/bin/"
SET SVPATH=D:\neoreads\data\backup\
SET DBNAME=neoreads
SET DBUSR=postgres
FOR /F "TOKENS=1,2,3 DELIMS=/ " %%i IN ('DATE /T') DO SET d=%%i%%j%%k
 
SET DBDUMP=%DBNAME%_%d%.bak
:: @ECHO OFF
%PGPATH%pg_dump -h localhost -U postgres %DBNAME% > %SVPATH%%DBDUMP%

echo Backup Taken Complete %SVPATH%%DBDUMP%

:: forfiles /p %SVPATH% /d -5 /c "cmd /c echo deleting @file ... && del /f @path"