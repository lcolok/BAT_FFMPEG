@REM for %%a in ("*.mp4") do ffmpeg -i %%~na.mp4 -i %%~na.m4a -vcodec copy -acodec copy %%~na(combined).mp4
@echo off
echo %1
echo %2
ffmpeg -i %1 -i %2 -vcodec copy -acodec copy %~n1(combined).mp4
pause