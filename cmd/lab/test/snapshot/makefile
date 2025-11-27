# Simple Makefile to build and run the snapshot test

.PHONY: build run

# Target to build the executable
build:
	go build -o cmd/lab/test/snapshot-test/main cmd/lab/test/snapshot-test/main.go

# Target to run the executable with specified environment variables and flags
run: build
	KJ_TEST_DB_PATH=./files/database/kidung-jemaat.db KJ_TEST_MUSICXML_PATH=./files/scores/musicxml/ ./cmd/lab/test/snapshot-test/main -snapshot=gen -path=./cmd/lab/test/snapshot-test/goldenfiles
