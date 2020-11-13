package fil

import (
	"encoding/json"
	"errors"
	addr "github.com/filecoin-project/go-address"
	"github.com/ipfs/go-cid"
	"math/big"
	"time"
)

func init() {
	addr.CurrentNetwork = addr.Mainnet
}

// ErrInsuffientFunds is an error that should be returned by WalletBackend.Send method
// if the address does not have enough funds.
var ErrInsuffientFunds = errors.New("insufficient funds")

const attoFilPerFilecoin = 1000000000000000000

// Transaction represents a Filecoin transaction.
type Transaction struct {
	ID        cid.Cid
	From      addr.Address
	To        addr.Address
	Amount    *big.Int
	Timestamp time.Time
}

// MarshalJSON is a custom JSON marshaller for Transaction.
func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID        string    `json:"transactionID"`
		From      string    `json:"from"`
		To        string    `json:"to"`
		Amount    float64   `json:"amount"`
		Timestamp time.Time `json:"timestamp"`
	}{
		ID:        t.ID.String(),
		From:      t.From.String(),
		To:        t.To.String(),
		Amount:    AttoFILToFIL(t.Amount),
		Timestamp: t.Timestamp,
	})
}

// WalletBackend is an interface for a Filecoin wallet that can hold the keys
// for multiple addresses and can make transactions.
type WalletBackend interface {
	// NewAddress generates a new address and store the key in the backend.
	NewAddress() (addr.Address, error)

	// Send filecoin from one address to another. Returns the cid of the
	// transaction.
	Send(from, to addr.Address, amount *big.Int) (cid.Cid, error)

	// Balance returns the balance for an address.
	Balance(addr addr.Address) (*big.Int, error)

	// Transactions returns the list of transactions for an address.
	Transactions(addr addr.Address, limit, offset int) ([]Transaction, error)
}

// FILtoAttoFIL converts a float containing an amount of Filecoin to
// a big.Int representation in the attoFil base unit.
func FILtoAttoFIL(fil float64) *big.Int {
	bigAtto := big.NewFloat(fil)
	bigAtto = bigAtto.Mul(bigAtto, big.NewFloat(attoFilPerFilecoin))
	ret, acc := bigAtto.Int(nil)

	if acc == big.Above {
		ret.Add(ret, big.NewInt(-1))
	}
	if acc == big.Below {
		ret.Add(ret, big.NewInt(1))
	}

	return ret
}

// AttoFILtoFIL converts a big.Int containing the attoFIL base unit to
// a float of the amount of Filecoin.
func AttoFILToFIL(attoFIL *big.Int) float64 {
	bigAtto := new(big.Float).SetInt(attoFIL)
	bigAtto = bigAtto.Quo(bigAtto, new(big.Float).SetInt(big.NewInt(attoFilPerFilecoin)))
	ret, _ := bigAtto.Float64()
	return ret
}
