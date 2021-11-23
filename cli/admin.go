package cli

import (
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/admin"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/logger"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmgaspricer"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmtransaction"
	"github.com/spf13/cobra"
)

var AdminCeloCmd = &cobra.Command{
	Use:   "admin",
	Short: "Admin-related instructions",
	Long:  "Admin-related instructions",
}

var pauseCmd = &cobra.Command{
	Use:   "pause",
	Short: "Pause deposits and proposals",
	Long:  "Pause deposits and proposals",
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.LoggerMetadata(cmd.Name(), cmd.Flags())
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return admin.PauseCmd(cmd, args, evmtransaction.NewTransaction, &evmgaspricer.LondonGasPriceDeterminant{})
	},
	Args: func(cmd *cobra.Command, args []string) error {
		err := admin.ValidatePauseCmdFlags(cmd, args)
		if err != nil {
			return err
		}

		admin.ProcessPauseCmdFlags(cmd, args)

		return nil
	},
}

var unpauseCmd = &cobra.Command{
	Use:   "unpause",
	Short: "Unpause deposits and proposals",
	Long:  "Unpause deposits and proposals",
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.LoggerMetadata(cmd.Name(), cmd.Flags())
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return admin.UnpauseCmd(cmd, args, evmtransaction.NewTransaction, &evmgaspricer.LondonGasPriceDeterminant{})
	},
	Args: func(cmd *cobra.Command, args []string) error {
		err := admin.ValidateUnpauseCmdFlags(cmd, args)
		if err != nil {
			return err
		}

		admin.ProcessUnpauseCmdFlags(cmd, args)

		return nil
	},
}

func init() {
	admin.BindPauseCmdFlags(pauseCmd)
	admin.BindUnpauseCmdFlags(unpauseCmd)

	AdminCeloCmd.AddCommand(
		pauseCmd,
		unpauseCmd,
	)
}
