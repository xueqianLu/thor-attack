// Copyright (c) 2018 The VeChainThor developers

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package poa

import (
	"github.com/vechain/thor/thor"
)

// SchedulerSimp to schedule the time when a proposer to produce a block.
// V2 is for post VIP-214 stage.
type SchedulerSimp struct {
	proposer        Proposer
	parentBlockTime uint64
	proposers       []Proposer
	geneTime        uint64
}

var _ Scheduler = (*SchedulerSimp)(nil)

// NewSchedulerSimp create a SchedulerSimp object.
// `addr` is the proposer to be scheduled.
// If `addr` is not listed in `proposers` or not active, an error returned.
func NewSchedulerSimp(
	addr thor.Address,
	proposers []Proposer,
	parentBlockNumber uint32,
	parentBlockTime uint64,
	seed []byte, geneTime uint64) (*SchedulerSimp, error) {

	var (
		proposer Proposer
	)

	for _, p := range proposers {
		if p.Address == addr {
			proposer = p
		}
	}
	log.Info("scheduler simp", "proposer list", proposers, "proposer", proposer)
	return &SchedulerSimp{
		proposer,
		parentBlockTime,
		proposers,
		geneTime,
	}, nil
}

// Schedule to determine time of the proposer to produce a block, according to `nowTime`.
// `newBlockTime` is promised to be >= nowTime and > parentBlockTime
func (s *SchedulerSimp) Schedule(nowTime uint64) (newBlockTime uint64) {
	const T = thor.BlockInterval

	newBlockTime = s.parentBlockTime + T
	if nowTime > newBlockTime {
		// ensure T aligned, and >= nowTime
		newBlockTime += (nowTime - newBlockTime + T - 1) / T * T
	}

	offset := (newBlockTime-s.geneTime)/T - 1
	for i, n := uint64(0), uint64(len(s.proposers)); i < n; i++ {
		index := (i + offset) % n
		if s.proposers[index].Address == s.proposer.Address {
			log.Info("scheduler schedule", "blockTime", newBlockTime+i*T, "proposer", s.proposer.Address, "index", index)
			return newBlockTime + i*T
		}
	}

	// should never happen
	panic("something wrong with proposers list")
}

// IsTheTime returns if the newBlockTime is correct for the proposer.
func (s *SchedulerSimp) IsTheTime(newBlockTime uint64) bool {
	return s.IsScheduled(newBlockTime, s.proposer.Address)
}

// IsScheduled returns if the schedule(proposer, blockTime) is correct.
func (s *SchedulerSimp) IsScheduled(blockTime uint64, proposer thor.Address) bool {
	if s.parentBlockTime >= blockTime {
		// invalid block time
		return false
	}

	T := thor.BlockInterval
	if (blockTime-s.parentBlockTime)%T != 0 {
		// invalid block time
		return false
	}

	index := (blockTime - s.geneTime - T) / T % uint64(len(s.proposers))
	return s.proposers[index].Address == proposer
}

// Updates returns proposers whose status are changed, and the score when new block time is assumed to be newBlockTime.
func (s *SchedulerSimp) Updates(newBlockTime uint64) (updates []Proposer, score uint64) {
	T := thor.BlockInterval

	for i := uint64(0); i < uint64(len(s.proposers)); i++ {
		if s.parentBlockTime+T+i*T >= newBlockTime {
			break
		}
		if s.proposers[i].Address != s.proposer.Address {
			updates = append(updates, Proposer{Address: s.proposers[i].Address, Active: false})
		}
	}

	score = uint64(len(s.proposers) - len(updates))

	if !s.proposer.Active {
		cpy := s.proposer
		cpy.Active = true
		updates = append(updates, cpy)
	}
	log.Info("scheduler updates", "blockTime", newBlockTime, "len(proposers)", len(s.proposers), "score", score)
	return
}
