package client

//
//import (
//	"bytes"
//	"context"
//	"math/big"
//
//	"github.com/ethereum/go-ethereum/crypto"
//
//	"github.com/ChainSafe/chainbridge-core/chains/evm"
//	"github.com/ChainSafe/chainbridge-core/relayer"
//	"github.com/ethereum/go-ethereum/common"
//	"github.com/ethereum/go-ethereum/core/types"
//	"golang.org/x/crypto/sha3"
//)
//
//const DefaultGasLimit = 6721975
//const DefaultGasPrice = 20000000000
//
//type CeloBridge struct {
//	Address common.Address
//}
//
//func (r *CeloBridge) ExecuteProposal(bridgeAddress string, proposal *evm.Proposal) error {
//	//executeProposal(uint8 chainID, uint64 depositNonce, bytes data, bytes32 resourceID, bytes signatureHeader, bytes aggregatePublicKey, bytes32 hashedMessage, bytes32 rootHash, bytes key, bytes nodes)
//	data, err := buildDataUnsafe(
//		[]byte("executeProposal(uint8,uint64,bytes,bytes32,bytes,bytes,bytes32,bytes32,bytes,bytes)"),
//		big.NewInt(int64(proposal.Source)).Bytes(),
//		big.NewInt(int64(proposal.DepositNonce)).Bytes(),
//		proposal.Data,
//		proposal.ResourceId[:],
//		[]byte{},
//		[]byte{},
//		[]byte{},
//		[]byte{},
//		[]byte{},
//		[]byte{})
//	if err != nil {
//		return err
//	}
//
//	err = r.SignAndSendTransaction(data)
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//func (r *CeloBridge) VoteProposal(bridgeAddress string, proposal *evm.Proposal) error {
//	//voteProposal(uint8 chainID, uint64 depositNonce, bytes32 resourceID, bytes32 dataHash)
//
//	dataHash := createProposalDataHash(data, handlerContract, m.MPParams, m.SVParams)
//
//	data, err := buildDataUnsafe(
//		[]byte("voteProposal(uint8,uint64,bytes,bytes32,bytes32)"),
//		big.NewInt(int64(proposal.Source)).Bytes(),
//		big.NewInt(int64(proposal.DepositNonce)).Bytes(),
//		proposal.ResourceId[:],
//		proposal.DataHash,
//	)
//
//	err = r.SignAndSendTransaction(data)
//	if err != nil {
//		return err
//	}
//}
//
//func (r *CeloBridge) SignAndSendTransaction(data []byte) error {
//	nonce := uint64(0)
//	tx := NewCeloTransaction(nonce, r.Address, big.NewInt(0), DefaultGasLimit, big.NewInt(DefaultGasPrice), nil, nil, nil, data)
//
//	signedTx, err := types.SignTx(tx, types.HomesteadSigner{})
//
//	err = sender.SendGeneralizedTransaction(context.TODO(), tx)
//	if err != nil {
//		return err
//	}
//
//}
//
//func (r *CeloBridge) MatchResourceIDToHandlerAddress(bridgeAddress string, rID [32]byte) (string, error) {
//	return "", nil
//}
//
//func (r *CeloBridge) ProposalStatus(bridgeAddress string, proposal *evm.Proposal) (relayer.ProposalStatus, error) {
//	return 0, nil
//}
//
//func (r *CeloBridge) VotedBy(bridgeAddress string, p *evm.Proposal) bool {
//	return false
//}
//
//// CreateProposalDataHash constructs and returns proposal data hash
//// https://github.com/ChainSafe/chainbridge-celo-solidity/blob/1fae9c66a07139c277b03a09877414024867a8d9/contracts/Bridge.sol#L452-L454
//func createProposalDataHash(data []byte, handler common.Address, mp *utils.MerkleProof, sv *utils.SignatureVerification) common.Hash {
//	b := bytes.NewBuffer(handler.Bytes())
//	b.Write(data)
//	b.Write(mp.TxRootHash[:])
//	b.Write(mp.Key)
//	b.Write(mp.Nodes)
//	b.Write(sv.AggregatePublicKey)
//	b.Write(sv.BlockHash[:])
//	b.Write(sv.Signature)
//	return crypto.Keccak256Hash(b.Bytes())
//}
//
//// buildDataUnsafe is unsafe function that collects encoded hex params into byte array
//func buildDataUnsafe(method []byte, params ...[]byte) ([]byte, error) {
//	hash := sha3.NewLegacyKeccak256()
//	_, err := hash.Write(method)
//	if err != nil {
//		return nil, err
//	}
//	methodID := hash.Sum(nil)[:4]
//
//	var data []byte
//	data = append(data, methodID...)
//	for _, v := range params {
//		paddedParam := common.LeftPadBytes(v, 32)
//		data = append(data, paddedParam...)
//	}
//	return data, nil
//}

//
//func (c *CeloClient) ExecuteProposal(bridgeAddress string, proposal *evm.Proposal) error {
//	for i := 0; i < TxRetryLimit; i++ {
//		err := c.lockAndUpdateOpts()
//		if err != nil {
//			log.Error().Err(err).Msgf("failed to update tx opts")
//			time.Sleep(TxRetryInterval)
//		}
//		b, err := Bridge.NewBridge(common.HexToAddress(bridgeAddress), c)
//		if err != nil {
//			return err
//		}
//		// Empty MPT verification params
//		tx, err := b.ExecuteProposal(
//			c.getOpts(),
//			uint8(proposal.Source),
//			uint64(proposal.DepositNonce),
//			proposal.Data,
//			proposal.ResourceId,
//			nil,
//			nil,
//			[32]byte{},
//			[32]byte{},
//			nil,
//			nil,
//			//proposal.SVParams.Signature,
//			//proposal.SVParams.AggregatePublicKey,
//			//proposal.SVParams.BlockHash,
//			//proposal.MPParams.TxRootHash,
//			//proposal.MPParams.Key,
//			//proposal.MPParams.Nodes,
//		)
//		c.unlockOpts()
//		if err == nil {
//			log.Info().Interface("source", proposal.Source).Interface("dest", proposal.Destination).Interface("nonce", proposal.DepositNonce).Str("tx", tx.Hash().Hex()).Msg("Submitted proposal execution")
//			return nil
//		}
//		if err.Error() == ErrNonceTooLow.Error() || err.Error() == ErrTxUnderpriced.Error() {
//			log.Error().Err(err).Msg("Nonce too low, will retry")
//			time.Sleep(TxRetryInterval)
//		} else {
//			// TODO: this part is unclear. Does sending transaction with contract binding response with error if transaction failed inside contract?
//			log.Error().Err(err).Msg("Execution failed, proposal may already be complete")
//			time.Sleep(TxRetryInterval)
//		}
//		// Checking proposal status one more time (Since it could be execute by some other bridge). If it is completed then we do not need to retry
//		s, err := c.ProposalStatus(bridgeAddress, proposal)
//		if err != nil {
//			log.Error().Err(err).Msgf("error getting proposal status %+v", proposal)
//			continue
//		}
//		if s == relayer.ProposalStatusPassed || s == relayer.ProposalStatusExecuted || s == relayer.ProposalStatusCanceled {
//			log.Info().Interface("source", proposal.Source).Interface("dest", proposal.Destination).Interface("nonce", proposal.DepositNonce).Msg("Proposal finalized on chain")
//			return nil
//		}
//	}
//	log.Error().Msgf("Submission of Execution transaction failed, source %v dest %v depNonce %v", proposal.Source, proposal.Destination, proposal.DepositNonce)
//	return ErrFatalTx
//}
//
//func (c *CeloClient) VoteProposal(bridgeAddress string, proposal *evm.Proposal) error {
//	for i := 0; i < TxRetryLimit; i++ {
//		err := c.lockAndUpdateOpts()
//		if err != nil {
//			log.Error().Err(err).Msgf("failed to update tx opts")
//		}
//		b, err := Bridge.NewBridge(common.HexToAddress(bridgeAddress), c)
//		if err != nil {
//			return err
//		}
//		tx, err := b.VoteProposal(
//			c.getOpts(),
//			uint8(proposal.Source),
//			uint64(proposal.DepositNonce),
//			proposal.ResourceId,
//			proposal.DataHash,
//		)
//		c.unlockOpts()
//		if err == nil {
//			log.Info().Interface("source", proposal.Source).Interface("dest", proposal.Destination).Interface("nonce", proposal.DepositNonce).Str("tx", tx.Hash().Hex()).Msg("Submitted proposal vote")
//			return nil
//		}
//		if err.Error() == ErrNonceTooLow.Error() || err.Error() == ErrTxUnderpriced.Error() {
//			log.Error().Err(err).Msg("Nonce too low, will retry")
//			time.Sleep(TxRetryInterval)
//		} else {
//			// TODO: this part is unclear. Does sending transaction with contract binding response with error if transaction failed inside contract?
//			log.Error().Err(err).Msg("Execution failed, proposal may already be complete")
//			time.Sleep(TxRetryInterval)
//		}
//		// Checking proposal status one more time (Since it could be execute by some other bridge). If it is completed then we do not need to retry
//		ps, err := c.ProposalStatus(bridgeAddress, proposal)
//		if err != nil {
//			log.Error().Err(err).Msgf("error getting proposal status %+v", proposal)
//			continue
//		}
//		if ps == relayer.ProposalStatusPassed {
//			log.Info().Interface("source", proposal.Source).Interface("dest", proposal.Destination).Interface("nonce", proposal.DepositNonce).Msg("Proposal is ready to be executed on chain")
//			return nil
//		}
//	}
//	log.Error().Msgf("Submission of vote transaction failed, source %v dest %v depNonce %v", proposal.Source, proposal.Destination, proposal.DepositNonce)
//	return ErrFatalTx
//}
//
//func (c *CeloClient) ProposalStatus(bridgeAddress string, p *evm.Proposal) (relayer.ProposalStatus, error) {
//	b, err := Bridge.NewBridge(common.HexToAddress(bridgeAddress), c)
//	if err != nil {
//		return 99, err
//	}
//	prop, err := b.GetProposal(&bind.CallOpts{}, p.Source, p.DepositNonce, p.DataHash)
//	if err != nil {
//		log.Error().Err(err).Msg("Failed to check proposal existence")
//		return 99, err
//	}
//	log.Debug().Msgf("Fetching proposal %+v", prop)
//	return relayer.ProposalStatus(prop.Status), nil
//}
//
//func (c *CeloClient) VotedBy(bridgeAddress string, p *evm.Proposal) bool {
//	b, err := Bridge.NewBridge(common.HexToAddress(bridgeAddress), c)
//	if err != nil {
//		return false
//	}
//	addr := common.HexToAddress(c.sender.Address())
//	hv, err := b.HasVotedOnProposal(&bind.CallOpts{}, evm.GetIDAndNonce(p), p.DataHash, addr)
//	if err != nil {
//		return false
//	}
//	return hv
//}
//
//func (c *CeloClient) MatchResourceIDToHandlerAddress(bridgeAddress string, rID [32]byte) (string, error) {
//	b, err := Bridge.NewBridge(common.HexToAddress(bridgeAddress), c)
//	if err != nil {
//		return "", err
//	}
//	addr, err := b.ResourceIDToHandlerAddress(&bind.CallOpts{}, rID)
//	if err != nil {
//		return "", fmt.Errorf("failed to get handler from resource ID %x, reason: %w", rID, err)
//	}
//	return addr.String(), nil
//}
//
//func (c *CeloClient) MatchResourceIDToHandlerAddress1(bridgeAddress string, rID [32]byte) (string, error) {
//	b, err := Bridge.NewBridge(common.HexToAddress(bridgeAddress), c)
//	if err != nil {
//		return "", err
//	}
//	addr, err := b.ResourceIDToHandlerAddress(&bind.CallOpts{}, rID)
//	if err != nil {
//		return "", fmt.Errorf("failed to get handler from resource ID %x, reason: %w", rID, err)
//	}
//	return addr.String(), nil
//}
