package voter

import (
	"errors"
	"math/big"

	"github.com/ChainSafe/chainbridge-celo-module/proposal"
	"github.com/ChainSafe/chainbridge-core/chains/evm/voter"
	"github.com/ChainSafe/chainbridge-core/relayer"
	"github.com/ethereum/go-ethereum/common"
)

func ERC20CeloMessageHandler(m *relayer.Message, handlerAddr, bridgeAddress common.Address) (voter.Proposer, error) {
	if len(m.Payload) != 2 {
		return nil, errors.New("malformed payload. Len  of payload should be 2")
	}
	amount, ok := m.Payload[0].([]byte)
	if !ok {
		return nil, errors.New("wrong payloads amount format")
	}

	recipient, ok := m.Payload[1].([]byte)
	if !ok {
		return nil, errors.New("wrong payloads recipient format")

	}
	var data []byte
	data = append(data, common.LeftPadBytes(amount, 32)...) // amount (uint256)

	recipientLen := big.NewInt(int64(len(recipient))).Bytes()
	data = append(data, common.LeftPadBytes(recipientLen, 32)...) // length of recipient (uint256)
	data = append(data, recipient...)                             // recipient ([]byte)

	return proposal.ProposalWithMPTVerification{
		Source:         m.Source,
		DepositNonce:   m.DepositNonce,
		ResourceId:     m.ResourceId,
		Data:           data,
		HandlerAddress: handlerAddr,
		BridgeAddress:  bridgeAddress,
	}, nil
}
