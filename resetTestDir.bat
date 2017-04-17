@echo off

echo removing old test folders
rd /s /q "test/"

echo creating new test folders
robocopy testSrc/ test/ /S >nul 2>&1
