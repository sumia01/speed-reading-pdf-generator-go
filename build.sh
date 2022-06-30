#!/usr/bin/env bash

env GOOS=darwin GOARCH=amd64 go build -o ./bin/osx_x64
env GOOS=darwin GOARCH=arm64 go build -o ./bin/osx_arm
env GOOS=linux GOARCH=amd64 go build -o ./bin/linux_x64
env GOOS=windows GOARCH=amd64 go build -o ./bin/windows_x64.exe
