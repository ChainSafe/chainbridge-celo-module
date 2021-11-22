package transaction

import (
	"crypto/ecdsa"
	"errors"
	"math/big"
	"sync/atomic"

	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
)

var (
	ErrInvalidSig = errors.New("invalid transaction v, r, s values")
	// ErrEthCompatibleTransactionIsntCompatible is returned if the transaction has EthCompatible: true
	// but has non-nil-or-0 values for some of the Celo-only fields
	ErrEthCompatibleTransactionIsntCompatible = errors.New("ethCompatible is true, but non-eth-compatible fields are present")
)

func NewCeloTransaction(nonce uint64, to *common.Address, amount *big.Int, gasLimit uint64, gasPrice []*big.Int, data []byte) (evmclient.CommonTransaction, error) {
	return newTransaction(nonce, to, amount, gasLimit, gasPrice, nil, nil, nil, data), nil
}

type CeloTransaction struct {
	data txdata
	// caches
	hash atomic.Value
	from atomic.Value
}

type txdata struct {
	AccountNonce        uint64          `json:"nonce"    gencodec:"required"`
	Price               *big.Int        `json:"gasPrice" gencodec:"required"`
	GasLimit            uint64          `json:"gas"      gencodec:"required"`
	FeeCurrency         *common.Address `json:"feeCurrency" rlp:"nil"`         // nil means native currency
	GatewayFeeRecipient *common.Address `json:"gatewayFeeRecipient" rlp:"nil"` // nil means no gateway fee is paid
	GatewayFee          *big.Int        `json:"gatewayFee"`
	Recipient           *common.Address `json:"to"       rlp:"nil"` // nil means contract creation
	Amount              *big.Int        `json:"value"    gencodec:"required"`
	Payload             []byte          `json:"input"    gencodec:"required"`

	// Signature values
	V *big.Int `json:"v" gencodec:"required"`
	R *big.Int `json:"r" gencodec:"required"`
	S *big.Int `json:"s" gencodec:"required"`

	// This is only used when marshaling to JSON.
	Hash *common.Hash `json:"hash" rlp:"-"`

	// Whether this is an ethereum-compatible transaction (i.e. with FeeCurrency, GatewayFeeRecipient and GatewayFee omitted)
	EthCompatible bool `json:"ethCompatible" rlp:"-"`
}

func newTransaction(nonce uint64, to *common.Address, amount *big.Int, gasLimit uint64, gasPrice []*big.Int, feeCurrency, gatewayFeeRecipient *common.Address, gatewayFee *big.Int, data []byte) *CeloTransaction {
	if len(data) > 0 {
		data = common.CopyBytes(data)
	}
	d := txdata{
		AccountNonce:        nonce,
		Recipient:           to,
		Payload:             data,
		Amount:              new(big.Int),
		GasLimit:            gasLimit,
		FeeCurrency:         feeCurrency,
		GatewayFeeRecipient: gatewayFeeRecipient,
		GatewayFee:          new(big.Int),
		Price:               new(big.Int),
		V:                   new(big.Int),
		R:                   new(big.Int),
		S:                   new(big.Int),
	}
	if amount != nil {
		d.Amount.Set(amount)
	}
	if gatewayFee != nil {
		d.GatewayFee.Set(gatewayFee)
	}
	if len(gasPrice) != 0 {
		d.Price.Set(gasPrice[0])
	}

	return &CeloTransaction{data: d}
}

// ChainId returns which chain id this transaction was signed for (if at all)
func (tx *CeloTransaction) ChainId() *big.Int {
	return deriveChainId(tx.data.V)
}

// Protected returns whether the transaction is protected from replay protection.
func (tx *CeloTransaction) Protected() bool {
	return isProtectedV(tx.data.V)
}

func isProtectedV(V *big.Int) bool {
	if V.BitLen() <= 8 {
		v := V.Uint64()
		return v != 27 && v != 28
	}
	// anything not 27 or 28 is considered protected
	return true
}

// To returns the recipient address of the transaction.
// It returns nil if the transaction is a contract creation.
func (tx *CeloTransaction) To() *common.Address {
	if tx.data.Recipient == nil {
		return nil
	}
	to := *tx.data.Recipient
	return &to
}

// Hash hashes the RLP encoding of tx.
// It uniquely identifies the transaction.
func (tx *CeloTransaction) Hash() common.Hash {
	if hash := tx.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	v := rlpHash(tx)
	tx.hash.Store(v)
	return v
}

// WithSignature returns a new transaction with the given signature.
// This signature needs to be in the [R || S || V] format where V is 0 or 1.
func (tx *CeloTransaction) WithSignature(signer CeloSigner, sig []byte) (*CeloTransaction, error) {
	r, s, v, err := signer.SignatureValues(tx, sig)
	if err != nil {
		return nil, err
	}
	cpy := &CeloTransaction{data: tx.data}
	cpy.data.R, cpy.data.S, cpy.data.V = r, s, v
	tx.data.R, tx.data.S, tx.data.V = r, s, v
	return cpy, nil
}

func (tx *CeloTransaction) RawWithSignature(key *ecdsa.PrivateKey, chainID *big.Int) ([]byte, error) {
	opts := NewKeyedTransactor(key)
	signedTx, err := opts.Signer(NewEIP155Signer(chainID), crypto.PubkeyToAddress(key.PublicKey), tx)
	if err != nil {
		return nil, err
	}
	rawTX, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		return nil, err
	}
	return rawTX, nil
}

func rlpHash(x interface{}) (h common.Hash) {
	hw := sha3.NewLegacyKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}
