@echo off
SET outputDir=%~dp1_hflipped
md "%outputDir%"
for %%i in (%*) do (
    ffmpeg -i %%i -vf "hflip" "%outputDir%\%%~ni_hflip.png" -y
)
@REM pause