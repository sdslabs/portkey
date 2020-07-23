 
#!/bin/bash

GO_FLAGS=${GO_FLAGS:-"-tags quic"}

pkgs=$(go list ./... | grep -v vendor)

echo "[*] Formatting code"

go vet ${pkgs} {$GO_FLAGS} > /dev/null 2>&1

vet_status="$?"

if [ "$vet_status" -ne 0 ]
then
	go fmt ${pkgs}
	echo "[-] Error while running golangci-lint formatters. Run 'make lint' to fix errors."
	exit 1
fi

# Run --fix with exit code 0 and don't display output of command
$(go env GOPATH)/bin/golangci-lint run --fix --issues-exit-code 0 > /dev/null 2>&1