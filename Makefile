.PHONY: check fmt lint

SHELL:=/bin/bash

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


