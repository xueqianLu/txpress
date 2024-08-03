package finalize

import (
	log "github.com/sirupsen/logrus"
	"github.com/xueqianLu/txpress/chains"
	"github.com/xueqianLu/txpress/types"
	"time"
)

type Finalize struct {
	chain types.ChainPlugin
	quit  chan struct{}
}

func NewFinalize(chain types.ChainPlugin) *Finalize {
	return &Finalize{
		chain: chain,
		quit:  make(chan struct{}),
	}
}

func (f *Finalize) Loop() {
	tm := time.NewTicker(time.Minute * 20)
	updateTm := time.NewTicker(time.Second * 10)
	defer tm.Stop()
	defer updateTm.Stop()
	lastfinalized := 0
	blockTime := make(map[int]time.Time)

	for {
		select {
		case <-f.quit:
			return
		case <-updateTm.C:
			finalized, err := f.chain.FinalizedBlock()
			if err != nil {
				continue
			}
			if _, ok := blockTime[finalized]; !ok {
				blockTime[finalized] = time.Now()
			}

		case <-tm.C:
			// 1. get latest block number.
			blk, err := f.chain.LatestBlockInfo()
			if err != nil {
				continue
			}

			// 2. get finalized block.
			finalized, err := f.chain.FinalizedBlock()
			if err != nil {
				continue
			}

			if finalized > lastfinalized {
				if _, exist := blockTime[lastfinalized]; !exist {
					blockTime[lastfinalized] = time.Now()
				}
				if _, exist := blockTime[finalized]; !exist {
					blockTime[finalized] = time.Now()
				}

				// if finalized block changed, calc tps from last finalized block to current finalized block.
				record := chains.CalcTps(f.chain, lastfinalized+1, finalized, blockTime[finalized].Sub(blockTime[lastfinalized]))
				lastfinalized = finalized
				log.WithFields(log.Fields{
					"chain":       f.chain.Id(),
					"tps":         record.Tps,
					"begin":       record.Begin,
					"end":         record.End,
					"txcount":     record.TotalTx,
					"latestBlock": blk.Number,
					"finalized":   finalized,
				}).Info("finalized tps info")
			} else {
				log.WithFields(log.Fields{
					"latestBlock": blk.Number,
					"finalized":   finalized,
				}).Info("finalized info on the chain")
			}
		}
	}
}

func (f *Finalize) Stop() {
	close(f.quit)
}
