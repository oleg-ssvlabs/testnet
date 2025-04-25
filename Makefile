EXEC_DIRECTORY=./cmd/testnet
BINARY_PATH=${EXEC_DIRECTORY}/bin/testnet
BINARY_DIR=${EXEC_DIRECTORY}/bin

.PHONY: default
default: run

.PHONY: blockchain-clean
blockchain-clean:
	kurtosis clean -a

.PHONY: build
build:
	go build -o ${BINARY_PATH} ${EXEC_DIRECTORY}

.PHONY: run
run: build
	${BINARY_PATH}

.PHONY: clean
clean:
	go clean
	
.PHONY: test
test:
	go test ./...

.PHONY: lint
lint:
	golangci-lint run -v ./...