// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package client

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ChainSafe/chainbridge-celo-module/bindings/mptp/Bridge"
	"github.com/ChainSafe/chainbridge-core/chains/evm"
	coreListener "github.com/ChainSafe/chainbridge-core/chains/evm/listener"
	"github.com/ChainSafe/chainbridge-core/crypto/secp256k1"
	"github.com/ChainSafe/chainbridge-core/relayer"
	ethereum "github.com/celo-org/celo-blockchain"
	"github.com/celo-org/celo-blockchain/accounts/abi/bind"
	"github.com/celo-org/celo-blockchain/common"
	"github.com/celo-org/celo-blockchain/crypto"
	"github.com/celo-org/celo-blockchain/ethclient"
	"github.com/celo-org/celo-blockchain/rpc"
	"github.com/rs/zerolog/log"
)

var ErrFatalTx = errors.New("submission of transaction failed")
var ErrNonceTooLow = errors.New("nonce too low")
var ErrTxUnderpriced = errors.New("replacement transaction underpriced")

// Time between retrying a failed tx
const TxRetryInterval = time.Second * 2

// Tries to retry sending transaction
const TxRetryLimit = 10

const DefaultGasLimit = 6721975
const DefaultGasPrice = 20000000000

func NewCeloClient(endpoint string, http bool, sender *secp256k1.Keypair) (*CeloClient, error) {
	c := &CeloClient{
		endpoint: endpoint,
		http:     http,
		sender:   sender,
	}
	if err := c.connect(); err != nil {
		return nil, err
	}
	return c, nil
}

type CeloClient struct {
	*ethclient.Client
	endpoint      string
	http          bool
	stop          <-chan struct{}
	errChn        chan<- error
	optsLock      sync.Mutex
	opts          *bind.TransactOpts
	sender        *secp256k1.Keypair
	maxGasPrice   *big.Int   // TODO
	gasMultiplier *big.Float // TODO
	gasLimit      *big.Int
}

// Connect starts the ethereum WS connection
func (c *CeloClient) connect() error {
	log.Info().Str("url", c.endpoint).Msg("Connecting to ethereum chain...")
	var rpcClient *rpc.Client
	var err error
	// Start http or ws client
	if c.http {
		rpcClient, err = rpc.DialHTTP(c.endpoint)
	} else {
		rpcClient, err = rpc.DialWebsocket(context.Background(), c.endpoint, "/ws")
	}
	if err != nil {
		return err
	}
	c.Client = ethclient.NewClient(rpcClient)
	// TODO: move to config
	opts, err := c.newTransactOpts(big.NewInt(0), big.NewInt(DefaultGasLimit), big.NewInt(DefaultGasPrice))
	if err != nil {
		return err
	}
	c.opts = opts
	return nil
}

// LatestBlock returns the latest block from the current chain
func (c *CeloClient) LatestBlock() (*big.Int, error) {
	header, err := c.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	return header.Number, nil
}

func (c *CeloClient) GetEthClient() *ethclient.Client {
	return c.Client
}

func (c *CeloClient) ExecuteProposal(bridgeAddress string, proposal *evm.Proposal) error {
	for i := 0; i < TxRetryLimit; i++ {
		err := c.lockAndUpdateOpts()
		if err != nil {
			log.Error().Err(err).Msgf("failed to update tx opts")
			time.Sleep(TxRetryInterval)
		}
		b, err := Bridge.NewBridge(common.HexToAddress(bridgeAddress), c)
		if err != nil {
			return err
		}
		// Empty MPT verification params
		tx, err := b.ExecuteProposal(
			c.getOpts(),
			uint8(proposal.Source),
			uint64(proposal.DepositNonce),
			proposal.Data,
			proposal.ResourceId,
			nil,
			nil,
			[32]byte{},
			[32]byte{},
			nil,
			nil,
			//proposal.SVParams.Signature,
			//proposal.SVParams.AggregatePublicKey,
			//proposal.SVParams.BlockHash,
			//proposal.MPParams.TxRootHash,
			//proposal.MPParams.Key,
			//proposal.MPParams.Nodes,
		)
		c.unlockOpts()
		if err == nil {
			log.Info().Interface("source", proposal.Source).Interface("dest", proposal.Destination).Interface("nonce", proposal.DepositNonce).Str("tx", tx.Hash().Hex()).Msg("Submitted proposal execution")
			return nil
		}
		if err.Error() == ErrNonceTooLow.Error() || err.Error() == ErrTxUnderpriced.Error() {
			log.Error().Err(err).Msg("Nonce too low, will retry")
			time.Sleep(TxRetryInterval)
		} else {
			// TODO: this part is unclear. Does sending transaction with contract binding response with error if transaction failed inside contract?
			log.Error().Err(err).Msg("Execution failed, proposal may already be complete")
			time.Sleep(TxRetryInterval)
		}
		// Checking proposal status one more time (Since it could be execute by some other bridge). If it is completed then we do not need to retry
		s, err := c.ProposalStatus(bridgeAddress, proposal)
		if err != nil {
			log.Error().Err(err).Msgf("error getting proposal status %+v", proposal)
			continue
		}
		if s == relayer.ProposalStatusPassed || s == relayer.ProposalStatusExecuted || s == relayer.ProposalStatusCanceled {
			log.Info().Interface("source", proposal.Source).Interface("dest", proposal.Destination).Interface("nonce", proposal.DepositNonce).Msg("Proposal finalized on chain")
			return nil
		}
	}
	log.Error().Msgf("Submission of Execution transaction failed, source %v dest %v depNonce %v", proposal.Source, proposal.Destination, proposal.DepositNonce)
	return ErrFatalTx
}

func (c *CeloClient) VoteProposal(bridgeAddress string, proposal *evm.Proposal) error {
	for i := 0; i < TxRetryLimit; i++ {
		err := c.lockAndUpdateOpts()
		if err != nil {
			log.Error().Err(err).Msgf("failed to update tx opts")
		}
		b, err := Bridge.NewBridge(common.HexToAddress(bridgeAddress), c)
		if err != nil {
			return err
		}
		tx, err := b.VoteProposal(
			c.getOpts(),
			uint8(proposal.Source),
			uint64(proposal.DepositNonce),
			proposal.ResourceId,
			proposal.DataHash,
		)
		c.unlockOpts()
		if err == nil {
			log.Info().Interface("source", proposal.Source).Interface("dest", proposal.Destination).Interface("nonce", proposal.DepositNonce).Str("tx", tx.Hash().Hex()).Msg("Submitted proposal vote")
			return nil
		}
		if err.Error() == ErrNonceTooLow.Error() || err.Error() == ErrTxUnderpriced.Error() {
			log.Error().Err(err).Msg("Nonce too low, will retry")
			time.Sleep(TxRetryInterval)
		} else {
			// TODO: this part is unclear. Does sending transaction with contract binding response with error if transaction failed inside contract?
			log.Error().Err(err).Msg("Execution failed, proposal may already be complete")
			time.Sleep(TxRetryInterval)
		}
		// Checking proposal status one more time (Since it could be execute by some other bridge). If it is completed then we do not need to retry
		ps, err := c.ProposalStatus(bridgeAddress, proposal)
		if err != nil {
			log.Error().Err(err).Msgf("error getting proposal status %+v", proposal)
			continue
		}
		if ps == relayer.ProposalStatusPassed {
			log.Info().Interface("source", proposal.Source).Interface("dest", proposal.Destination).Interface("nonce", proposal.DepositNonce).Msg("Proposal is ready to be executed on chain")
			return nil
		}
	}
	log.Error().Msgf("Submission of vote transaction failed, source %v dest %v depNonce %v", proposal.Source, proposal.Destination, proposal.DepositNonce)
	return ErrFatalTx
}

func (c *CeloClient) ProposalStatus(bridgeAddress string, p *evm.Proposal) (relayer.ProposalStatus, error) {
	b, err := Bridge.NewBridge(common.HexToAddress(bridgeAddress), c)
	if err != nil {
		return 99, err
	}
	prop, err := b.GetProposal(&bind.CallOpts{}, p.Source, p.DepositNonce, p.DataHash)
	if err != nil {
		log.Error().Err(err).Msg("Failed to check proposal existence")
		return 99, err
	}
	log.Debug().Msgf("Fetching proposal %+v", prop)
	return relayer.ProposalStatus(prop.Status), nil
}

func (c *CeloClient) VotedBy(bridgeAddress string, p *evm.Proposal) bool {
	b, err := Bridge.NewBridge(common.HexToAddress(bridgeAddress), c)
	if err != nil {
		return false
	}
	addr := common.Address(c.sender.CommonAddress())
	hv, err := b.HasVotedOnProposal(&bind.CallOpts{}, evm.GetIDAndNonce(p), p.DataHash, addr)
	if err != nil {
		return false
	}
	return hv
}

func (c *CeloClient) MatchResourceIDToHandlerAddress(bridgeAddress string, rID [32]byte) (string, error) {
	b, err := Bridge.NewBridge(common.HexToAddress(bridgeAddress), c)
	if err != nil {
		return "", err
	}
	addr, err := b.ResourceIDToHandlerAddress(&bind.CallOpts{}, rID)
	if err != nil {
		return "", fmt.Errorf("failed to get handler from resource ID %x, reason: %w", rID, err)
	}
	return addr.String(), nil
}

// newTransactOpts builds the TransactOpts for the connection's keypair.
func (c *CeloClient) newTransactOpts(value, gasLimit, gasPrice *big.Int) (*bind.TransactOpts, error) {
	privateKey := c.sender.PrivateKey()
	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	nonce, err := c.PendingNonceAt(context.Background(), address)
	if err != nil {
		return nil, err
	}

	auth := bind.NewKeyedTransactor(privateKey)
	if err != nil {
		return nil, err
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = value
	auth.GasLimit = uint64(gasLimit.Int64())
	auth.GasPrice = gasPrice
	auth.Context = context.Background()
	return auth, nil
}

func (c *CeloClient) unlockOpts() {
	c.optsLock.Unlock()
}

func (c *CeloClient) lockAndUpdateOpts() error {
	c.optsLock.Lock()

	gasPrice, err := c.safeEstimateGas(context.TODO())
	if err != nil {
		return err
	}
	c.opts.GasPrice = gasPrice

	nonce, err := c.PendingNonceAt(context.Background(), c.opts.From)
	if err != nil {
		c.optsLock.Unlock()
		return err
	}
	c.opts.Nonce.SetUint64(nonce)
	return nil
}

func (c *CeloClient) getOpts() *bind.TransactOpts {
	return c.opts
}

func (c *CeloClient) safeEstimateGas(ctx context.Context) (*big.Int, error) {
	suggestedGasPrice, err := c.SuggestGasPrice(context.TODO())
	if err != nil {
		return nil, err
	}

	//gasPrice := multiplyGasPrice(suggestedGasPrice, c.gasMultiplier)

	// Check we aren't exceeding our limit
	//if suggestedGasPrice.Cmp(c.maxGasPrice) == 1 {
	//	return c.maxGasPrice, nil
	//} else {
	return suggestedGasPrice, nil
	//}
}

func multiplyGasPrice(gasEstimate *big.Int, gasMultiplier *big.Float) *big.Int {
	gasEstimateFloat := new(big.Float).SetInt(gasEstimate)
	result := gasEstimateFloat.Mul(gasEstimateFloat, gasMultiplier)
	gasPrice := new(big.Int)
	result.Int(gasPrice)
	return gasPrice
}

func (c *CeloClient) FetchDepositLogs(ctx context.Context, contractAddress string, sig string, startBlock *big.Int, endBlock *big.Int) ([]*coreListener.DepositLogs, error) {
	logs, err := c.FilterLogs(ctx, buildQuery(common.HexToAddress(contractAddress), sig, startBlock, endBlock))
	if err != nil {
		return nil, err
	}
	depositLogs := make([]*coreListener.DepositLogs, 0)

	for _, l := range logs {
		dl := &coreListener.DepositLogs{
			DestinationID: uint8(l.Topics[1].Big().Uint64()),
			ResourceID:    l.Topics[2],
			DepositNonce:  l.Topics[3].Big().Uint64(),
		}
		depositLogs = append(depositLogs, dl)
	}
	return depositLogs, nil
}

// buildQuery constructs a query for the bridgeContract by hashing sig to get the event topic
func buildQuery(contract common.Address, sig string, startBlock *big.Int, endBlock *big.Int) ethereum.FilterQuery {
	query := ethereum.FilterQuery{
		FromBlock: startBlock,
		ToBlock:   endBlock,
		Addresses: []common.Address{contract},
		Topics: [][]common.Hash{
			{crypto.Keccak256Hash([]byte(sig))},
		},
	}
	return query
}
