@ECHO OFF
@setlocal enableextensions
@cd /d "%~dp0"
 
SET PGPATH="D:/Program Files/PostgreSQL/10/bin/"
SET SVPATH=D:\neoreads\data\backup\
SET DBNAME=%1%
SET DBUSR=postgres
SET DBDUMP=neoreads_%2%.bak
 

if "%1%" == "" (
    echo " dbname is empty"
) else (

    if "%2%" == "" (
        echo " date is empty"
    ) else (

        echo "%PGPATH%psql.exe -h localhost -U postgres -d %DBNAME% < %SVPATH%%DBDUMP%"
        %PGPATH%psql.exe -h localhost -U postgres -d %DBNAME% < %SVPATH%%DBDUMP%
        echo Restore Complete %SVPATH%%DBDUMP%
    )
)
