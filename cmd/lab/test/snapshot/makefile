# Simple Makefile to build and run the snapshot test

.PHONY: build run

# Target to build the executable
build:
	go build -o cmd/lab/test/snapshot/main cmd/lab/test/snapshot/main.go

# Target to run the executable with specified environment variables and flags
gen: build
	KJ_TEST_DB_PATH=./files/database/kidung-jemaat.db KJ_TEST_MUSICXML_PATH=./files/scores/musicxml/ ./cmd/lab/test/snapshot/main -snapshot=gen -path=./cmd/lab/test/snapshot

assert: build
	KJ_TEST_DB_PATH=./files/database/kidung-jemaat.db KJ_TEST_MUSICXML_PATH=./files/scores/musicxml/ ./cmd/lab/test/snapshot/main -snapshot=assert -path=./cmd/lab/test/snapshot