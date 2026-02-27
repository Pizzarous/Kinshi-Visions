@echo off
cd ..
call git reset --hard origin/master
if errorlevel 1 (echo git reset failed & pause & exit /b 1)
call git pull --all
if errorlevel 1 (echo git pull failed & pause & exit /b 1)
call go build
if errorlevel 1 (echo go build failed & pause & exit /b 1)
echo Build successful!
pause
