package blockcache

import (
	"github.com/vechain/thor/thor"
	"sync"
)

type cacheMinedBlock struct {
	blocks map[string]bool
	mux    sync.Mutex
}

var (
	cacheMinedBlockInstance *cacheMinedBlock
)

func init() {
	cacheMinedBlockInstance = &cacheMinedBlock{
		blocks: make(map[string]bool),
	}
}

func UpdateBlockBroadcasted(id thor.Bytes32) {
	cacheMinedBlockInstance.mux.Lock()
	defer cacheMinedBlockInstance.mux.Unlock()
	cacheMinedBlockInstance.blocks[id.String()] = true
}

func GetBlockBroadCasted(id thor.Bytes32) (bool, bool) {
	cacheMinedBlockInstance.mux.Lock()
	defer cacheMinedBlockInstance.mux.Unlock()
	exist, ok := cacheMinedBlockInstance.blocks[id.String()]
	return exist, ok
}

func AddNewBlock(id thor.Bytes32) {
	cacheMinedBlockInstance.mux.Lock()
	defer cacheMinedBlockInstance.mux.Unlock()
	cacheMinedBlockInstance.blocks[id.String()] = false
}
