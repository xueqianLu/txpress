package workflow

import (
	log "github.com/sirupsen/logrus"
	"github.com/xueqianLu/txpress/types"
	"sync"
	"time"
)

type Task struct {
	baseCount int
	batch     int
	interval  time.Duration
}

type Record struct {
	Begin     int
	End       int
	TotalTime int
	TotalTx   int
	Tps       int
}

type Result struct {
	chain    string
	minBlock int
	maxBlock int
}

type Workflow struct {
	chains []types.ChainPlugin
	conf   types.RunConfig
	quit   chan struct{}
}

func NewWorkFlow(chains []types.ChainPlugin, conf types.RunConfig) *Workflow {
	return &Workflow{
		chains: chains,
		conf:   conf,
		quit:   make(chan struct{}),
	}
}
func (w *Workflow) Start() {
	resultCh := make(chan Result, 1)

	wg := sync.WaitGroup{}
	taskChs := make([]chan Task, len(w.chains))
	for i, chain := range w.chains {
		taskChs[i] = make(chan Task, 1)
		wg.Add(1)
		go w.loop(chain, taskChs[i], resultCh, &wg)
	}

	baseTxCount := w.conf.BaseCount
	lastTps := 0
	noincrease := 0

	for {
		for _, ch := range taskChs {
			ch <- Task{
				baseCount: baseTxCount,
				batch:     w.conf.Batch,
				interval:  w.conf.Interval,
			}
		}
		// wait task finished.
		var minBlock, maxBlock int
		results := make([]Result, 0)
		for len(results) < len(w.chains) {
			select {
			case result := <-resultCh:
				log.Infof("chain %s finished, min block: %d, max block: %d", result.chain, result.minBlock, result.maxBlock)
				results = append(results, result)
				if minBlock == 0 {
					minBlock = result.minBlock
				}
				if maxBlock == 0 {
					maxBlock = result.maxBlock
				}

				if result.minBlock < minBlock {
					minBlock = result.minBlock
				}
				if result.maxBlock > maxBlock {
					maxBlock = result.maxBlock
				}
			case <-time.After(time.Second * 5):
				log.Infof("wait test finished")
			}
		}
		log.Info("all chain run task finished")
		// calculate tps
		record := w.calculateTps(w.chains[0], minBlock, maxBlock)
		if record.Tps > 0 && record.Tps >= lastTps {
			baseTxCount *= 2
			lastTps = record.Tps
		} else {
			noincrease++
			if noincrease >= 2 {
				break
			}
		}
		log.WithFields(log.Fields{
			"begin":     record.Begin,
			"end":       record.End,
			"totaltime": record.TotalTime,
			"totaltx":   record.TotalTx,
			"tps":       record.Tps,
		}).Info("test one round finished")
	}

	close(w.quit)
	wg.Wait()

}

func (w *Workflow) calculateTps(chain types.ChainPlugin, minBlock, maxBlock int) Record {
	start := int64(0)
	end := int64(0)
	txCount := int64(0)
	for i := minBlock; i <= maxBlock; i++ {
		block, err := chain.GetBlockInfo(int64(i))
		if err != nil {
			log.Errorf("get block info failed: %s", err)
			continue
		}
		if i == minBlock {
			start = block.Timestamp
		}
		if i == maxBlock {
			end = block.Timestamp
		}
		txCount += block.TxCount
	}
	record := Record{
		Begin:     int(minBlock),
		End:       int(maxBlock),
		TotalTime: int(end-start) + chain.SecondPerBlock(),
		TotalTx:   int(txCount),
	}
	if record.TotalTime > 0 {
		record.Tps = int(txCount) / (record.TotalTime)
	}

	return record
}

func (w *Workflow) makeTx(chain types.ChainPlugin, baseCount int, batch int, checkNonce bool) [][]types.ChainTx {
	txs := make([][]types.ChainTx, batch)
	for i := 0; i < batch; i++ {
		if i > 0 {
			checkNonce = false
		}
		mtxs, err := chain.CreateTxs(baseCount, checkNonce)
		if err != nil {
			return nil
		}
		txs[i] = mtxs
	}
	return txs
}

func (w *Workflow) loop(chain types.ChainPlugin, taskCh chan Task, result chan Result, wg *sync.WaitGroup) {

	defer wg.Done()
	first := true

	for {
		select {
		case task := <-taskCh:
			txs := w.makeTx(chain, task.baseCount, task.batch, first)
			_min, _max := w.runTest(chain, txs, task.interval)
			result <- Result{
				chain:    chain.Id(),
				minBlock: _min,
				maxBlock: _max,
			}
			if first {
				first = false
			}

		case <-w.quit:
			return
		}
	}
}

func (w *Workflow) runTest(chain types.ChainPlugin, txs [][]types.ChainTx, interval time.Duration) (int, int) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	var waited = make(map[string]int)
	hashes := make([][]string, 0)
	for _, batch := range txs {
		hashs, err := chain.SendTxs(batch)
		if err != nil {
			log.Errorf("send txs error: %s", err)
		} else {
			log.Infof("send txs success, count: %d", len(hashs))
			hashes = append(hashes, hashs)
			for _, hash := range hashs {
				waited[hash] = 0
			}
		}
		<-ticker.C
	}
	var _min, _max int
	for len(waited) > 0 {
		time.Sleep(time.Second)
		for hash, _ := range waited {
			time.Sleep(time.Millisecond * 10)
			block, err := chain.TxBlock(hash)
			if err != nil {
				continue
			}
			if _min == 0 {
				_min = block
			}

			if block < _min {
				_min = block
			}

			if _max == 0 {
				_max = block
			}

			if block > _max {
				_max = block
			}
			delete(waited, hash)
		}
	}

	return _min, _max
}
