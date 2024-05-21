package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strings"
)

var (
	verifyUrl  = flag.String("url", "http://localhost:8669", "rest url")
	hackBlocks = flag.String("blocks", "", "hack blocks info file")
)

func main() {
	flag.Parse()
	content, err := os.ReadFile(*hackBlocks)
	if err != nil {
		panic(err)
	}
	var hackBlockList []blockSummary
	if err := json.Unmarshal(content, &hackBlockList); err != nil {
		panic(err)
	}
	ratio := verifyBlock(hackBlockList, *verifyUrl)
	log.Println("verify ratio:", ratio)
}

type blockSummary struct {
	Number int64  `json:"number"`
	Id     string `json:"id"`
}

func verifyBlock(hackBlockList []blockSummary, verifyNodeUrl string) float32 {
	hackBlockCount := len(hackBlockList)
	succeedCount := 0
	for _, hackBlock := range hackBlockList {
		remoteBlock := blockInfo(verifyNodeUrl, hackBlock.Number)
		if remoteBlock != nil && strings.Contains(remoteBlock.ID.String(), hackBlock.Id) {
			succeedCount++
		}
	}
	return float32(succeedCount) / float32(hackBlockCount)
}
