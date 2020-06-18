// Copyright (c) 2020 The Meter developers

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package meter

import (
	"encoding/hex"
	"errors"
	"strings"
)

// Bytes32 array of 32 bytes.
type Bytes32 [32]byte

// String implements stringer
func (b Bytes32) String() string {
	return "0x" + hex.EncodeToString(b[:])
}

// Bytes returns byte slice form of Bytes32.
func (b Bytes32) Bytes() []byte {
	return b[:]
}

// ParseBytes32 convert string presented into Bytes32 type
func ParseBytes32(s string) (Bytes32, error) {
	if len(s) == 32*2 {
	} else if len(s) == 32*2+2 {
		if strings.ToLower(s[:2]) != "0x" {
			return Bytes32{}, errors.New("invalid prefix")
		}
		s = s[2:]
	} else {
		return Bytes32{}, errors.New("invalid length")
	}

	var b Bytes32
	_, err := hex.Decode(b[:], []byte(s))
	if err != nil {
		return Bytes32{}, err
	}
	return b, nil
}

// MustParseBytes32 convert string presented into Bytes32 type, panic on error.
func MustParseBytes32(s string) Bytes32 {
	b32, err := ParseBytes32(s)
	if err != nil {
		panic(err)
	}
	return b32
}
