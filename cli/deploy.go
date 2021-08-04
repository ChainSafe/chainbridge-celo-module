package cli

import (
	"github.com/ChainSafe/chainbridge-celo-module/transaction"
	coreDeployCLI "github.com/ChainSafe/chainbridge-core/chains/evm/cli/deploy"
	"github.com/spf13/cobra"
)

var DeployCELO = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy smart contracts",
	Long:  "This command can be used to deploy all or some of the contracts required for bridging. Selection of contracts can be made by either specifying --all or a subset of flags",
	RunE:   func(cmd *cobra.Command, args []string) error {
		txFabric := transaction.NewCeloTransaction
		return coreDeployCLI.DeployCLI(cmd, args, txFabric)
	},
}

func init() {
	DeployCELO.Flags().Bool("bridge", false, "deploy bridge")
	DeployCELO.Flags().Bool("erc20Handler", false, "deploy ERC20 handler")
	//DeployEVM.Flags().Bool("erc721Handler", false, "deploy ERC721 handler")
	//DeployEVM.Flags().Bool("genericHandler", false, "deploy generic handler")
	DeployCELO.Flags().Bool("erc20", false, "deploy ERC20")
	DeployCELO.Flags().Bool("erc721", false, "deploy ERC721")
	DeployCELO.Flags().Bool("all", false, "deploy all")
	DeployCELO.Flags().Int64("relayerThreshold", 1, "number of votes required for a proposal to pass")
	DeployCELO.Flags().String("chainId", "1", "chain ID for the instance")
	DeployCELO.Flags().StringSlice("relayers", []string{}, "list of initial relayers")
	DeployCELO.Flags().String("fee", "0", "fee to be taken when making a deposit (in ETH, decimas are allowed)")
	DeployCELO.Flags().String("bridgeAddress", "", "bridge contract address. Should be provided if handlers are deployed separately")
	DeployCELO.Flags().String("erc20Symbol", "", "ERC20 contract symbol")
	DeployCELO.Flags().String("erc20Name", "", "ERC20 contract name")
	DeployCELO.Flags().String("url", "ws://localhost:8545", "node url")
}

