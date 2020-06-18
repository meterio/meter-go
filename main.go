// Copyright (c) 2020 The Meter developers

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"

	"meter-go/meter"
	"meter-go/tx"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

type RawTx struct {
	Raw string `json:"raw"`
}

var (
	ToAddress      = meter.MustParseAddress("0xf3dd5c55b96889369f714143f213403464a268a6")
	TestPrivateKey = os.Getenv("TEST_PRIVATE_KEY") //  hex string without leading 0x
)

func sendTx(blockRef uint32) {
	chainTag := byte(88) // chainTag is NOT the same across chains
	var expiration = uint32(100)
	var gas = uint64(21000)
	clause := tx.NewClause(&ToAddress).
		WithValue(big.NewInt(2e18)).   // value in Wei
		WithToken(byte(tx.MeterToken)) // choose which token to send

	tx := new(tx.Builder).
		BlockRef(tx.NewBlockRef(blockRef)).
		ChainTag(chainTag).
		Expiration(expiration).
		GasPriceCoef(128).
		Gas(gas).
		Clause(clause).
		Nonce(1234567).
		Build()
	privKey, err := crypto.HexToECDSA(TestPrivateKey)
	if err != nil {
		fmt.Println(err)
		return
	}
	sig, err := crypto.Sign(tx.SigningHash().Bytes(), privKey)
	if err != nil {
		fmt.Println(err)
		return
	}
	tx = tx.WithSignature(sig)
	rlpTx, err := rlp.EncodeToBytes(tx)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Built Tx: ", tx.String())
	fmt.Println("Raw Tx:", hexutil.Encode(rlpTx))

	fmt.Println("Send tx to warringstakes network")
	res := httpPost("http://warringstakes.meter.io:8669/transactions", RawTx{Raw: hexutil.Encode(rlpTx)})
	fmt.Println("Received response: ", string(res))
	var txObj map[string]string
	if err = json.Unmarshal(res, &txObj); err != nil {
		fmt.Println(err)
		return
	}
}

func httpPost(url string, obj interface{}) []byte {
	data, err := json.Marshal(obj)
	if err != nil {
		fmt.Println("http post error:", err)
		return make([]byte, 0)
	}
	res, err := http.Post(url, "application/x-www-form-urlencoded", bytes.NewReader(data))
	if err != nil {
		fmt.Println("http post error:", err)
		return make([]byte, 0)
	}
	r, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		fmt.Println("http post error:", err)
		return make([]byte, 0)
	}
	return r
}

type ApiBlock struct {
	Number uint32 `json:"number"`
	ID     string `json:"id"`
	Size   uint32 `json:"size"`
}

func getBestBlock(url string) *ApiBlock {
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("http get error:", err)
		return nil
	}
	r, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	block := &ApiBlock{}
	err = json.Unmarshal(r, block)
	if err != nil {
		fmt.Println("http post error:", err)
		return nil
	}
	return block
}

func main() {
	bestBlock := getBestBlock("http://warringstakes.meter.io:8669/blocks/best")

	sendTx(bestBlock.Number)
}
