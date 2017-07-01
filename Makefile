.PHONY: check fmt lint errcheck test build

SHELL:=/bin/bash

build: check
	go build github.com/Netflix/chaosmonkey/cmd/chaosmonkey

check: fmt lint errcheck

gofmt: fmt

fmt: 
	diff -u <(echo -n) <(gofmt -d `find . -name '*.go' | grep -Ev '/vendor/|/migration'`)

lint:
	go list ./... | grep -Ev '/vendor/|/migration' | xargs -L1 golint

errcheck:
	errcheck -ignore 'io:Close' -ignoretests `go list ./... | grep -v /vendor/`

test:
	go test -v  ./...

fix:
	gofmt -w `find . -name '*.go' | grep -Ev '/vendor/|/migration'`
