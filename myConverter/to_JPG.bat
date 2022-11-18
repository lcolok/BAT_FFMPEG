@echo off
SET outputDir=%~dp1toJPG
md "%outputDir%"
for %%i in (%*) do (
    ffmpeg -i %%i -y "%outputDir%\%%~ni.jpg"
)
@REM pause