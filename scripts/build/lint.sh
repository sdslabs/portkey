#!/bin/bash

tags="quic"
pkgs=$(go list ./... | grep -v vendor)

echo "[*] Linting code for errors"

# Vet first since golangci-lint returns unclear errors if packages don't build
go vet --tags ${tags} ${pkgs} 

# Exit if vet throws an error
vet_status="$?"

if [ "$vet_status" -ne 0 ]
then
	exit "$vet_status"
fi

$(go env GOPATH)/bin/golangci-lint run --build-tags ${tags}