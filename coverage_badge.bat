@echo off
setlocal enabledelayedexpansion

REM Place coverage stats in coverage.txt
go tool cover -func=coverage.out > coverage.txt

REM Get the total coverage percentage
for /F "delims=" %%a in (coverage.txt) do (
   set "line="
   for %%b in (%%a) do set "line=!line!,%%b"
)
SET line=%line:,total:,(statements),=%
SET "cov=!line:%%=!"
REM echo %cov%

gopherbadger -md="README.md" -manualcov=%cov%

DEL coverage.txt
