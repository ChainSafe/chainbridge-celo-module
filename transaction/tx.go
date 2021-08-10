package transaction

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"io"
	"math/big"
	"sync/atomic"

	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
)

const (
	ethCompatibleTxNumFields = 9
)

var (
	ErrInvalidSig = errors.New("invalid transaction v, r, s values")
	// ErrEthCompatibleTransactionIsntCompatible is returned if the transaction has EthCompatible: true
	// but has non-nil-or-0 values for some of the Celo-only fields
	ErrEthCompatibleTransactionIsntCompatible = errors.New("ethCompatible is true, but non-eth-compatible fields are present")
)

func NewCeloTransaction(nonce uint64, to *common.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte) evmclient.CommonTransaction {
	return newTransaction(nonce, to, amount, gasLimit, gasPrice, nil, nil, nil, data)
}

type CeloTransaction struct {
	data txdata
	// caches
	hash atomic.Value
	size atomic.Value
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

type txdataMarshaling struct {
	AccountNonce        hexutil.Uint64
	Price               *hexutil.Big
	GasLimit            hexutil.Uint64
	FeeCurrency         *common.Address
	GatewayFeeRecipient *common.Address
	GatewayFee          *hexutil.Big
	Amount              *hexutil.Big
	Payload             hexutil.Bytes
	V                   *hexutil.Big
	R                   *hexutil.Big
	S                   *hexutil.Big
	EthCompatible       bool
}

// ethCompatibleTxRlpList is used for RLP encoding/decoding of eth-compatible transactions.
// As such, it:
// (a) excludes the Celo-only fields,
// (b) doesn't need the Hash or EthCompatible fields, and
// (c) doesn't need the `json` or `gencodec` tags
type ethCompatibleTxRlpList struct {
	AccountNonce uint64
	Price        *big.Int
	GasLimit     uint64
	Recipient    *common.Address `rlp:"nil"` // nil means contract creation
	Amount       *big.Int
	Payload      []byte
	V            *big.Int
	R            *big.Int
	S            *big.Int
}

func toEthCompatibleRlpList(data txdata) ethCompatibleTxRlpList {
	return ethCompatibleTxRlpList{
		AccountNonce: data.AccountNonce,
		Price:        data.Price,
		GasLimit:     data.GasLimit,
		Recipient:    data.Recipient,
		Amount:       data.Amount,
		Payload:      data.Payload,
		V:            data.V,
		R:            data.R,
		S:            data.S,
	}
}

func fromEthCompatibleRlpList(data ethCompatibleTxRlpList) txdata {
	return txdata{
		AccountNonce:        data.AccountNonce,
		Price:               data.Price,
		GasLimit:            data.GasLimit,
		FeeCurrency:         nil,
		GatewayFeeRecipient: nil,
		GatewayFee:          big.NewInt(0),
		Recipient:           data.Recipient,
		Amount:              data.Amount,
		Payload:             data.Payload,
		V:                   data.V,
		R:                   data.R,
		S:                   data.S,
		Hash:                nil, // txdata.Hash is calculated and saved inside tx.Hash()
		EthCompatible:       true,
	}
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

// EncodeRLP implements celorlp.Encoder
func (tx *CeloTransaction) EncodeRLP(w io.Writer) error {
	if tx.data.EthCompatible {
		return rlp.Encode(w, toEthCompatibleRlpList(tx.data))
	} else {
		return rlp.Encode(w, &tx.data)
	}
}

// DecodeRLP implements celorlp.Decoder
func (tx *CeloTransaction) DecodeRLP(s *rlp.Stream) (err error) {
	_, size, _ := s.Kind()
	var raw rlp.RawValue
	err = s.Decode(&raw)
	if err != nil {
		return err
	}
	headerSize := len(raw) - int(size)
	numElems, err := rlp.CountValues(raw[headerSize:])
	if err != nil {
		return err
	}
	if numElems == ethCompatibleTxNumFields {
		rlpList := ethCompatibleTxRlpList{}
		err = rlp.DecodeBytes(raw, &rlpList)
		tx.data = fromEthCompatibleRlpList(rlpList)
	} else {
		err = rlp.DecodeBytes(raw, &tx.data)
	}
	if err == nil {
		tx.size.Store(common.StorageSize(rlp.ListSize(size)))
	}

	return err
}

// MarshalJSON encodes the web3 RPC transaction format.
func (tx *CeloTransaction) MarshalJSON() ([]byte, error) {
	hash := tx.Hash()
	data := tx.data
	data.Hash = &hash
	return data.MarshalJSON()
}

// UnmarshalJSON decodes the web3 RPC transaction format.
func (tx *CeloTransaction) UnmarshalJSON(input []byte) error {
	var dec txdata
	if err := dec.UnmarshalJSON(input); err != nil {
		return err
	}

	withSignature := dec.V.Sign() != 0 || dec.R.Sign() != 0 || dec.S.Sign() != 0
	if withSignature {
		var V byte
		if isProtectedV(dec.V) {
			chainID := deriveChainId(dec.V).Uint64()
			V = byte(dec.V.Uint64() - 35 - 2*chainID)
		} else {
			V = byte(dec.V.Uint64() - 27)
		}
		if !crypto.ValidateSignatureValues(V, dec.R, dec.S, false) {
			return ErrInvalidSig
		}
	}

	*tx = CeloTransaction{data: dec}
	return nil
}

func (tx *CeloTransaction) Data() []byte                         { return common.CopyBytes(tx.data.Payload) }
func (tx *CeloTransaction) Gas() uint64                          { return tx.data.GasLimit }
func (tx *CeloTransaction) GasPrice() *big.Int                   { return new(big.Int).Set(tx.data.Price) }
func (tx *CeloTransaction) FeeCurrency() *common.Address         { return tx.data.FeeCurrency }
func (tx *CeloTransaction) GatewayFeeRecipient() *common.Address { return tx.data.GatewayFeeRecipient }
func (tx *CeloTransaction) GatewayFee() *big.Int                 { return tx.data.GatewayFee }
func (tx *CeloTransaction) Value() *big.Int                      { return new(big.Int).Set(tx.data.Amount) }
func (tx *CeloTransaction) Nonce() uint64                        { return tx.data.AccountNonce }
func (tx *CeloTransaction) CheckNonce() bool                     { return true }
func (tx *CeloTransaction) EthCompatible() bool                  { return tx.data.EthCompatible }
func (tx *CeloTransaction) Fee() *big.Int {
	gasFee := new(big.Int).Mul(tx.data.Price, big.NewInt(int64(tx.data.GasLimit)))
	return gasFee.Add(gasFee, tx.data.GatewayFee)
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

// Size returns the true RLP encoded storage size of the transaction, either by
// encoding and returning it, or returning a previsouly cached value.
func (tx *CeloTransaction) Size() common.StorageSize {
	if size := tx.size.Load(); size != nil {
		return size.(common.StorageSize)
	}
	c := writeCounter(0)
	rlp.Encode(&c, &tx.data)
	tx.size.Store(common.StorageSize(c))
	return common.StorageSize(c)
}

// CheckEthCompatibility checks that the Celo-only fields are nil-or-0 if EthCompatible is true
func (tx *CeloTransaction) CheckEthCompatibility() error {
	if tx.EthCompatible() && !(tx.FeeCurrency() == nil && tx.GatewayFeeRecipient() == nil && tx.GatewayFee().Sign() == 0) {
		return ErrEthCompatibleTransactionIsntCompatible
	}
	return nil
}

// AsMessage returns the transaction as a core.Message.
//
// AsMessage requires a signer to derive the sender.
//
// XXX Rename message to something less arbitrary?
func (tx *CeloTransaction) AsMessage(s CeloSigner) (Message, error) {
	msg := Message{
		nonce:               tx.data.AccountNonce,
		gasLimit:            tx.data.GasLimit,
		gasPrice:            new(big.Int).Set(tx.data.Price),
		feeCurrency:         tx.data.FeeCurrency,
		gatewayFeeRecipient: tx.data.GatewayFeeRecipient,
		gatewayFee:          tx.data.GatewayFee,
		to:                  tx.data.Recipient,
		amount:              tx.data.Amount,
		data:                tx.data.Payload,
		checkNonce:          true,
		ethCompatible:       tx.data.EthCompatible,
	}

	var err error
	msg.from, err = Sender(s, tx)
	return msg, err
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

// Cost returns amount + gasprice * gaslimit + gatewayfee.
func (tx *CeloTransaction) Cost() *big.Int {
	total := new(big.Int).Mul(tx.data.Price, new(big.Int).SetUint64(tx.data.GasLimit))
	total.Add(total, tx.data.Amount)
	total.Add(total, tx.data.GatewayFee)
	return total
}

// RawSignatureValues returns the V, R, S signature values of the transaction.
// The return values should not be modified by the caller.
func (tx *CeloTransaction) RawSignatureValues() (v, r, s *big.Int) {
	return tx.data.V, tx.data.R, tx.data.S
}

// Transactions is a Transaction slice type for basic sorting.
type CeloTransactions []*CeloTransaction

// Len returns the length of s.
func (s CeloTransactions) Len() int { return len(s) }

// Swap swaps the i'th and the j'th element in s.
func (s CeloTransactions) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// GetRlp implements Rlpable and returns the i'th element of s in celorlp.
func (s CeloTransactions) GetRlp(i int) []byte {
	enc, _ := rlp.EncodeToBytes(s[i])
	return enc
}

// TxDifference returns a new set which is the difference between a and b.
func TxDifference(a, b CeloTransactions) CeloTransactions {
	keep := make(CeloTransactions, 0, len(a))

	remove := make(map[common.Hash]struct{})
	for _, tx := range b {
		remove[tx.Hash()] = struct{}{}
	}

	for _, tx := range a {
		if _, ok := remove[tx.Hash()]; !ok {
			keep = append(keep, tx)
		}
	}

	return keep
}

// Message is a fully derived transaction and implements core.Message
//
// NOTE: In a future PR this will be removed.
type Message struct {
	to                  *common.Address
	from                common.Address
	nonce               uint64
	amount              *big.Int
	gasLimit            uint64
	gasPrice            *big.Int
	feeCurrency         *common.Address
	gatewayFeeRecipient *common.Address
	gatewayFee          *big.Int
	data                []byte
	ethCompatible       bool
	checkNonce          bool
}

func NewMessage(from common.Address, to *common.Address, nonce uint64, amount *big.Int, gasLimit uint64, gasPrice *big.Int, feeCurrency, gatewayFeeRecipient *common.Address, gatewayFee *big.Int, data []byte, ethCompatible, checkNonce bool) Message {
	return Message{
		from:                from,
		to:                  to,
		nonce:               nonce,
		amount:              amount,
		gasLimit:            gasLimit,
		gasPrice:            gasPrice,
		feeCurrency:         feeCurrency,
		gatewayFeeRecipient: gatewayFeeRecipient,
		gatewayFee:          gatewayFee,
		data:                data,
		ethCompatible:       ethCompatible,
		checkNonce:          checkNonce,
	}
}

func (m Message) From() common.Address                 { return m.from }
func (m Message) To() *common.Address                  { return m.to }
func (m Message) GasPrice() *big.Int                   { return m.gasPrice }
func (m Message) EthCompatible() bool                  { return m.ethCompatible }
func (m Message) FeeCurrency() *common.Address         { return m.feeCurrency }
func (m Message) GatewayFeeRecipient() *common.Address { return m.gatewayFeeRecipient }
func (m Message) GatewayFee() *big.Int                 { return m.gatewayFee }
func (m Message) Value() *big.Int                      { return m.amount }
func (m Message) Gas() uint64                          { return m.gasLimit }
func (m Message) Nonce() uint64                        { return m.nonce }
func (m Message) Data() []byte                         { return m.data }
func (m Message) CheckNonce() bool                     { return m.checkNonce }
func (m Message) Fee() *big.Int {
	gasFee := new(big.Int).Mul(m.gasPrice, big.NewInt(int64(m.gasLimit)))
	return gasFee.Add(gasFee, m.gatewayFee)
}

var _ = (*txdataMarshaling)(nil)

// MarshalJSON marshals as JSON.
func (t txdata) MarshalJSON() ([]byte, error) {
	type txdata struct {
		AccountNonce        hexutil.Uint64  `json:"nonce"    gencodec:"required"`
		Price               *hexutil.Big    `json:"gasPrice" gencodec:"required"`
		GasLimit            hexutil.Uint64  `json:"gas"      gencodec:"required"`
		FeeCurrency         *common.Address `json:"feeCurrency" rlp:"nil"`
		GatewayFeeRecipient *common.Address `json:"gatewayFeeRecipient" rlp:"nil"`
		GatewayFee          *hexutil.Big    `json:"gatewayFee"`
		Recipient           *common.Address `json:"to"       rlp:"nil"`
		Amount              *hexutil.Big    `json:"value"    gencodec:"required"`
		Payload             hexutil.Bytes   `json:"input"    gencodec:"required"`
		V                   *hexutil.Big    `json:"v" gencodec:"required"`
		R                   *hexutil.Big    `json:"r" gencodec:"required"`
		S                   *hexutil.Big    `json:"s" gencodec:"required"`
		Hash                *common.Hash    `json:"hash" rlp:"-"`
		EthCompatible       bool            `json:"ethCompatible" rlp:"-"`
	}
	var enc txdata
	enc.AccountNonce = hexutil.Uint64(t.AccountNonce)
	enc.Price = (*hexutil.Big)(t.Price)
	enc.GasLimit = hexutil.Uint64(t.GasLimit)
	enc.FeeCurrency = t.FeeCurrency
	enc.GatewayFeeRecipient = t.GatewayFeeRecipient
	enc.GatewayFee = (*hexutil.Big)(t.GatewayFee)
	enc.Recipient = t.Recipient
	enc.Amount = (*hexutil.Big)(t.Amount)
	enc.Payload = t.Payload
	enc.V = (*hexutil.Big)(t.V)
	enc.R = (*hexutil.Big)(t.R)
	enc.S = (*hexutil.Big)(t.S)
	enc.Hash = t.Hash
	enc.EthCompatible = t.EthCompatible
	return json.Marshal(&enc)
}

// UnmarshalJSON unmarshals from JSON.
func (t *txdata) UnmarshalJSON(input []byte) error {
	type txdata struct {
		AccountNonce        *hexutil.Uint64 `json:"nonce"    gencodec:"required"`
		Price               *hexutil.Big    `json:"gasPrice" gencodec:"required"`
		GasLimit            *hexutil.Uint64 `json:"gas"      gencodec:"required"`
		FeeCurrency         *common.Address `json:"feeCurrency" rlp:"nil"`
		GatewayFeeRecipient *common.Address `json:"gatewayFeeRecipient" rlp:"nil"`
		GatewayFee          *hexutil.Big    `json:"gatewayFee"`
		Recipient           *common.Address `json:"to"       rlp:"nil"`
		Amount              *hexutil.Big    `json:"value"    gencodec:"required"`
		Payload             *hexutil.Bytes  `json:"input"    gencodec:"required"`
		V                   *hexutil.Big    `json:"v" gencodec:"required"`
		R                   *hexutil.Big    `json:"r" gencodec:"required"`
		S                   *hexutil.Big    `json:"s" gencodec:"required"`
		Hash                *common.Hash    `json:"hash" rlp:"-"`
		EthCompatible       *bool           `json:"ethCompatible" rlp:"-"`
	}
	var dec txdata
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.AccountNonce == nil {
		return errors.New("missing required field 'nonce' for txdata")
	}
	t.AccountNonce = uint64(*dec.AccountNonce)
	if dec.Price == nil {
		return errors.New("missing required field 'gasPrice' for txdata")
	}
	t.Price = (*big.Int)(dec.Price)
	if dec.GasLimit == nil {
		return errors.New("missing required field 'gas' for txdata")
	}
	t.GasLimit = uint64(*dec.GasLimit)
	if dec.FeeCurrency != nil {
		t.FeeCurrency = dec.FeeCurrency
	}
	if dec.GatewayFeeRecipient != nil {
		t.GatewayFeeRecipient = dec.GatewayFeeRecipient
	}
	if dec.GatewayFee != nil {
		t.GatewayFee = (*big.Int)(dec.GatewayFee)
	}
	if dec.Recipient != nil {
		t.Recipient = dec.Recipient
	}
	if dec.Amount == nil {
		return errors.New("missing required field 'value' for txdata")
	}
	t.Amount = (*big.Int)(dec.Amount)
	if dec.Payload == nil {
		return errors.New("missing required field 'input' for txdata")
	}
	t.Payload = *dec.Payload
	if dec.V == nil {
		return errors.New("missing required field 'v' for txdata")
	}
	t.V = (*big.Int)(dec.V)
	if dec.R == nil {
		return errors.New("missing required field 'r' for txdata")
	}
	t.R = (*big.Int)(dec.R)
	if dec.S == nil {
		return errors.New("missing required field 's' for txdata")
	}
	t.S = (*big.Int)(dec.S)
	if dec.Hash != nil {
		t.Hash = dec.Hash
	}
	if dec.EthCompatible != nil {
		t.EthCompatible = *dec.EthCompatible
	}
	return nil
}

func rlpHash(x interface{}) (h common.Hash) {
	hw := sha3.NewLegacyKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}

type writeCounter common.StorageSize

func (c *writeCounter) Write(b []byte) (int, error) {
	*c += writeCounter(len(b))
	return len(b), nil
}
