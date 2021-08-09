package cli

import (
	"github.com/ChainSafe/chainbridge-celo-module/transaction"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/erc20"
	"github.com/spf13/cobra"
)

var ERC20CeloCMD = &cobra.Command{
	Use: "erc20",
	Short: "erc20-related instructions",
	Long:  "erc20-related instructions",
}

var approveCmd = &cobra.Command{
	Use:   "approve",
	Short: "Approve tokens in an ERC20 contract for transfer",
	Long:  "Approve tokens in an ERC20 contract for transfer",
	RunE:   func(cmd *cobra.Command, args []string) error {
		txFabric := transaction.NewCeloTransaction
		return erc20.Approve(cmd, args, txFabric)
	},
}

var deposit = &cobra.Command{
	Use: "deposit",
	Short: "Deposit",
	Long: "Deposit",
	RunE:   func(cmd *cobra.Command, args []string) error {
		txFabric := transaction.NewCeloTransaction
		return erc20.DepositCMD(cmd, args, txFabric)
	},

}

func init() {
	erc20.BindApproveCLIFlags(approveCmd)
	erc20.BindDepositCMDFlags(deposit)
	ERC20CeloCMD.AddCommand(approveCmd, deposit)
}