#!/usr/bin/env bash
all:
	go build -mod="vendor" -o bin/mock_log main.go

clean:
	go clean
	rm bin/*