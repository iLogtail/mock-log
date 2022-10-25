#!/usr/bin/env bash
all: bin/mock_log

bin/mock_log: main.go
	go build -o bin/mock_log main.go

clean:
	go clean
	rm bin/*

docker: bin/mock_log
	docker build -t mock_log:local .
