
.PHONY: help run build install license
all: help

## license: Adds license header to missing files.
license:
	@echo "  >  \033[32mAdding license headers...\033[0m "
	GO111MODULE=off go get -u github.com/google/addlicense
	addlicense -c "ChainSafe Systems" -f ./scripts/header.txt -y 2021 .

## license-check: Checks for missing license headers
license-check:
	@echo "  >  \033[Checking for license headers...\033[0m "
	GO111MODULE=off go get -u github.com/google/addlicense
	addlicense -check -c "ChainSafe Systems" -f ./scripts/header.txt -y 2021 .

local:
	go build -o chainbridge-celo-relayer e2e/celo-celo/example/main.go

e2e-setup:
	docker-compose --file=./e2e/celo-celo/docker-compose.e2e.yml up

e2e-test:
	./scripts/int_tests.sh

local-setup:
	./scripts/local_setup.sh
