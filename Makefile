default: test

updatedeps:
	go get -v -u ./...

build: 
	go build -o bin/gotests

test: 
	go test -v -parallel 5

test-race:
	go test -v -race -parallel 5

test-all: vet lint test cover

vet:
	@go get golang.org/x/tools/cmd/vet
	@go tool vet *.go

lint:
	@go get github.com/golang/lint/golint
	golint . | grep -v "comment"

# cover shows test coverages
cover:
	@go get golang.org/x/tools/cmd/cover		
	go test -coverprofile=cover.out
	go tool cover -html cover.out
	rm cover.out

generate: build
	@go generate

doc: generate	
	@open http://127.0.0.1:9999/pkg/github.com/tcnksm/gotests/
	godoc -http=:9999
