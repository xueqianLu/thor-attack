package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"flag"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/vechain/thor/api/transactions"
	"github.com/vechain/thor/thor"
	"github.com/vechain/thor/tx"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"strings"
	"time"
)

var (
	accountFile  = flag.String("account", "/root/account.json", "account file")
	accountIndex = flag.Int("index", 0, "account index")
	url          = flag.String("url", "http://127.0.0.1:8669", "node rpc url")
)

type AccountInfo struct {
	Address string `json:"address"`
	Private string `json:"private"`
}

var (
	chainTag = byte(0)
)

type JsonBlockInfo struct {
	Number int    `json:"number"`
	ID     string `json:"id"`
}

func getChainTag(url string) {
	for {
		getBlockZeroInfo := httpGet(url + "/blocks/0")
		if len(getBlockZeroInfo) == 0 {
			time.Sleep(time.Second)
			continue
		}
		var blockInfo JsonBlockInfo
		if err := json.Unmarshal(getBlockZeroInfo, &blockInfo); err != nil {
			log.Printf("getChainTag json.Unmarshal: %v", err)
			time.Sleep(time.Second)
			continue
		}
		if strings.HasPrefix(blockInfo.ID, "0x") {
			blockInfo.ID = blockInfo.ID[2:]
		}
		idbytes, _ := hex.DecodeString(blockInfo.ID)
		chainTag = idbytes[len(idbytes)-1]
		break
	}

}

// implement function load account
func loadAccount() *AccountInfo {
	data, err := ioutil.ReadFile(*accountFile)
	if err != nil {
		log.Fatalf("ioutil.ReadFile: %v", err)
	}
	var accounts []*AccountInfo
	if err = json.Unmarshal(data, &accounts); err != nil {
		log.Fatalf("loadAccount json.Unmarshal: %v", err)
	}
	return accounts[*accountIndex]
}

func main() {
	flag.Parse()
	account := loadAccount()
	tc := time.NewTicker(time.Second * 5)
	nonce := uint64(time.Now().Unix())
	defer tc.Stop()
	getChainTag(*url)
	for {
		select {
		case <-tc.C:
			sendTx(*url, account, nonce)
			nonce++
			tc.Reset(time.Second)
		}
	}
}

func sendTx(url string, account *AccountInfo, nonce uint64) {
	addr := thor.BytesToAddress([]byte("to"))
	cla := tx.NewClause(&addr).WithValue(big.NewInt(100000000000000000))
	tx := new(tx.Builder).
		ChainTag(chainTag).
		GasPriceCoef(1).
		Expiration(10000).
		Gas(21000).
		Nonce(nonce).
		Clause(cla).
		BlockRef(tx.NewBlockRef(0)).
		Build()

	pk, err := crypto.HexToECDSA(account.Private)

	sig, err := crypto.Sign(tx.SigningHash().Bytes(), pk)
	if err != nil {
		log.Fatalf("crypto.Sign: %v", err)
	}
	tx = tx.WithSignature(sig)
	rlpTx, err := rlp.EncodeToBytes(tx)
	if err != nil {
		log.Fatalf("rlp.EncodeToBytes: %v", err)
	}

	res := httpPost(url+"/transactions", transactions.RawTx{Raw: hexutil.Encode(rlpTx)})
	var txObj map[string]string
	if err = json.Unmarshal(res, &txObj); err != nil {
		log.Fatalf("parse transaction response json.Unmarshal: %v", err)
	}
	log.Printf("txid: %s", txObj["id"])
}

func httpPost(url string, obj interface{}) []byte {
	data, err := json.Marshal(obj)
	if err != nil {
		log.Fatalf("json.Marshal: %v", err)
	}
	res, err := http.Post(url, "application/x-www-form-urlencoded", bytes.NewReader(data))
	if err != nil {
		log.Fatalf("http.Post: %v", err)
	}
	r, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatalf("ioutil.ReadAll: %v", err)
	}
	log.Printf("response: %s", string(r))
	return r
}

func httpGet(url string) []byte {
	res, err := http.Get(url)
	if err != nil {
		log.Printf("http.Get: %v", err)
		return nil
	}
	r, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Printf("ioutil.ReadAll: %v", err)
		return nil
	}
	return r
}
