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
	tm := time.NewTicker(time.Second * 60)
	defer tm.Stop()
	lastfinalized := 0

	for {
		select {
		case <-f.quit:
			return
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
				// if finalized block changed, calc tps from last finalized block to current finalized block.
				record := chains.CalcTps(f.chain, lastfinalized+1, finalized)
				lastfinalized = finalized
				log.WithFields(log.Fields{
					"chain":   f.chain.Id(),
					"tps":     record.Tps,
					"begin":   record.Begin,
					"end":     record.End,
					"txcount": record.TotalTx,
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
