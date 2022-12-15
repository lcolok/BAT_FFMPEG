@echo off
SET outputDir=%~dp1toPNG_768x768px
md "%outputDir%"
for %%i in (%*) do (
    ffmpeg -i %%i -s 768x768 "%outputDir%\%%~ni.png" -y
)
@REM pause