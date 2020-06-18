// Copyright (c) 2020 The Meter developers

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package meter

import (
	"encoding/hex"
	"errors"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

const (
	// AddressLength length of address in bytes.
	AddressLength = common.AddressLength
)

// Address address of account.
type Address common.Address

// String implements the stringer interface
func (a Address) String() string {
	return "0x" + hex.EncodeToString(a[:])
}

// Bytes returns byte slice form of address.
func (a Address) Bytes() []byte {
	return a[:]
}

// ParseAddress convert string presented address into Address type.
func ParseAddress(s string) (Address, error) {
	if len(s) == AddressLength*2 {
	} else if len(s) == AddressLength*2+2 {
		if strings.ToLower(s[:2]) != "0x" {
			return Address{}, errors.New("invalid prefix")
		}
		s = s[2:]
	} else {
		return Address{}, errors.New("invalid length")
	}

	var addr Address
	_, err := hex.Decode(addr[:], []byte(s))
	if err != nil {
		return Address{}, err
	}
	return addr, nil
}

// MustParseAddress convert string presented address into Address type, panic on error.
func MustParseAddress(s string) Address {
	addr, err := ParseAddress(s)
	if err != nil {
		panic(err)
	}
	return addr
}
