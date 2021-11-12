module github.com/ChainSafe/chainbridge-celo-module

go 1.15

replace github.com/ChainSafe/chainbridge-core => ../chainbridge-core

require (
	github.com/ChainSafe/chainbridge-core v0.0.0-20210702085934-c073bc8c16a4
	github.com/ethereum/go-ethereum v1.10.9
	github.com/pierrec/xxHash v0.1.5 // indirect
	github.com/rs/zerolog v1.25.0
	github.com/spf13/cobra v1.2.1
	github.com/status-im/keycard-go v0.0.0-20211004132608-c32310e39b86
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519
)
