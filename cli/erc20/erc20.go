package erc20

import (
	"fmt"

	"github.com/ChainSafe/chainbridge-celo-module/transaction"
	bridgeContract "github.com/ChainSafe/chainbridge-core/chains/evm/calls/contracts/bridge"
	erc20Contract "github.com/ChainSafe/chainbridge-core/chains/evm/calls/contracts/erc20"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/erc20"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/initialize"
	"github.com/spf13/cobra"
)

var ERC20CeloCmd = &cobra.Command{
	Use:   "erc20",
	Short: "Set of commands for interacting with an ERC20 contract",
	Long:  "Set of commands for interacting with an ERC20 contract",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		// fetch global flag values
		url, gasLimit, gasPrice, senderKeyPair, prepare, err = flags.GlobalFlagValues(cmd)
		if err != nil {
			return fmt.Errorf("could not get global flags: %v", err)
		}
		return nil
	},
}

var addMinterCmd = &cobra.Command{
	Use:   "add-minter",
	Short: "Add a new ERC20 minter",
	Long:  "The add-minter subcommand adds a minter to an ERC20 mintable contract",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := initialize.InitializeClient(url, senderKeyPair)
		if err != nil {
			return err
		}
		t, err := initialize.InitializeTransactor(gasPrice, transaction.NewCeloTransaction, c, prepare)
		if err != nil {
			return err
		}
		return erc20.AddMinterCmd(
			cmd,
			args,
			erc20Contract.NewERC20Contract(
				c,
				Erc20Addr,
				t,
			))
	},
	Args: func(cmd *cobra.Command, args []string) error {
		err := erc20.ValidateAddMinterFlags(cmd, args)
		if err != nil {
			return err
		}
		erc20.ProcessAddMinterFlags(cmd, args)
		return nil
	},
}

var allowanceCmd = &cobra.Command{
	Use:   "get-allowance",
	Short: "Get the allowance of a spender for an address",
	Long:  "The get-allowance subcommand returns the allowance of a spender for an address",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := initialize.InitializeClient(url, senderKeyPair)
		if err != nil {
			return err
		}
		t, err := initialize.InitializeTransactor(gasPrice, transaction.NewCeloTransaction, c, prepare)
		if err != nil {
			return err
		}
		return erc20.GetAllowanceCmd(
			cmd,
			args,
			erc20Contract.NewERC20Contract(
				c,
				Erc20Addr,
				t,
			))
	},
}

var approveCmd = &cobra.Command{
	Use:   "approve",
	Short: "Approve an ERC20 tokens",
	Long:  "The approve subcommand approves tokens in an ERC20 contract for transfer",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := initialize.InitializeClient(url, senderKeyPair)
		if err != nil {
			return err
		}
		t, err := initialize.InitializeTransactor(gasPrice, transaction.NewCeloTransaction, c, prepare)
		if err != nil {
			return err
		}
		return erc20.ApproveCmd(
			cmd,
			args,
			erc20Contract.NewERC20Contract(
				c,
				Erc20Addr,
				t,
			))
	},
	Args: func(cmd *cobra.Command, args []string) error {
		err := erc20.ValidateApproveFlags(cmd, args)
		if err != nil {
			return err
		}

		err = erc20.ProcessApproveFlags(cmd, args)
		return err
	},
}

var depositCmd = &cobra.Command{
	Use:   "deposit",
	Short: "Deposit an ERC20 token",
	Long:  "The deposit subcommand creates a new ERC20 token deposit on the bridge contract",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := initialize.InitializeClient(url, senderKeyPair)
		if err != nil {
			return err
		}
		t, err := initialize.InitializeTransactor(gasPrice, transaction.NewCeloTransaction, c, prepare)
		if err != nil {
			return err
		}
		return erc20.DepositCmd(
			cmd,
			args,
			bridgeContract.NewBridgeContract(
				c,
				BridgeAddr,
				t,
			))
	},
	Args: func(cmd *cobra.Command, args []string) error {
		err := erc20.ValidateDepositFlags(cmd, args)
		if err != nil {
			return err
		}

		err = erc20.ProcessDepositFlags(cmd, args)
		if err != nil {
			return err
		}
		return nil
	},
}

var mintCmd = &cobra.Command{
	Use:   "mint",
	Short: "Mint an ERC20 token",
	Long:  "The mint subcommand mints a token on an ERC20 mintable contract",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := initialize.InitializeClient(url, senderKeyPair)
		if err != nil {
			return err
		}
		t, err := initialize.InitializeTransactor(gasPrice, transaction.NewCeloTransaction, c, prepare)
		if err != nil {
			return err
		}
		return erc20.MintCmd(
			cmd,
			args,
			erc20Contract.NewERC20Contract(
				c,
				Erc20Addr,
				t,
			))
	},
	Args: func(cmd *cobra.Command, args []string) error {
		err := erc20.ValidateMintFlags(cmd, args)
		if err != nil {
			return err
		}

		err = erc20.ProcessMintFlags(cmd, args)
		return err
	},
}

var balanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "Query an ERC20 token balance",
	Long:  "The balance subcommand queries the balance of an account in an ERC20 contract",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := initialize.InitializeClient(url, senderKeyPair)
		if err != nil {
			return err
		}
		t, err := initialize.InitializeTransactor(gasPrice, transaction.NewCeloTransaction, c, prepare)
		if err != nil {
			return err
		}
		return erc20.BalanceCmd(
			cmd,
			args,
			erc20Contract.NewERC20Contract(
				c,
				Erc20Addr,
				t,
			))
	},
	Args: func(cmd *cobra.Command, args []string) error {
		err := erc20.ValidateBalanceFlags(cmd, args)
		if err != nil {
			return err
		}

		erc20.ProcessBalanceFlags(cmd, args)
		return nil
	},
}

func init() {
	erc20.BindApproveFlags(approveCmd)
	erc20.BindDepositFlags(depositCmd)
	erc20.BindAddMinterFlags(addMinterCmd)
	erc20.BindGetAllowanceFlags(allowanceCmd)
	erc20.BindMintFlags(mintCmd)
	erc20.BindBalanceFlags(balanceCmd)
	ERC20CeloCmd.AddCommand(
		approveCmd,
		depositCmd,
		addMinterCmd,
		allowanceCmd,
		mintCmd,
		balanceCmd,
	)
}
