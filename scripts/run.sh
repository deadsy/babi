#!/bin/bash
./genfft.py > fft.go
goimports -w .
go build
