package transaction

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
)

func NewCeloTransaction(nonce uint64, to common.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, feeCurrency, gatewayFeeRecipient *common.Address, gatewayFee *big.Int, data []byte) *CeloTransaction {
	return newTransaction(nonce, &to, amount, gasLimit, gasPrice, feeCurrency, gatewayFeeRecipient, gatewayFee, data)
}

type CeloTransaction struct {
	data txdata
	// caches
	hash atomic.Value
	size atomic.Value
	from atomic.Value
}

func (a *CeloTransaction) RawWithSignature(key *ecdsa.PrivateKey, chainID *big.Int) ([]byte, error) {
	opts := NewKeyedTransactor(key)

	signedTx, err := opts.Signer(CeloSigner{}, crypto.PubkeyToAddress(key.PublicKey), a)
	if err != nil {
		return nil, err
	}
	rawTX, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		return nil, err
	}
	return rawTX, nil
}

func (tx *CeloTransaction) Hash() common.Hash {
	if hash := tx.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	v := rlpHash(tx)
	tx.hash.Store(v)
	return v
}

func (tx *CeloTransaction) WithSignature(signer CeloSigner, sig []byte) (*CeloTransaction, error) {
	r, s, v, err := signer.SignatureValues(tx, sig)
	if err != nil {
		return nil, err
	}
	cpy := &CeloTransaction{data: tx.data}
	cpy.data.R, cpy.data.S, cpy.data.V = r, s, v
	return cpy, nil
}

// MarshalBinary returns the canonical encoding of the transaction.
// For legacy transactions, it returns the RLP encoding. For EIP-2718 typed
// transactions, it returns the type and payload.
func (tx *CeloTransaction) MarshalBinary() ([]byte, error) {
	return rlp.EncodeToBytes(tx)
}

func newTransaction(nonce uint64, to *common.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, feeCurrency, gatewayFeeRecipient *common.Address, gatewayFee *big.Int, data []byte) *CeloTransaction {
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
	if gasPrice != nil {
		d.Price.Set(gasPrice)
	}

	return &CeloTransaction{data: d}
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

type SignerFn func(CeloSigner, common.Address, *CeloTransaction) (*CeloTransaction, error)

type TransactOpts struct {
	From   common.Address // Ethereum account to send the transaction from
	Nonce  *big.Int       // Nonce to use for the transaction execution (nil = use pending state)
	Signer SignerFn       // Method to use for signing the transaction (mandatory)

	Value               *big.Int        // Funds to transfer along along the transaction (nil = 0 = no funds)
	GasPrice            *big.Int        // Gas price to use for the transaction execution (nil = gas price oracle)
	FeeCurrency         *common.Address // Fee currency to be used for transaction (nil = default currency = Celo Gold)
	GatewayFeeRecipient *common.Address // Address to which gateway fees should be paid (nil = no gateway fees are paid)
	GatewayFee          *big.Int        // Value of gateway fees to be paid (nil = no gateway fees are paid)
	GasLimit            uint64          // Gas limit to set for the transaction execution (0 = estimate)

	Context context.Context // Network context to support cancellation and timeouts (nil = no timeout)
}

func rlpHash(x interface{}) (h common.Hash) {
	hw := sha3.NewLegacyKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}

func NewKeyedTransactor(key *ecdsa.PrivateKey) *TransactOpts {
	keyAddr := crypto.PubkeyToAddress(key.PublicKey)
	return &TransactOpts{
		From: keyAddr,
		Signer: func(signer CeloSigner, address common.Address, tx *CeloTransaction) (*CeloTransaction, error) {
			if address != keyAddr {
				return nil, errors.New("not authorized to sign this account")
			}
			signature, err := crypto.Sign(signer.Hash(tx).Bytes(), key)
			if err != nil {
				return nil, err
			}
			return tx.WithSignature(signer, signature)
		},
	}
}
