module github.com/ChainSafe/chainbridge-celo-module

go 1.16

replace (
	github.com/ChainSafe/chainbridge-core v0.0.0-20210602125535-8f78a5e6de69 => ../chainbridge-core
	github.com/celo-org/celo-bls-go => github.com/celo-org/celo-bls-go v0.1.7
)

require (
	github.com/ChainSafe/chainbridge-core v0.0.0-20210602125535-8f78a5e6de69
	github.com/celo-org/celo-blockchain v1.3.2
	github.com/celo-org/celo-bls-go v0.2.4
	github.com/golang/mock v1.4.4
	github.com/pkg/errors v0.9.1
	github.com/rs/zerolog v1.21.0
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	github.com/syndtr/goleveldb v1.0.1-0.20200815110645-5c35d600f0ca
)
