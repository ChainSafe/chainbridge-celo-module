module github.com/ChainSafe/chainbridge-celo-module

go 1.15

replace github.com/ChainSafe/chainbridge-core => ../chainbridge-core

require (
	github.com/ChainSafe/chainbridge-core v0.0.0-20210702085934-c073bc8c16a4
	github.com/ethereum/go-ethereum v1.10.9
	github.com/rs/zerolog v1.23.0
	github.com/spf13/cobra v0.0.3
	github.com/status-im/keycard-go v0.0.0-20200402102358-957c09536969
	golang.org/x/crypto v0.0.0-20210616213533-5ff15b29337e
)
