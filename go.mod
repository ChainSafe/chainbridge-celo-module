module github.com/ChainSafe/chainbridge-celo-module

go 1.16

replace github.com/celo-org/celo-bls-go => github.com/celo-org/celo-bls-go v0.1.7

require (
	github.com/ChainSafe/chainbridge-core v0.0.0-00010101000000-000000000000
	github.com/celo-org/celo-blockchain v1.3.2
	github.com/celo-org/celo-bls-go v0.2.4
	github.com/go-kit/kit v0.9.0 // indirect
	github.com/golang/mock v1.4.4
	github.com/kr/pretty v0.2.0 // indirect
	github.com/mattn/go-colorable v0.1.2 // indirect
	github.com/pkg/errors v0.9.1
	github.com/rs/zerolog v1.21.0
	github.com/stretchr/testify v1.7.0
	github.com/syndtr/goleveldb v1.0.1-0.20200815110645-5c35d600f0ca
	golang.org/x/crypto v0.0.0-20200709230013-948cd5f35899 // indirect
	golang.org/x/net v0.0.0-20210119194325-5f4716e94777 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
)
