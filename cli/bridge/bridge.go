package bridge

import (
	"fmt"

	"github.com/ChainSafe/chainbridge-celo-module/transaction"
	bridgeContract "github.com/ChainSafe/chainbridge-core/chains/evm/calls/contracts/bridge"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/bridge"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/initialize"

	"github.com/spf13/cobra"
)

var BridgeCeloCmd = &cobra.Command{
	Use:   "bridge",
	Short: "Set of commands for interacting with a bridge",
	Long:  "Set of commands for interacting with a bridge",
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

var registerResourceCmd = &cobra.Command{
	Use:   "register-resource",
	Short: "Register a resource ID",
	Long:  "The register-resource subcommand registers a resource ID with a contract address for a handler",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := initialize.InitializeClient(
			url,
			senderKeyPair,
		)
		if err != nil {
			return err
		}
		t, err := initialize.InitializeTransactor(gasPrice, transaction.NewCeloTransaction, c, prepare)
		if err != nil {
			return err
		}

		return bridge.RegisterResourceCmd(
			cmd,
			args,
			bridgeContract.NewBridgeContract(
				c,
				BridgeAddr,
				t,
			))
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
	Long:  "The set-burn subcommand sets a token contract as mintable/burnable in a handler",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := initialize.InitializeClient(
			url,
			senderKeyPair,
		)
		if err != nil {
			return err
		}
		t, err := initialize.InitializeTransactor(gasPrice, transaction.NewCeloTransaction, c, prepare)
		if err != nil {
			return err
		}

		return bridge.SetBurnCmd(
			cmd,
			args,
			bridgeContract.NewBridgeContract(
				c,
				BridgeAddr,
				t,
			))
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
	bridge.BindRegisterResourceFlags(registerResourceCmd)
	bridge.BindSetBurnFlags(setBurnCmd)

	BridgeCeloCmd.AddCommand(registerResourceCmd, setBurnCmd)
}
