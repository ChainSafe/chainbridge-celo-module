package deploy

import (
	"github.com/ChainSafe/chainbridge-celo-module/transaction"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/evmgaspricer"
	coreDeployCLI "github.com/ChainSafe/chainbridge-core/chains/evm/cli/deploy"
	"github.com/spf13/cobra"
)

var DeployCeloCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy smart contracts",
	Long:  "This command can be used to deploy all or some of the contracts required for bridging. Selection of contracts can be made by either specifying --all or a subset of flags",
	RunE: func(cmd *cobra.Command, args []string) error {
		txFabric := transaction.NewCeloTransaction
		return coreDeployCLI.DeployCLI(cmd, args, txFabric, &evmgaspricer.StaticGasPriceDeterminant{})
	},
	Args: func(cmd *cobra.Command, args []string) error {
		err := coreDeployCLI.ValidateDeployFlags(cmd, args)
		if err != nil {
			return err
		}
		err = coreDeployCLI.ProcessDeployFlags(cmd, args)
		return err
	},
}

func init() {
	coreDeployCLI.BindDeployEVMFlags(DeployCeloCmd)
}
