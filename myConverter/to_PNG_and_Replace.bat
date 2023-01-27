@echo off
for %%i in (%*) do (
ffmpeg -i %%i -y "%~dpn1.png"
del %%i
)
@REM pause