#!/usr/bin/env sh

# Run integration tests
go test -v ./test/...
if [ $? -ne 0 ]; then
    echo "Integration tests failed"
    exit 1
fi

# Run unit tests
cd form3
go test -v ./...
if [ $? -ne 0 ]; then
    echo "Unit tests failed"
    exit 1
fi