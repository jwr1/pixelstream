#!/usr/bin/env bash

GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" .
tar -czf pixelstream-linux-amd64.tar.gz  pixelstream
rm pixelstream

GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" .
zip pixelstream-windows-amd64.zip pixelstream.exe
rm pixelstream.exe
