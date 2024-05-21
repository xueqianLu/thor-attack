package hackblock

import (
	"bytes"
	"fmt"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/fsnotify/fsnotify"
	lru "github.com/hashicorp/golang-lru"
	"github.com/pborman/uuid"
	"github.com/vechain/thor/block"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type BlockWatch struct {
	id   string
	ch   chan *block.Block
	stop func()
}

func (b *BlockWatch) Stop() {
	b.stop()
}

func (b *BlockWatch) Add(blk *block.Block) {
	select {
	case b.ch <- blk:
	default:
	}
}

func (b *BlockWatch) Watch() <-chan *block.Block {
	return b.ch
}

type HackBlock struct {
	path     string
	eventCh  chan fsnotify.Event
	watchers map[string]*BlockWatch
	mux      sync.Mutex
	cache    *lru.Cache
}

func NewHackBlock(path string) *HackBlock {
	cache, _ := lru.New(100)
	return &HackBlock{
		path:     path,
		eventCh:  make(chan fsnotify.Event, 10),
		watchers: make(map[string]*BlockWatch),
		cache:    cache,
	}
}

func (h *HackBlock) addWatched(uuid string) *BlockWatch {
	w := &BlockWatch{
		id: uuid,
		ch: make(chan *block.Block, 10),
		stop: func() {
			h.mux.Lock()
			defer h.mux.Unlock()
			delete(h.watchers, uuid)
		},
	}
	h.mux.Lock()
	defer h.mux.Unlock()
	h.watchers[uuid] = w
	return w
}
func (h *HackBlock) StartWatch() {
	h.watch()
}

func (h *HackBlock) SaveBlock(blk *block.Block) error {
	f := filepath.Join(h.path, fmt.Sprintf("block_%d.bin", blk.Header().Number()))
	file, err := os.OpenFile(f, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	if err = blk.EncodeRLP(file); err != nil {
		return err
	}
	return nil
}

func (h *HackBlock) GetBlock(height uint64) (*block.Block, error) {
	if blk, ok := h.cache.Get(height); ok {
		return blk.(*block.Block), nil
	}
	blk, err := h.getBlock(fmt.Sprintf("block_%d.bin", height))
	if err != nil {
		return nil, err
	}
	h.cache.Add(height, blk)
	return blk, nil
}

func (h *HackBlock) getBlock(name string) (*block.Block, error) {
	f := filepath.Join(h.path, name)
	data, err := os.ReadFile(f)
	if err != nil {
		return nil, err
	}

	stream := rlp.NewStream(bytes.NewReader(data), 0)
	blk := new(block.Block)
	if err = blk.DecodeRLP(stream); err != nil {
		return nil, err
	}
	return blk, nil
}

func (h *HackBlock) watch() error {
	watch, err := fsnotify.NewWatcher()
	err = watch.Add(h.path)
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case ev := <-h.eventCh:
				block, err := h.getBlock(ev.Name)
				if err != nil {
					log.Println("get block failed", err)
					continue
				}
				h.mux.Lock()
				for _, w := range h.watchers {
					w.Add(block)
				}
				h.mux.Unlock()
			}
		}
	}()
	go func() {
		for {
			select {
			case ev := <-watch.Events:
				if ev.Op&fsnotify.Create == fsnotify.Create {
					h.eventCh <- ev
				}
			case err := <-watch.Errors:
				{
					log.Println("error : ", err)
					return
				}
			}
		}
	}()
	return nil
}

func (h *HackBlock) WatchNewBlock() *BlockWatch {
	uid := uuid.New()
	return h.addWatched(uid)
}
