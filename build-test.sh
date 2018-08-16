#!/usr/bin/env bash

VERSION=0.2.1

# Mac OS X 64bit
GOOS=darwin GOARCH=amd64 go build -o out/eztest-mac64-${VERSION} ez-test.go &&

# linux 64bit
GOOS=linux GOARCH=amd64 go build -o out/eztest-linux64-${VERSION} ez-test.go &&

# linux 32bit
GOOS=linux GOARCH=386 go build -o out/eztest-linux32-${VERSION} ez-test.go &&

# windows 64bit
GOOS=windows GOARCH=amd64 go build -o out/eztest-win64-${VERSION}.exe ez-test.go &&

# windows 32bit
GOOS=windows GOARCH=386 go build -o out/eztest-win32-${VERSION}.exe ez-test.go