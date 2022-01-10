package main

import (
	"testing"

	"github.com/ChainSafe/chainbridge-celo-module/transaction"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/local"
	"github.com/ChainSafe/chainbridge-core/e2e/evm"
	"github.com/stretchr/testify/suite"
)

const CeloEndpoint1 = "ws://localhost:8546"
const CeloEndpoint2 = "ws://localhost:8548"

// Alice key is used by the relayer, Eve key is used as admin and depositter
func TestRunE2ETests(t *testing.T) {
	suite.Run(t, evm.SetupEVM2EVMTestSuite(transaction.NewCeloTransaction, transaction.NewCeloTransaction, CeloEndpoint1, CeloEndpoint2, local.EveKp))
}
