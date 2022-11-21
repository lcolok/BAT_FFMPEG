@echo off
SET outputDir=%~dp1toJPG
md "%outputDir%"
for %%i in (%*) do (
    ffmpeg -i %%i -y -q 1 "%outputDir%\%%~ni.jpg"
)
@REM pause