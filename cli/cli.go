package cli

import (
	"github.com/ChainSafe/chainbridge-celo-module/cli/admin"
	"github.com/ChainSafe/chainbridge-celo-module/cli/bridge"
	"github.com/ChainSafe/chainbridge-celo-module/cli/deploy"
	"github.com/ChainSafe/chainbridge-celo-module/cli/erc20"
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
	CeloRootCLI.AddCommand(deploy.DeployCeloCmd)

	// // admin
	CeloRootCLI.AddCommand(admin.AdminCeloCmd)

	// // bridge
	CeloRootCLI.AddCommand(bridge.BridgeCeloCmd)

	// // erc20
	CeloRootCLI.AddCommand(erc20.ERC20CeloCmd)

	// // erc721
	// celoRootCLI.AddCommand(erc721.ERC721Cmd)
}
