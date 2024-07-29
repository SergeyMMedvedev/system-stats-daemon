
install-lint-deps:
	export GOROOT=$(go env GOROOT)
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.55.2

lint: install-lint-deps
	golangci-lint run ./...

test:
	go test -race -count 2 ./...

generate:
	rm -rf internal/pb
	mkdir -p internal/pb

	protoc \
	--proto_path=./api \
	--go_out=./ \
	--go-grpc_out=./ \
	api/*.proto