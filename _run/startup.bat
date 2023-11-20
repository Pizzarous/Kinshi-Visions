cd ..
call start-ai.bat
call git reset --hard origin/main
call git pull --all
call go build
call .\kinshi_vision_bot.exe