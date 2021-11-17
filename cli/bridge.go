package cli

import (
	"github.com/ChainSafe/chainbridge-celo-module/transaction"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/bridge"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmgaspricer"
	"github.com/spf13/cobra"
)

var BridgeCeloCmd = &cobra.Command{
	Use:   "bridge",
	Short: "Bridge-related instructions",
	Long:  "Bridge-related instructions",
}

var registerResourceCmd = &cobra.Command{
	Use:   "register-resource",
	Short: "Register a resource ID",
	Long:  "Register a resource ID with a contract address for a handler",
	RunE: func(cmd *cobra.Command, args []string) error {
		txFabric := transaction.NewCeloTransaction
		return bridge.RegisterResourceCmd(cmd, args, txFabric, &evmgaspricer.StaticGasPriceDeterminant{})
	},
	Args: func(cmd *cobra.Command, args []string) error {
		err := bridge.ValidateRegisterResourceFlags(cmd, args)
		if err != nil {
			return err
		}

		err = bridge.ProcessRegisterResourceFlags(cmd, args)
		return err
	},
}

var setBurnCmd = &cobra.Command{
	Use:   "set-burn",
	Short: "Set a token contract as mintable/burnable",
	Long:  "Set a token contract as mintable/burnable in a handler",
	RunE: func(cmd *cobra.Command, args []string) error {
		txFabric := transaction.NewCeloTransaction
		return bridge.SetBurnCmd(cmd, args, txFabric, &evmgaspricer.StaticGasPriceDeterminant{})
	},
	Args: func(cmd *cobra.Command, args []string) error {
		err := bridge.ValidateSetBurnFlags(cmd, args)
		if err != nil {
			return err
		}

		bridge.ProcessSetBurnFlags(cmd, args)
		return nil
	},
}

// TODO:
// cancel-proposal
// query-proposal
// query-resource
// register-generic-resource

func init() {
	bridge.BindRegisterResourceCmdFlags(registerResourceCmd)
	bridge.BindSetBurnCmdFlags(setBurnCmd)

	BridgeCeloCmd.AddCommand(registerResourceCmd, setBurnCmd)
}
