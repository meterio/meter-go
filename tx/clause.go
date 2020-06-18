// Copyright (c) 2020 The Meter developers

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package tx

import (
	"fmt"
	"io"
	"math/big"

	"meter-go/meter"

	"github.com/ethereum/go-ethereum/rlp"
)

type clauseBody struct {
	To    *meter.Address `rlp:"nil"`
	Value *big.Int
	Token byte
	Data  []byte
}

// Clause is the basic execution unit of a transaction.
type Clause struct {
	body clauseBody
}

// NewClause create a new clause instance.
func NewClause(to *meter.Address) *Clause {
	if to != nil {
		// make a copy of 'to'
		cpy := *to
		to = &cpy
	}
	return &Clause{
		clauseBody{
			to,
			&big.Int{},
			byte(MeterToken),
			nil,
		},
	}
}

// WithValue create a new clause copy with value changed.
func (c *Clause) WithValue(value *big.Int) *Clause {
	newClause := *c
	newClause.body.Value = new(big.Int).Set(value)
	return &newClause
}

// WithData create a new clause copy with data changed.
func (c *Clause) WithData(data []byte) *Clause {
	newClause := *c
	newClause.body.Data = append([]byte(nil), data...)
	return &newClause
}

// WithToken create a new clause copy with value changed.
func (c *Clause) WithToken(token byte) *Clause {
	newClause := *c
	newClause.body.Token = token
	return &newClause
}

// To returns 'To' address.
func (c *Clause) To() *meter.Address {
	if c.body.To == nil {
		return nil
	}
	cpy := *c.body.To
	return &cpy
}

// Value returns 'Value'.
func (c *Clause) Value() *big.Int {
	return new(big.Int).Set(c.body.Value)
}

// Data returns 'Data'.
func (c *Clause) Data() []byte {
	return append([]byte(nil), c.body.Data...)
}

// Data returns 'Token'.
func (c *Clause) Token() byte {
	return c.body.Token
}

// EncodeRLP implements rlp.Encoder
func (c *Clause) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, &c.body)
}

// DecodeRLP implements rlp.Decoder
func (c *Clause) DecodeRLP(s *rlp.Stream) error {
	var body clauseBody
	if err := s.Decode(&body); err != nil {
		return err
	}
	*c = Clause{body}
	return nil
}

func (c *Clause) String() string {
	var to string
	if c.body.To == nil {
		to = "nil"
	} else {
		to = c.body.To.String()
	}
	return fmt.Sprintf(`
    (To: %v, Value: %v, Token: %v, Data: 0x%x)`, to, c.body.Value, c.body.Token, c.body.Data)
}
