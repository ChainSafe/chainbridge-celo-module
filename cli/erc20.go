package cli

import (
	"github.com/ChainSafe/chainbridge-celo-module/transaction"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/erc20"
	"github.com/spf13/cobra"
)

var ERC20CELOCMD = &cobra.Command{
	Use:   "erc20",
	Short: "ERC20-related instructions",
	Long:  "ERC20-related instructions",
}

var AddMinterCeloCMD = &cobra.Command{
	Use:   "add-minter",
	Short: "Add a minter to an Erc20 mintable contract",
	Long:  "Add a minter to an Erc20 mintable contract",
	RunE: func(cmd *cobra.Command, args []string) error {
		txFabric := transaction.NewCeloTransaction
		return erc20.AddMinterEVMCMD(cmd, args, txFabric)
	},
}

// var AllowanceCeloCMD = &cobra.Command{
// 	Use:   "allowance",
// 	Short: "Set a token contract as mintable/burnable",
// 	Long:  "Set a token contract as mintable/burnable in a handler",
// 	RunE: func(cmd *cobra.Command, args []string) error {
// 		txFabric := transaction.NewCeloTransaction
// 		return erc20.AllowanceEVMCMD(cmd, args, txFabric)
// 	},
// }

var DepositCeloCMD = &cobra.Command{
	Use:   "deposit",
	Short: "Initiate a transfer of ERC20 tokens",
	Long:  "Initiate a transfer of ERC20 tokens",
	RunE: func(cmd *cobra.Command, args []string) error {
		txFabric := transaction.NewCeloTransaction
		return erc20.DepositEVMCMD(cmd, args, txFabric)
	},
}

var MintCeloCMD = &cobra.Command{
	Use:   "mint",
	Short: "Mint tokens on an ERC20 mintable contract",
	Long:  "Mint tokens on an ERC20 mintable contract",
	RunE: func(cmd *cobra.Command, args []string) error {
		txFabric := transaction.NewCeloTransaction
		return erc20.MintEVMCMD(cmd, args, txFabric)
	},
}

// TODO:
// allowance
// approve
// balance

func init() {
	erc20.BindERC20AddMinterCLIFlags(AddMinterCeloCMD)
	erc20.BindERC20DepositCLIFlags(DepositCeloCMD)
	erc20.BindERC20MintCLIFlags(MintCeloCMD)
	// erc20.BindERC20AllowanceCLIFlags(AllowanceCeloCMD)
	ERC20CELOCMD.AddCommand(AddMinterCeloCMD, DepositCeloCMD, MintCeloCMD)
}
