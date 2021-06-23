module github.com/ChainSafe/chainbridge-celo-module

go 1.16

replace github.com/ChainSafe/chainbridge-core => /var/www/ChainSafe/rnd/chainbridge-core

replace github.com/celo-org/celo-bls-go => github.com/celo-org/celo-bls-go v0.1.7

require (
	github.com/ChainSafe/chainbridge-core v0.0.0-00010101000000-000000000000
	github.com/ethereum/go-ethereum v1.10.4
	github.com/rs/zerolog v1.23.0
	github.com/status-im/keycard-go v0.0.0-20200402102358-957c09536969
)
