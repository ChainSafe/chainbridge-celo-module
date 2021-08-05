package cli

import (
	"github.com/ChainSafe/chainbridge-celo-module/transaction"
	coreDeployCLI "github.com/ChainSafe/chainbridge-core/chains/evm/cli/deploy"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/spf13/cobra"
)


var DeployCELO = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy smart contracts",
	Long:  "This command can be used to deploy all or some of the contracts required for bridging. Selection of contracts can be made by either specifying --all or a subset of flags",
	RunE: func(cmd *cobra.Command, args []string) error {
		txFabric := transaction.NewCeloTransaction
		return coreDeployCLI.DeployCLI(cmd, args, txFabric)
	},
}

func init() {
	flags.BindDeployEVMFlags(DeployCELO)
}
