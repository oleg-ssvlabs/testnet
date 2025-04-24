EXEC_DIRECTORY=./cmd/cli
BINARY_PATH=${EXEC_DIRECTORY}/bin/testnet
BINARY_DIR=${EXEC_DIRECTORY}/bin

CONFIG_FILE=./internal/blockchain/params.yaml

.PHONY: default
default: run

.PHONY: blockchain-clean
blockchain-clean:
	kurtosis clean -a

.PHONY: build
build:
	go build -o ${BINARY_PATH} ${EXEC_DIRECTORY}
	@cp $(CONFIG_FILE) $(BINARY_DIR)/

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