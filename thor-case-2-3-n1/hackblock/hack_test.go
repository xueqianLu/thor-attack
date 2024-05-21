package hackblock

import (
	"crypto/rand"
	"fmt"
	"github.com/vechain/thor/block"
	"github.com/vechain/thor/thor"
	"math/big"
	"testing"
	"time"
)

func TestNewHackBlock(t *testing.T) {
	h := NewHackBlock("./")
	if h == nil {
		t.Error("NewHackBlock failed")
	} else {
		t.Log("NewHackBlock success")
	}

}

func TestHackBlock_WatchBlock(t *testing.T) {
	build := block.Builder{}
	parent, _ := rand.Int(rand.Reader, big.NewInt(100000))
	build.ParentID(thor.BytesToBytes32(parent.Bytes()))
	build.Timestamp(1)
	build.TotalScore(1)
	build.GasLimit(1)
	build.GasUsed(1)
	build.Beneficiary([20]byte{1})
	build.StateRoot([32]byte{1})
	build.ReceiptsRoot([32]byte{1})

	blk := build.Build()
	h := NewHackBlock("./")
	if h == nil {
		t.Error("NewHackBlock failed")
		return
	}
	h.StartWatch()
	watch := h.WatchNewBlock()
	go func() {
		for {
			select {
			case blk := <-watch.Watch():
				fmt.Println("watch new block", blk)
			}
		}
	}()

	if err := h.SaveBlock(blk); err != nil {
		t.Error("SaveBlock failed", err)
	} else {
		t.Log("SaveBlock success")
	}
	time.Sleep(time.Second)
}

func TestHackBlock_SaveBlock(t *testing.T) {
	build := block.Builder{}
	build.ParentID(thor.BytesToBytes32([]byte{1}))
	build.Timestamp(1)
	build.TotalScore(1)
	build.GasLimit(1)
	build.GasUsed(1)
	build.Beneficiary([20]byte{1})
	build.StateRoot([32]byte{1})
	build.ReceiptsRoot([32]byte{1})

	blk := build.Build()
	h := NewHackBlock("./")
	if h == nil {
		t.Error("NewHackBlock failed")
		return
	}

	if err := h.SaveBlock(blk); err != nil {
		t.Error("SaveBlock failed", err)
	} else {
		t.Log("SaveBlock success")
	}
}

func TestHackBlock_GetBlock(t *testing.T) {
	h := NewHackBlock("./")
	if h == nil {
		t.Error("NewHackBlock failed")
		return
	}
	blk, err := h.GetBlock(1)
	if err != nil {
		t.Error("GetBlock failed", err)
	} else {
		t.Log("GetBlock success")
	}
	t.Log("blk", blk)
}
