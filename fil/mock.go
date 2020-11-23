package fil

import (
	"crypto/rand"
	addr "github.com/filecoin-project/go-address"
	"github.com/gcash/bchd/bchec"
	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multihash"
	"math/big"
	"sync"
	"time"
)

// MockWalletBackend is a mock backend for the wallet that allows
// for making mock transactions and generating mock blocks.
type MockWalletBackend struct {
	transactions map[addr.Address][]Transaction
	nextAddr     *addr.Address
	nextTxid     *cid.Cid
	nextTime     *time.Time
	mtx          sync.RWMutex
}

// NewMockWalletBackend instantiates a new WalletBackend.
func NewMockWalletBackend() *MockWalletBackend {
	return &MockWalletBackend{
		transactions: make(map[addr.Address][]Transaction),
		mtx:          sync.RWMutex{},
	}
}

// GenerateToAddress creates mock coins and sends them to the address.
func (w *MockWalletBackend) GenerateToAddress(addr addr.Address, amount *big.Int) {
	w.mtx.Lock()
	defer w.mtx.Unlock()

	var txid cid.Cid
	if w.nextTxid != nil {
		txid = *w.nextTxid
		w.nextTxid = nil
	} else {
		txid, _ = randTxid()
	}

	ts := time.Now()
	if w.nextTime != nil {
		ts = *w.nextTime
	}

	tx := Transaction{
		ID:        txid,
		To:        addr,
		Timestamp: ts,
		Amount:    amount,
	}

	w.transactions[addr] = append(w.transactions[addr], tx)
}

// NewAddress generates a new address and store the key in the backend.
func (w *MockWalletBackend) NewAddress() (addr.Address, error) {
	w.mtx.Lock()
	defer w.mtx.Unlock()

	if w.nextAddr != nil {
		ret := *w.nextAddr
		w.nextAddr = nil
		return ret, nil
	}

	priv, err := bchec.NewPrivateKey(bchec.S256())
	if err != nil {
		return addr.Address{}, err
	}
	return addr.NewSecp256k1Address(priv.PubKey().SerializeCompressed())
}

func (w *MockWalletBackend) SetNextAddress(addr addr.Address) {
	w.mtx.Lock()
	defer w.mtx.Unlock()

	w.nextAddr = &addr
}

func (w *MockWalletBackend) SetNextTxid(id cid.Cid) {
	w.mtx.Lock()
	defer w.mtx.Unlock()

	w.nextTxid = &id
}

func (w *MockWalletBackend) SetNextTime(timestamp time.Time) {
	w.mtx.Lock()
	defer w.mtx.Unlock()

	w.nextTime = &timestamp
}

// Send filecoin from one address to another. Returns the cid of the
// transaction.
func (w *MockWalletBackend) Send(from, to addr.Address, amount *big.Int) (cid.Cid, error) {
	w.mtx.Lock()
	defer w.mtx.Unlock()

	balance, err := w.balance(from)
	if err != nil {
		return cid.Cid{}, err
	}

	if amount.Cmp(balance) > 0 {
		return cid.Cid{}, ErrInsuffientFunds
	}

	var txid cid.Cid
	if w.nextTxid != nil {
		txid = *w.nextTxid
		w.nextTxid = nil
	} else {
		txid, err = randTxid()
		if err != nil {
			return cid.Cid{}, err
		}
	}

	ts := time.Now()
	if w.nextTime != nil {
		ts = *w.nextTime
	}

	tx := Transaction{
		ID:        txid,
		To:        to,
		From:      from,
		Timestamp: ts,
		Amount:    amount,
	}

	w.transactions[to] = append(w.transactions[to], tx)
	if to.String() != from.String() {
		w.transactions[from] = append(w.transactions[from], tx)
	}

	return txid, nil
}

// Balance returns the balance for an address.
func (w *MockWalletBackend) Balance(addr addr.Address) (*big.Int, error) {
	w.mtx.RLock()
	defer w.mtx.RUnlock()

	return w.balance(addr)
}

func (w *MockWalletBackend) balance(addr addr.Address) (*big.Int, error) {
	incoming, outgoing := big.NewInt(0), big.NewInt(0)

	txs := w.transactions[addr]
	for _, tx := range txs {
		if tx.To.String() == addr.String() {
			incoming.Add(incoming, tx.Amount)
		} else if tx.From.String() == addr.String() {
			outgoing.Add(outgoing, tx.Amount)
		}
	}

	return new(big.Int).Sub(incoming, outgoing), nil
}

// Transactions returns the list of transactions for an address.
func (w *MockWalletBackend) Transactions(addr addr.Address, limit, offset int) ([]Transaction, error) {
	w.mtx.RLock()
	defer w.mtx.RUnlock()

	txs := w.transactions[addr]

	if limit < 0 {
		limit = len(txs)
	}
	if offset < 0 {
		offset = 0
	}
	if offset+limit > len(txs) {
		limit = len(txs) - offset
	}

	return txs[offset : offset+limit], nil
}

func randTxid() (cid.Cid, error) {
	r := make([]byte, 32)
	rand.Read(r)

	mh, err := multihash.Encode(r, multihash.SHA2_256)
	if err != nil {
		return cid.Cid{}, err
	}

	return cid.NewCidV1(cid.Raw, mh), nil
}
