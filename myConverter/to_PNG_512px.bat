@echo off
SET outputDir=%~dp1toPNG_512x512px
md "%outputDir%"
for %%i in (%*) do (
    ffmpeg -i %%i -s 512x512 "%outputDir%\%%~ni.png" -y
)
@REM pause