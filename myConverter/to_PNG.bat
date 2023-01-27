@echo off
SET outputDir=%~dp1toPNG
md "%outputDir%"
for %%i in (%*) do (
    ffmpeg -i %%i -y "%outputDir%\%%~ni.png"
)
@REM pause