ffmpeg ^
-i %1 ^
-f mp4 ^
-c:v libx264 -preset veryfast -profile:v baseline -level 1.3 -crf 15.0 ^
-c:a aac -ac 2 -b:a 192k ^
-sn -map_metadata -1 -map_chapters -1 ^
-y ^
"[н╒пе]%~n1.mp4"
pause
