package admin

import (
	"github.com/ChainSafe/chainbridge-celo-module/transaction"
	bridgeContract "github.com/ChainSafe/chainbridge-core/chains/evm/calls/contracts/bridge"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/admin"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/initialize"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/logger"
	"github.com/spf13/cobra"
)

var AdminCeloCmd = &cobra.Command{
	Use:   "admin",
	Short: "Set of commands for executing various admin actions",
	Long:  "Set of commands for executing various admin actions",
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
		transactor, err := initialize.InitializeTransactor(
			gasPrice,
			transaction.NewCeloTransaction,
			client,
		)
		if err != nil {
			return err
		}
		return admin.PauseCmd(
			cmd,
			args,
			bridgeContract.NewBridgeContract(
				client,
				BridgeAddr,
				transactor,
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
		transactor, err := initialize.InitializeTransactor(
			gasPrice,
			transaction.NewCeloTransaction,
			client,
		)
		if err != nil {
			return err
		}
		return admin.UnpauseCmd(
			cmd,
			args,
			bridgeContract.NewBridgeContract(
				client,
				BridgeAddr,
				transactor,
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
