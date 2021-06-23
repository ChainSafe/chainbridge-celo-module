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
//
//func (c *CeloClient) ReturnErc20HandlerFabric() listener.EventHandler {
//	return func(sourceID, destId uint8, nonce uint64, handlerContractAddress string) (*relayer.Message, error) {
//		contract, err := erc20Handler.NewERC20HandlerCaller(common.HexToAddress(handlerContractAddress), c.Client)
//		if err != nil {
//			return nil, err
//		}
//		record, err := contract.GetDepositRecord(&bind.CallOpts{}, uint64(nonce), uint8(destId))
//		if err != nil {
//			return nil, err
//		}
//		return &relayer.Message{
//			Source:       sourceID,
//			Destination:  destId,
//			DepositNonce: nonce,
//			ResourceId:   record.ResourceID,
//			Type:         relayer.FungibleTransfer,
//			Payload: []interface{}{
//				record.Amount.Bytes(),
//				record.DestinationRecipientAddress,
//			},
//		}, nil
//	}
//}
//
//func (c *CeloClient) ReturnErc721HandlerFabric() listener.EventHandler {
//	return func(sourceID, destId uint8, nonce uint64, handlerContractAddress string) (*relayer.Message, error) {
//		contract, err := erc721Handler.NewERC721HandlerCaller(common.HexToAddress(handlerContractAddress), c.Client)
//		if err != nil {
//			return nil, err
//		}
//		record, err := contract.GetDepositRecord(&bind.CallOpts{}, uint64(nonce), uint8(destId))
//		if err != nil {
//			return nil, err
//		}
//		return &relayer.Message{
//			Source:       sourceID,
//			Destination:  destId,
//			DepositNonce: nonce,
//			ResourceId:   record.ResourceID,
//			Type:         relayer.NonFungibleTransfer,
//			Payload: []interface{}{
//				record.TokenID.Bytes(),
//				record.DestinationRecipientAddress,
//				record.MetaData,
//			},
//		}, nil
//	}
//}
//
//func (c *CeloClient) ReturnGenericHandlerFabric() listener.EventHandler {
//	return func(sourceID, destId uint8, nonce uint64, handlerContractAddress string) (*relayer.Message, error) {
//		contract, err := genericHandler.NewGenericHandlerCaller(common.HexToAddress(handlerContractAddress), c.Client)
//		if err != nil {
//			return nil, err
//		}
//		record, err := contract.GetDepositRecord(&bind.CallOpts{}, uint64(nonce), uint8(destId))
//		if err != nil {
//			return nil, err
//		}
//		return &relayer.Message{
//			Source:       sourceID,
//			Destination:  destId,
//			DepositNonce: nonce,
//			ResourceId:   record.ResourceID,
//			Type:         relayer.GenericTransfer,
//			Payload: []interface{}{
//				record.MetaData,
//			},
//		}, nil
//	}
//}
