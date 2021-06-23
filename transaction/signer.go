package transaction

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var ErrInvalidSig = errors.New("invalid transaction v, r, s values")

type CeloSigner struct{}

//
//func (s CeloSigner) Equal(s2 Signer) bool {
//	_, ok := s2.(CeloSigner)
//	return ok
//}

// SignatureValues returns signature values. This signature
// needs to be in the [R || S || V] format where V is 0 or 1.
func (fs CeloSigner) SignatureValues(tx *CeloTransaction, sig []byte) (r, s, v *big.Int, err error) {
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
func (fs CeloSigner) Hash(tx *CeloTransaction) common.Hash {
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

func (fs CeloSigner) Sender(tx *CeloTransaction) (common.Address, error) {
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
