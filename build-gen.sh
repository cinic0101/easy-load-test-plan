#!/usr/bin/env bash

VERSION=0.3.0

# Mac OS X 64bit
GOOS=darwin GOARCH=amd64 go build -o out/ezgen-mac64-${VERSION} ez-plan-gen.go &&

# linux 64bit
GOOS=linux GOARCH=amd64 go build -o out/ezgen-linux64-${VERSION} ez-plan-gen.go &&

# linux 32bit
GOOS=linux GOARCH=386 go build -o out/ezgen-linux32-${VERSION} ez-plan-gen.go &&

# windows 64bit
GOOS=windows GOARCH=amd64 go build -o out/ezgen-win64-${VERSION}.exe ez-plan-gen.go &&

# windows 32bit
GOOS=windows GOARCH=386 go build -o out/ezgen-win32-${VERSION}.exe ez-plan-gen.go