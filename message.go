package client

import (
	"github.com/ChainSafe/chainbridge-core/relayer"
	"github.com/ethereum/go-ethereum/common"
)

type Message struct {
	//Source       uint8  // Source where message was initiated
	//Destination  uint8  // Destination chain of message
	//DepositNonce uint64 // Nonce for the deposit
	//ResourceId   [32]byte
	//Type         relayer.TransferType
	relayer.Message
	MPParams *MerkleProof
	SVParams *SignatureVerification
	//Payload      []interface{} // data associated with event sequence
}

type MerkleProof struct {
	TxRootHash [32]byte // Expected root of trie, in our case should be transactionsRoot from block
	Key        []byte   // RLP encoding of tx index, for the tx we want to prove
	Nodes      []byte   // The actual proof, all the nodes of the trie that between leaf value and root
}

type SignatureVerification struct {
	AggregatePublicKey []byte      // Aggregated public key of block validators
	BlockHash          common.Hash // Hash of block we are proving
	Signature          []byte      // Signature of block we are proving
	RLPHeader          []byte      // RLP encoding of header data
}

func sliceTo32Bytes(in []byte) [32]byte {
	var res [32]byte
	copy(res[:], in)
	return res
}
