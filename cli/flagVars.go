package cli

import (
	"math/big"

	"github.com/ChainSafe/chainbridge-core/crypto/secp256k1"
)

// global flags
var (
	url           string
	gasLimit      uint64
	gasPrice      *big.Int
	prepare       bool
	senderKeyPair *secp256k1.Keypair
)
