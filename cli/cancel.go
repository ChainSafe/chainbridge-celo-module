package cli

//var CancelProposalEVM = &cobra.Command{
//	Use: "cancel-proposal",
//	Short: "deploy smart contracts",
//	Long: "This command can be used to deploy all or some of the contracts required for bridging. Selection of contracts can be made by either specifying --all or a subset of flags",
//	Run: func(cmd *cobra.Command, args []string) {
//		err := ReturnCancelProposalCLICELO()(cmd, args)
//		log.Err(err)
//	},
//}
//
//
//func ReturnCancelProposalCLICELO() func(cmd *cobra.Command, args []string) error {
//	txFabric := transaction.NewCeloTransaction
//
//
//	return func(cmd *cobra.Command, args[] string) error {
//		err = bridge.CancelProposal(ethClient, bridgeAddress, uint8(chainID), depositNonce, dataHashBytes, fabric)
//		if err != nil{
//			return err
//		}
//		log.Info().Msgf("Setting proposal with chain ID %v and deposit nonce %v status to 'Cancelled", chainID, depositNonce)
//		return nil
//	}
//}
