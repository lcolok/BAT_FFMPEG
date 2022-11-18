@echo off
SET outputDir=%~dp1_hflipped
@REM md "%outputDir%"
for %%i in (%*) do (
    @REM ffmpeg -i %%i -vf "hflip" "%outputDir%\%%~ni_hflip.png" -y
    ffmpeg -i %%i -vf "hflip" "%%~ni_hflip.png" -y
)
@REM pause