#!/usr/bin/env sh

# Run integration tests
go test -v ./test/...

# Run unit tests
cd form3
go test -v ./...