package main

import (
	"testing"

	"github.com/ChainSafe/chainbridge-celo-module/transaction"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/evmclient"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/local"
	"github.com/ChainSafe/chainbridge-core/e2e/evm"
	"github.com/stretchr/testify/suite"
)

const ETHEndpoint1 = "ws://localhost:8648"
const ETHEndpoint2 = "ws://localhost:8546"

// Alice key is used by the relayer, Eve key is used as admin and depositter
func TestRunE2ETests(t *testing.T) {
	celo1, err := evmclient.NewEVMClientFromParams(ETHEndpoint1, local.EveKp.PrivateKey())
	if err != nil {
		panic(err)
	}

	celo2, err := evmclient.NewEVMClientFromParams(ETHEndpoint2, local.EveKp.PrivateKey())
	if err != nil {
		panic(err)
	}

	suite.Run(t, evm.SetupEVM2EVMTestSuite(
		transaction.NewCeloTransaction,
		transaction.NewCeloTransaction,
		celo1,
		celo2,
		local.DefaultRelayerAddresses,
		local.DefaultRelayerAddresses,
	))
}
