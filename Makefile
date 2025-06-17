EXEC_DIRECTORY=./cmd/testnet
BINARY_PATH=${EXEC_DIRECTORY}/bin/testnet
BINARY_DIR=${EXEC_DIRECTORY}/bin
OBSERVABILITY_LABEL=stack=localnet-observability

.PHONY: default
default: run

.PHONY: kurtosis-clean
kurtosis-clean:
	kurtosis clean -a

.PHONY: kurtosis-show
kurtosis-show:
	kurtosis enclave inspect localnet

.PHONY: build
build:
	go build -o ${BINARY_PATH} ${EXEC_DIRECTORY}

.PHONY: run
run: build
	${BINARY_PATH}

.PHONY: observability-show
observability-show:
	docker ps -a --filter "label=${OBSERVABILITY_LABEL}"

.PHONY: observability-clean
observability-clean:
	docker ps -aq --filter "label=${OBSERVABILITY_LABEL}" | xargs -r docker rm -f

.PHONY: clean
clean: observability-clean kurtosis-clean
	
.PHONY: test
test:
	go test ./...

.PHONY: lint
lint:
	golangci-lint run -v ./...