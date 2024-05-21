package main

import (
	"encoding/csv"
	"flag"
	"log"
	"math/big"
	"os"
	"strconv"
)

var (
	restUrl = flag.String("url", "http://localhost:8669", "rest url")
	report  = flag.String("report", "/root/node/report.csv", "report file")
)

var (
	nodeBenefitList = []struct {
		nodeName string
		benefit  string
	}{
		{"node0", "0x0000000000000000000000000000000000000010"},
		{"node1", "0x0000000000000000000000000000000000000011"},
		{"node2", "0x0000000000000000000000000000000000000012"},
		{"node3", "0x0000000000000000000000000000000000000013"},
		{"node4", "0x0000000000000000000000000000000000000014"},
		{"node5", "0x0000000000000000000000000000000000000015"},
		{"node6", "0x0000000000000000000000000000000000000016"},
	}
)

func main() {
	flag.Parse()
	height := bestBlock(*restUrl).Number
	f, err := os.OpenFile(*report, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		return
	}
	defer func() {
		f.Close()
	}()
	write := csv.NewWriter(f)
	write.Write([]string{"block", "node0", "node1", "node2", "node3", "node4", "node5", "node6"})
	for i := uint32(0); i < height; i++ {
		record := make([]string, 0)
		record = append(record, strconv.Itoa(int(i)))
		for _, node := range nodeBenefitList {
			acc := accountInfo(*restUrl, node.benefit, strconv.Itoa(int(i)))
			energy := big.Int(acc.Energy)
			record = append(record, energy.Text(10))
			log.Printf("%s has benefit %v at block \t%d", node.nodeName, energy.Text(10), i)
		}
		write.Write(record)
	}
	write.Flush()
}
