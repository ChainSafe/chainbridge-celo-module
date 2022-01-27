package admin

import (
	"fmt"

	"github.com/ChainSafe/chainbridge-celo-module/transaction"
	bridgeContract "github.com/ChainSafe/chainbridge-core/chains/evm/calls/contracts/bridge"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/admin"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/initialize"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/logger"
	"github.com/spf13/cobra"
)

var AdminCeloCmd = &cobra.Command{
	Use:   "admin",
	Short: "Set of commands for executing various admin actions",
	Long:  "Set of commands for executing various admin actions",
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

var pauseCmd = &cobra.Command{
	Use:   "pause",
	Short: "Pause deposits and proposals",
	Long:  "The pause subcommand pauses deposits and proposals",
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.LoggerMetadata(cmd.Name(), cmd.Flags())
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := initialize.InitializeClient(
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
		return admin.PauseCmd(
			cmd,
			args,
			bridgeContract.NewBridgeContract(
				client,
				BridgeAddr,
				t,
			))
	},
	Args: func(cmd *cobra.Command, args []string) error {
		err := admin.ValidatePauseFlags(cmd, args)
		if err != nil {
			return err
		}

		admin.ProcessPauseFlags(cmd, args)

		return nil
	},
}

var unpauseCmd = &cobra.Command{
	Use:   "unpause",
	Short: "Unpause deposits and proposals",
	Long:  "The unpause subcommand unpauses deposits and proposals",
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.LoggerMetadata(cmd.Name(), cmd.Flags())
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := initialize.InitializeClient(
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
		return admin.UnpauseCmd(
			cmd,
			args,
			bridgeContract.NewBridgeContract(
				client,
				BridgeAddr,
				t,
			))
	},
	Args: func(cmd *cobra.Command, args []string) error {
		err := admin.ValidateUnpauseFlags(cmd, args)
		if err != nil {
			return err
		}

		admin.ProcessUnpauseFlags(cmd, args)

		return nil
	},
}

func init() {
	admin.BindPauseFlags(pauseCmd)
	admin.BindUnpauseFlags(unpauseCmd)

	AdminCeloCmd.AddCommand(
		pauseCmd,
		unpauseCmd,
	)
}
