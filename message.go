package client

import (
	"github.com/celo-org/celo-blockchain/common"
)

type MerkleProof struct {
	TxRootHash [32]byte // Expected root of trie, in our case should be transactionsRoot from block
	Key        []byte   // RLP encoding of tx index, for the tx we want to prove
	Nodes      []byte   // The actual proof, all the nodes of the trie that between leaf value and root
}

type SignatureVerification struct {
	AggregatePublicKey []byte      // Aggregated public key of block validators
	BlockHash          common.Hash // Hash of block we are proving
	Signature          []byte      // Signature of block we are proving
}

func sliceTo32Bytes(in []byte) [32]byte {
	var res [32]byte
	copy(res[:], in)
	return res
}
