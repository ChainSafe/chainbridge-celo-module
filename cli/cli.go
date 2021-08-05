package cli

import (
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/spf13/cobra"
)

var CeloRootCLI = &cobra.Command{
	Use:   "celo-cli",
	Short: "Celo CLI",
	Long:  "Root command for starting Celo CLI",
}

func init() {
	// persistent flags
	flags.BindEVMCLIFlags(CeloRootCLI)

	// add commands to celo-cli root
	// deploy
	CeloRootCLI.AddCommand(DeployCELO)

	// // admin
	// celoRootCLI.AddCommand(admin.AdminCmd)

	// // bridge
	// celoRootCLI.AddCommand(bridge.BridgeCmd)

	// // erc20
	// celoRootCLI.AddCommand(erc20.ERC20Cmd)

	// // erc721
	// celoRootCLI.AddCommand(erc721.ERC721Cmd)
}
