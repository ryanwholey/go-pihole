default: test

.PHONY: testall test lint docs acceptance

test:
	go test ./...

acceptance:
	TEST_ACC=1 go test -race -v ./...

lint:
	golangci-lint run ./...

fmt:
	go fmt ./...

vet:
	go vet ./...	
