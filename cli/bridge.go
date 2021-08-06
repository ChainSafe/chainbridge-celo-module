package cli

import (
	"github.com/ChainSafe/chainbridge-celo-module/transaction"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/bridge"
	"github.com/spf13/cobra"
)

var BridgeCELOCMD = &cobra.Command{
	Use:   "bridge",
	Short: "Bridge-related instructions",
	Long:  "Bridge-related instructions",
}

var registerResourceCeloCMD = &cobra.Command{
	Use:   "register-resource",
	Short: "Register a resource ID",
	Long:  "Register a resource ID with a contract address for a handler",
	RunE:  func(cmd *cobra.Command, args []string) error {
		txFabric := transaction.NewCeloTransaction
		return bridge.RegisterResourceEVMCMD(cmd, args, txFabric)
	},
}

var setBurnCeloCMD = &cobra.Command{
	Use:   "set-burn",
	Short: "Set a token contract as mintable/burnable",
	Long:  "Set a token contract as mintable/burnable in a handler",
	RunE:  func(cmd *cobra.Command, args []string) error {
		txFabric := transaction.NewCeloTransaction
		return bridge.SetBurnEVMCMD(cmd, args, txFabric)
	},
}

func init() {
	bridge.BindBridgeRegisterResourceCLIFlags(registerResourceCeloCMD)
	bridge.BindBridgeSetBurnCLIFlags(setBurnCeloCMD)
	BridgeCELOCMD.AddCommand(registerResourceCeloCMD, setBurnCeloCMD)
}
