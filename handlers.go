package client

import (
	"math/big"
)


type ValidatorsAggregator interface {
	GetAPKForBlock(block *big.Int, chainID uint8, epochSize uint64) ([]byte, error)
}

//func HandleErc20DepositedEventCelo(sourceID, destId uint8, nonce uint64, handlerContractAddress string, backend listener.ChainReader) (*relayer.Message, error)  {
//	contract, err := erc20Handler.NewERC20HandlerCaller(common.HexToAddress(handlerContractAddress), backend)
//	if err != nil {
//		return nil, err
//	}
//	record, err := contract.GetDepositRecord(&bind.CallOpts{}, uint64(nonce), destId)
//	if err != nil {
//		return nil, err
//	}
//	m := &relayer.Message{
//		Source:       sourceID,
//		Destination:  destId,
//		Type:         relayer.FungibleTransfer,
//		DepositNonce: nonce,
//		ResourceId:   record.ResourceID,
//		Payload: []interface{}{
//			record.Amount.Bytes(),
//			record.DestinationRecipientAddress,
//		},
//	}

	//b := backend.(celo.ChainReader)
	//blockData, err := b.BlockByNumber(context.Background(), block)
	//if err != nil {
	//	return nil, err
	//}
	//trie, err := txtrie.CreateNewTrie(blockData.TxHash(), blockData.Transactions())
	//if err != nil {
	//	return nil, err
	//}
	//apk, err := l.valsAggr.GetAPKForBlock(block, uint8(l.cfg.ID), l.cfg.EpochSize)
	//if err != nil {
	//	return nil, err
	//
	//}
	//keyRlp, err := rlp.EncodeToBytes(txIndex)
	//if err != nil {
	//	return nil, fmt.Errorf("encoding TxIndex to rlp: %w", err)
	//}
	//proof, key, err := txtrie.RetrieveProof(trie, keyRlp)
	//if err != nil {
	//	return nil, err
	//}
	//m.SVParams = &SignatureVerification{AggregatePublicKey: apk, BlockHash: blockData.Header().Hash(), Signature: blockData.EpochSnarkData().Signature}
	//m.MPParams = &MerkleProof{TxRootHash: sliceTo32Bytes(blockData.TxHash().Bytes()), Nodes: proof, Key: key}
	//return m, nil
//}
