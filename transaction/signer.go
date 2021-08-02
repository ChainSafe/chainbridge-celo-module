package transaction

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/params"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	ErrInvalidChainId = errors.New("invalid chain id for signer")
)

// Signer encapsulates transaction signature handling. Note that this interface is not a
// stable API and may change at any time to accommodate new protocol rules.
type CeloSigner interface {
	// Sender returns the sender address of the transaction.
	Sender(tx *CeloTransaction) (common.Address, error)
	// SignatureValues returns the raw R, S, V values corresponding to the
	// given signature.
	SignatureValues(tx *CeloTransaction, sig []byte) (r, s, v *big.Int, err error)
	// Hash returns the hash to be signed.
	Hash(tx *CeloTransaction) common.Hash
	// Equal returns true if the given signer is the same as the receiver.
	Equal(CeloSigner) bool
}


// sigCache is used to cache the derived sender and contains
// the signer used to derive it.
type sigCache struct {
	signer CeloSigner
	from   common.Address
}

// MakeSigner returns a Signer based on the given chain config and block number.
func MakeSigner(config *params.ChainConfig, blockNumber *big.Int) CeloSigner {
	var signer CeloSigner
	switch {
	case config.IsEIP155(blockNumber):
		signer = NewEIP155Signer(config.ChainID)
	case config.IsHomestead(blockNumber):
		signer = HomesteadSigner{}
	default:
		signer = FrontierSigner{}
	}
	return signer
}

// SignTx signs the transaction using the given signer and private key
func SignTx(tx *CeloTransaction, s CeloSigner, prv *ecdsa.PrivateKey) (*CeloTransaction, error) {
	h := s.Hash(tx)
	sig, err := crypto.Sign(h[:], prv)
	if err != nil {
		return nil, err
	}
	return tx.WithSignature(s, sig)
}

// Sender returns the address derived from the signature (V, R, S) using secp256k1
// elliptic curve and an error if it failed deriving or upon an incorrect
// signature.
//
// Sender may cache the address, allowing it to be used regardless of
// signing method. The cache is invalidated if the cached signer does
// not match the signer used in the current call.
func Sender(signer CeloSigner, tx *CeloTransaction) (common.Address, error) {
	if sc := tx.from.Load(); sc != nil {
		sigCache := sc.(sigCache)
		// If the signer used to derive from in a previous
		// call is not the same as used current, invalidate
		// the cache.
		if sigCache.signer.Equal(signer) {
			return sigCache.from, nil
		}
	}

	addr, err := signer.Sender(tx)
	if err != nil {
		return common.Address{}, err
	}
	tx.from.Store(sigCache{signer: signer, from: addr})
	return addr, nil
}

// EIP155Transaction implements Signer using the EIP155 rules.
type EIP155Signer struct {
	chainId, chainIdMul *big.Int
}

func NewEIP155Signer(chainId *big.Int) EIP155Signer {
	if chainId == nil {
		chainId = new(big.Int)
	}
	return EIP155Signer{
		chainId:    chainId,
		chainIdMul: new(big.Int).Mul(chainId, big.NewInt(2)),
	}
}

func (s EIP155Signer) Equal(s2 CeloSigner) bool {
	eip155, ok := s2.(EIP155Signer)
	return ok && eip155.chainId.Cmp(s.chainId) == 0
}

var big8 = big.NewInt(8)

func (s EIP155Signer) Sender(tx *CeloTransaction) (common.Address, error) {
	if !tx.Protected() {
		return HomesteadSigner{}.Sender(tx)
	}
	if tx.ChainId().Cmp(s.chainId) != 0 {
		return common.Address{}, ErrInvalidChainId
	}
	V := new(big.Int).Sub(tx.data.V, s.chainIdMul)
	V.Sub(V, big8)
	addr, _, err := recoverPlain(s.Hash(tx), tx.data.R, tx.data.S, V, true)
	return addr, err
}

// SignatureValues returns signature values. This signature
// needs to be in the [R || S || V] format where V is 0 or 1.
func (s EIP155Signer) SignatureValues(tx *CeloTransaction, sig []byte) (R, S, V *big.Int, err error) {
	R, S, V, err = HomesteadSigner{}.SignatureValues(tx, sig)
	if err != nil {
		return nil, nil, nil, err
	}
	if s.chainId.Sign() != 0 {
		V = big.NewInt(int64(sig[64] + 35))
		V.Add(V, s.chainIdMul)
	}
	return R, S, V, nil
}

// Hash returns the hash to be signed by the sender.
// It does not uniquely identify the transaction.
func (s EIP155Signer) Hash(tx *CeloTransaction) common.Hash {
	if tx.data.EthCompatible {
		return rlpHash([]interface{}{
			tx.data.AccountNonce,
			tx.data.Price,
			tx.data.GasLimit,
			tx.data.Recipient,
			tx.data.Amount,
			tx.data.Payload,
			s.chainId, uint(0), uint(0),
		})
	} else {
		return rlpHash([]interface{}{
			tx.data.AccountNonce,
			tx.data.Price,
			tx.data.GasLimit,
			tx.data.FeeCurrency,
			tx.data.GatewayFeeRecipient,
			tx.data.GatewayFee,
			tx.data.Recipient,
			tx.data.Amount,
			tx.data.Payload,
			s.chainId, uint(0), uint(0),
		})
	}
}

// HomesteadTransaction implements TransactionInterface using the
// homestead rules.
type HomesteadSigner struct{ FrontierSigner }

func (s HomesteadSigner) Equal(s2 CeloSigner) bool {
	_, ok := s2.(HomesteadSigner)
	return ok
}

// SignatureValues returns signature values. This signature
// needs to be in the [R || S || V] format where V is 0 or 1.
func (hs HomesteadSigner) SignatureValues(tx *CeloTransaction, sig []byte) (r, s, v *big.Int, err error) {
	return hs.FrontierSigner.SignatureValues(tx, sig)
}

func (hs HomesteadSigner) Sender(tx *CeloTransaction) (common.Address, error) {
	addr, _, err := recoverPlain(hs.Hash(tx), tx.data.R, tx.data.S, tx.data.V, true)
	return addr, err
}

func (hs HomesteadSigner) SenderData(data common.Hash, sig []byte) (common.Address, []byte, error) {
	r, s, v, err := hs.SignatureValues(nil, sig)
	v = new(big.Int).Sub(v, big.NewInt(27))
	if err != nil {
		return common.Address{}, nil, err
	}
	return recoverPlain(data, r, s, v, true)
}

type FrontierSigner struct{}

func (s FrontierSigner) Equal(s2 CeloSigner) bool {
	_, ok := s2.(FrontierSigner)
	return ok
}

// SignatureValues returns signature values. This signature
// needs to be in the [R || S || V] format where V is 0 or 1.
func (fs FrontierSigner) SignatureValues(tx *CeloTransaction, sig []byte) (r, s, v *big.Int, err error) {
	if len(sig) != crypto.SignatureLength {
		panic(fmt.Sprintf("wrong size for signature: got %d, want %d", len(sig), crypto.SignatureLength))
	}
	r = new(big.Int).SetBytes(sig[:32])
	s = new(big.Int).SetBytes(sig[32:64])
	v = new(big.Int).SetBytes([]byte{sig[64] + 27})
	return r, s, v, nil
}

// Hash returns the hash to be signed by the sender.
// It does not uniquely identify the transaction.
func (fs FrontierSigner) Hash(tx *CeloTransaction) common.Hash {
	if tx.data.EthCompatible {
		return rlpHash([]interface{}{
			tx.data.AccountNonce,
			tx.data.Price,
			tx.data.GasLimit,
			tx.data.Recipient,
			tx.data.Amount,
			tx.data.Payload,
		})
	} else {
		return rlpHash([]interface{}{
			tx.data.AccountNonce,
			tx.data.Price,
			tx.data.GasLimit,
			tx.data.FeeCurrency,
			tx.data.GatewayFeeRecipient,
			tx.data.GatewayFee,
			tx.data.Recipient,
			tx.data.Amount,
			tx.data.Payload,
		})
	}
}

func (fs FrontierSigner) Sender(tx *CeloTransaction) (common.Address, error) {
	addr, _, err := recoverPlain(fs.Hash(tx), tx.data.R, tx.data.S, tx.data.V, false)
	return addr, err
}

func recoverPlain(sighash common.Hash, R, S, Vb *big.Int, homestead bool) (common.Address, []byte, error) {
	if Vb.BitLen() > 8 {
		return common.Address{}, nil, ErrInvalidSig
	}
	V := byte(Vb.Uint64() - 27)
	if !crypto.ValidateSignatureValues(V, R, S, homestead) {
		return common.Address{}, nil, ErrInvalidSig
	}
	// encode the signature in uncompressed format
	r, s := R.Bytes(), S.Bytes()
	sig := make([]byte, crypto.SignatureLength)
	copy(sig[32-len(r):32], r)
	copy(sig[64-len(s):64], s)
	sig[64] = V
	// recover the public key from the signature
	pub, err := crypto.Ecrecover(sighash[:], sig)
	if err != nil {
		return common.Address{}, nil, err
	}
	if len(pub) == 0 || pub[0] != 4 {
		return common.Address{}, nil, errors.New("invalid public key")
	}
	var addr common.Address
	copy(addr[:], crypto.Keccak256(pub[1:])[12:])
	return addr, pub, nil
}

// deriveChainId derives the chain id from the given v parameter
func deriveChainId(v *big.Int) *big.Int {
	if v.BitLen() <= 64 {
		v := v.Uint64()
		if v == 27 || v == 28 {
			return new(big.Int)
		}
		return new(big.Int).SetUint64((v - 35) / 2)
	}
	v = new(big.Int).Sub(v, big.NewInt(35))
	return v.Div(v, big.NewInt(2))
}
