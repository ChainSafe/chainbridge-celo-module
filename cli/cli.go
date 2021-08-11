package cli

import (
	evmCLI "github.com/ChainSafe/chainbridge-core/chains/evm/cli"
	"github.com/spf13/cobra"
)

var CeloRootCLI = &cobra.Command{
	Use:   "celo-cli",
	Short: "Celo CLI",
	Long:  "Root command for starting Celo CLI",
}

func init() {
	// persistent flags
	evmCLI.BindEVMCLIFlags(CeloRootCLI)

	// add commands to celo-cli root
	// deploy
	CeloRootCLI.AddCommand(DeployCeloCmd)

	// // admin
	// celoRootCLI.AddCommand(admin.AdminCmd)

	// // bridge
	CeloRootCLI.AddCommand(BridgeCeloCmd)

	// // erc20
	CeloRootCLI.AddCommand(ERC20CeloCmd)

	// // erc721
	// celoRootCLI.AddCommand(erc721.ERC721Cmd)
}
