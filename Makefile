default: test

updatedeps:
	go get -v -u ./...

build: 
	go build -o bin/go-test-generate

test: 
	go test -v -parallel 5

test-race:
	go test -v -race -parallel 5

test-all: vet lint test cover

vet:
	@go get golang.org/x/tools/cmd/vet
	go tool vet .

lint:
	@go get github.com/golang/lint/golint
	golint ./...

# cover shows test coverages
cover:
	@go get golang.org/x/tools/cmd/cover		
	go test -coverprofile=cover.out
	go tool cover -html cover.out
	rm cover.out
