package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/vechain/thor/api/blocks"
	"io/ioutil"
	"log"
	"net/http"
)

func blockInfo(url string, number int64) *blocks.JSONBlockSummary {
	rb := new(blocks.JSONCollapsedBlock)
	reqUrl := fmt.Sprintf("%s/blocks/%d", url, number)
	res := httpGet(reqUrl)
	if err := json.Unmarshal(res, &rb); err != nil {
		log.Fatalf("json.Unmarshal: %v", err)
	}
	return rb.JSONBlockSummary
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
