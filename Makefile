
install-lint-deps:
	export GOROOT=$(go env GOROOT)
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.55.2

lint: install-lint-deps
	golangci-lint run ./...

# Все тесты кроме интеграционных
test:
	go test $(shell go list ./... | grep -v github.com/SergeyMMedvedev/system-stats-daemon/integration) -race -count 2 ./...

# Интеграционные тесты
test-integration:
	go test -timeout 10m -count 1 ./integration/integration_test.go


generate:
	rm -rf internal/pb
	mkdir -p internal/pb

	protoc \
	--proto_path=./api \
	--go_out=./ \
	--go-grpc_out=./ \
	api/*.proto