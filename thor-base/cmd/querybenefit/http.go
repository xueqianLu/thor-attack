package main

import (
	"bytes"
	"encoding/json"
	"github.com/vechain/thor/api/accounts"
	"github.com/vechain/thor/api/blocks"
	"io/ioutil"
	"log"
	"net/http"
)

func bestBlock(url string) *blocks.JSONBlockSummary {
	rb := new(blocks.JSONCollapsedBlock)
	res := httpGet(url + "/blocks/best")
	if err := json.Unmarshal(res, &rb); err != nil {
		log.Fatalf("json.Unmarshal: %v", err)
	}
	return rb.JSONBlockSummary
}

func accountInfo(url string, addr string, revision string) accounts.Account {
	res := httpGet(url + "/accounts/" + addr + "?revision=" + revision)
	var acc accounts.Account
	if err := json.Unmarshal(res, &acc); err != nil {
		log.Fatalf("json.Unmarshal: %v", err)
	}
	return acc
}

func httpGet(url string) []byte {
	res, err := http.Get(url)
	if err != nil {
		log.Fatalf("http.Get: %v", err)
	}
	r, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatalf("ioutil.ReadAll: %v", err)
	}
	return r
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
	return r
}
