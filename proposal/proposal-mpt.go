package proposal

import (
	"bytes"
	"context"
	"math/big"
	"strconv"
	"strings"

	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmtransaction"
	"github.com/ChainSafe/chainbridge-core/relayer"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
	"github.com/status-im/keycard-go/hexutils"
)

type Proposal struct {
	Source         uint8  // Source where message was initiated
	Destination    uint8  // Destination chain of message
	DepositNonce   uint64 // Nonce for the deposit
	ResourceId     [32]byte
	Payload        []interface{} // data associated with event sequence
	Data           []byte
	HandlerAddress common.Address
	BridgeAddress  common.Address
}

type CeloChainClient interface {
	SignAndSendTransaction(ctx context.Context, tx evmclient.CommonTransaction) (common.Hash, error)
	CallContract(ctx context.Context, callArgs map[string]interface{}, blockNumber *big.Int) ([]byte, error)
	UnsafeNonce() (*big.Int, error)
	LockNonce()
	UnlockNonce()
	UnsafeIncreaseNonce() error
	GasPrice() (*big.Int, error)
}

func (p *Proposal) Status(evmCaller CeloChainClient) (relayer.ProposalStatus, error) {
	definition := "[{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"originChainID\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"depositNonce\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"dataHash\",\"type\":\"bytes32\"}],\"name\":\"getProposal\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"_resourceID\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_dataHash\",\"type\":\"bytes32\"},{\"internalType\":\"address[]\",\"name\":\"_yesVotes\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_noVotes\",\"type\":\"address[]\"},{\"internalType\":\"enumBridge.ProposalStatus\",\"name\":\"_status\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"_proposedBlock\",\"type\":\"uint256\"}],\"internalType\":\"structBridge.Proposal\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true}]"
	a, err := abi.JSON(strings.NewReader(definition))
	if err != nil {
		return relayer.ProposalStatusInactive, err // Not sure what status to use here
	}
	input, err := a.Pack("getProposal", p.Source, p.DepositNonce, p.GetDataHash())
	if err != nil {
		return relayer.ProposalStatusInactive, err
	}
	log.Debug().Msg(hexutils.BytesToHex(input))

	msg := ethereum.CallMsg{From: common.Address{}, To: &p.BridgeAddress, Data: input}
	out, err := evmCaller.CallContract(context.TODO(), toCallArg(msg), nil)
	log.Debug().Msg(strconv.Itoa(len(out)))
	if err != nil {
		return relayer.ProposalStatusInactive, err
	}
	type bridgeProposal struct {
		ResourceID    [32]byte
		DataHash      [32]byte
		YesVotes      []common.Address
		NoVotes       []common.Address
		Status        uint8
		ProposedBlock *big.Int
	}
	res, err := a.Unpack("getProposal", out)
	out0 := *abi.ConvertType(res[0], new(bridgeProposal)).(*bridgeProposal)
	return relayer.ProposalStatus(out0.Status), nil
}

func (p *Proposal) VotedBy(evmCaller CeloChainClient, by common.Address) (bool, error) {
	definition := "[{\"inputs\":[{\"internalType\":\"uint72\",\"name\":\"\",\"type\":\"uint72\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"_hasVotedOnProposal\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"
	a, err := abi.JSON(strings.NewReader(definition))
	if err != nil {
		return false, err // Not sure what status to use here
	}
	input, err := a.Pack("_hasVotedOnProposal", idAndNonce(p.Source, p.DepositNonce), p.GetDataHash(), by)
	if err != nil {
		return false, err
	}
	msg := ethereum.CallMsg{From: common.Address{}, To: &p.BridgeAddress, Data: input}
	out, err := evmCaller.CallContract(context.TODO(), toCallArg(msg), nil)
	if err != nil {
		return false, err
	}
	res, err := a.Unpack("_hasVotedOnProposal", out)
	out0 := *abi.ConvertType(res[0], new(bool)).(*bool)
	return out0, nil
}

func (p *Proposal) Execute(client CeloChainClient) error {
	definition := "[{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"chainID\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"depositNonce\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"resourceID\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"signatureHeader\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"aggregatePublicKey\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"hashedMessage\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"rootHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"key\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"nodes\",\"type\":\"bytes\"}],\"name\":\"executeProposal\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\",\"constant\":false}]"
	a, err := abi.JSON(strings.NewReader(definition))
	if err != nil {
		return err // Not sure what status to use here
	}
	input, err := a.Pack("executeProposal", p.Source, p.DepositNonce, p.Data, p.ResourceId, []byte{}, []byte{}, []byte{}, []byte{}, []byte{}, []byte{})
	if err != nil {
		return err
	}
	gasLimit := uint64(2000000)
	gp, err := client.GasPrice()
	if err != nil {
		return err
	}
	client.LockNonce()
	n, err := client.UnsafeNonce()
	if err != nil {
		return err
	}
	tx := evmtransaction.NewTransaction(n.Uint64(), p.BridgeAddress, big.NewInt(0), gasLimit, gp, input)
	h, err := client.SignAndSendTransaction(context.TODO(), tx)
	if err != nil {
		return err
	}
	err = client.UnsafeIncreaseNonce()
	if err != nil {
		return err
	}
	client.UnlockNonce()
	log.Debug().Str("hash", h.Hex()).Msgf("Executed")
	return nil
}

func (p *Proposal) Vote(client CeloChainClient) error {
	definition := "[{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"chainID\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"depositNonce\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"resourceID\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"dataHash\",\"type\":\"bytes32\"}],\"name\":\"voteProposal\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"
	a, err := abi.JSON(strings.NewReader(definition))
	if err != nil {
		return err // Not sure what status to use here
	}
	input, err := a.Pack("voteProposal", p.Source, p.DepositNonce, p.ResourceId, p.GetDataHash())
	if err != nil {
		return err
	}
	gasLimit := uint64(1000000)
	gp, err := client.GasPrice()
	if err != nil {
		return err
	}
	client.LockNonce()
	n, err := client.UnsafeNonce()
	if err != nil {
		return err
	}
	tx := evmtransaction.NewTransaction(n.Uint64(), p.BridgeAddress, big.NewInt(0), gasLimit, gp, input)
	h, err := client.SignAndSendTransaction(context.TODO(), tx)
	if err != nil {
		return err
	}
	err = client.UnsafeIncreaseNonce()
	if err != nil {
		return err
	}
	client.UnlockNonce()
	log.Debug().Str("hash", h.Hex()).Msgf("Voted")
	return nil
}

// CreateProposalDataHash constructs and returns proposal data hash
func (p *Proposal) GetDataHash() common.Hash {
	b := bytes.NewBuffer(p.HandlerAddress.Bytes())
	b.Write(p.Data)
	//b.Write(mp.TxRootHash[:])
	//b.Write(mp.Key)
	//b.Write(mp.Nodes)
	//b.Write(sv.AggregatePublicKey)
	//b.Write(sv.BlockHash[:])
	//b.Write(sv.Signature)
	return crypto.Keccak256Hash(b.Bytes())
}

func idAndNonce(srcId uint8, nonce uint64) *big.Int {
	var data []byte
	data = append(data, big.NewInt(int64(nonce)).Bytes()...)
	data = append(data, uint8(srcId))
	return big.NewInt(0).SetBytes(data)
}

func toCallArg(msg ethereum.CallMsg) map[string]interface{} {
	arg := map[string]interface{}{
		"from": msg.From,
		"to":   msg.To,
	}
	if len(msg.Data) > 0 {
		arg["data"] = hexutil.Bytes(msg.Data)
	}
	if msg.Value != nil {
		arg["value"] = (*hexutil.Big)(msg.Value)
	}
	if msg.Gas != 0 {
		arg["gas"] = hexutil.Uint64(msg.Gas)
	}
	if msg.GasPrice != nil {
		arg["gasPrice"] = (*hexutil.Big)(msg.GasPrice)
	}
	return arg
}
