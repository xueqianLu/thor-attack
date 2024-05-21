package staticpeer

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/vechain/thor/p2psrv"
	"log"
	"os"
)

func ParseAndAddStaticPeers(server *p2psrv.Server, nodefile string) {
	var nodeInfos = make([]string, 0)
	if data, err := os.ReadFile(nodefile); err != nil {
		log.Println("ParseAndAddStaticPeers Error reading file: ", err)
		return
	} else {
		if err := json.Unmarshal(data, &nodeInfos); err != nil {
			log.Println("ParseAndAddStaticPeers Error parsing json: ", err)
			return
		}
		for _, url := range nodeInfos {
			if node, err := discover.ParseNode(url); err == nil {
				server.AddStatic(node)
				log.Println("ParseAndAddStaticPeers Added node: ", url)
			}
		}
	}
}
