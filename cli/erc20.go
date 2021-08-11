package cli

import (
	"github.com/ChainSafe/chainbridge-celo-module/transaction"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/erc20"
	"github.com/spf13/cobra"
)

var ERC20CeloCmd = &cobra.Command{
	Use:   "erc20",
	Short: "erc20-related instructions",
	Long:  "erc20-related instructions",
}

var addMinterCmd = &cobra.Command{
	Use:   "add-minter",
	Short: "Add a minter to an Erc20 mintable contract",
	Long:  "Add a minter to an Erc20 mintable contract",
	RunE: func(cmd *cobra.Command, args []string) error {
		txFabric := transaction.NewCeloTransaction
		return erc20.AddMinterCmd(cmd, args, txFabric)
	},
}

var allowanceCmd = &cobra.Command{
	Use:   "allowance",
	Short: "Set a token contract as mintable/burnable",
	Long:  "Set a token contract as mintable/burnable in a handler",
	RunE: func(cmd *cobra.Command, args []string) error {
		txFabric := transaction.NewCeloTransaction
		return erc20.AllowanceCmd(cmd, args, txFabric)
	},
}

var approveCmd = &cobra.Command{
	Use:   "approve",
	Short: "Approve tokens in an ERC20 contract for transfer",
	Long:  "Approve tokens in an ERC20 contract for transfer",
	RunE: func(cmd *cobra.Command, args []string) error {
		txFabric := transaction.NewCeloTransaction
		return erc20.ApproveCmd(cmd, args, txFabric)
	},
}

var depositCmd = &cobra.Command{
	Use:   "deposit",
	Short: "Initiate a transfer of ERC20 tokens",
	Long:  "Initiate a transfer of ERC20 tokens",
	RunE: func(cmd *cobra.Command, args []string) error {
		txFabric := transaction.NewCeloTransaction
		return erc20.DepositCmd(cmd, args, txFabric)
	},
}

var mintCmd = &cobra.Command{
	Use:   "deposit",
	Short: "Initiate a transfer of ERC20 tokens",
	Long:  "Initiate a transfer of ERC20 tokens",
	RunE: func(cmd *cobra.Command, args []string) error {
		txFabric := transaction.NewCeloTransaction
		return erc20.MintCmd(cmd, args, txFabric)
	},
}

func init() {
	erc20.BindApproveCmdFlags(approveCmd)
	erc20.BindDepositCmdFlags(depositCmd)
	erc20.BindAddMinterCmdFlags(addMinterCmd)
	erc20.BindAllowanceCmdFlags(allowanceCmd)
	erc20.BindMintCmdFlags(mintCmd)
	ERC20CeloCmd.AddCommand(approveCmd, depositCmd, addMinterCmd, mintCmd)
}
