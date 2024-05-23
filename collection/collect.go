package collection

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"
	"github.com/xueqianLu/txpress/clientpool"
	"github.com/xueqianLu/txpress/config"
	"math/big"
	"sync"
	"sync/atomic"
	"time"
)

type Info struct {
	sendTime  time.Time
	blockinfo *types.Block
}

type Collect struct {
	cfg        *config.Config
	startBlock int64
	endblock   int64
	latest     atomic.Value
	txs        sync.Map // save txhash => types.Transaction
	txsinfo    sync.Map // save txhash => info
	totalsend  int32
}

func NewCollect(cfg *config.Config) *Collect {
	return &Collect{cfg: cfg}
}

func (c *Collect) TxCount() int32 {
	return c.totalsend
}

func (c *Collect) SetBeginBlock(start int64) {
	c.startBlock = start
}

func (c *Collect) latestTx() *types.Transaction {
	v := c.latest.Load()
	return v.(*types.Transaction)
}

func (c *Collect) SetLatestTx(tx *types.Transaction) {
	c.latest.Store(tx)
}

func (c *Collect) SetSendTime(tx *types.Transaction, tm time.Time) {
	c.txs.Store(tx.Hash(), tx)
	c.txsinfo.Store(tx.Hash(), &Info{sendTime: tm})
	atomic.AddInt32(&c.totalsend, 1)
}

func (c *Collect) SetBlockInfo(txhash common.Hash, block *types.Block) bool {
	i, exist := c.txsinfo.Load(txhash)
	if !exist {
		return false
	}
	info := i.(*Info)
	info.blockinfo = block
	c.txsinfo.Store(txhash, info)
	return true
}

func (c *Collect) getBlockCostTime(blocks []*types.Block, index int) uint64 {
	block := blocks[index]
	if index == (len(blocks) - 1) {
		client := clientpool.GetClient()
		next := new(big.Int).Add(block.Number(), big.NewInt(1))
		nextblock, err := client.BlockByNumber(context.Background(), next)
		if err != nil {
			log.WithFields(log.Fields{
				"err":    err,
				"number": next,
			}).Error("calc block cost time failed")
			return 0
		}
		return nextblock.Time() - block.Time()
	} else {
		nextblock := blocks[index+1]
		return nextblock.Time() - block.Time()
	}
}

func (c *Collect) Run() {
	client := clientpool.GetClient()
	// 1. check latest tx receipt.
	duration := time.Second * 5
	if c.totalsend <= 1000 {
		duration = time.Second * 5
	} else if c.totalsend > 1000 {
		duration = time.Second * 10 * time.Duration(c.totalsend/1000)
	}
	tm := time.NewTicker(duration)
	defer tm.Stop()
	bwait := true
	var waitLatest = sync.OnceFunc(func() { log.Info("wait latest transaction receipt") })
	for bwait {
		select {
		case <-tm.C:
			c.endblock = c.startBlock + 10
			bwait = false
		default:
			if r, err := client.TransactionReceipt(context.Background(), c.latestTx().Hash()); err == nil {
				c.endblock = r.BlockNumber.Int64()
				bwait = false
			} else {
				waitLatest()
				time.Sleep(time.Second)
			}
		}
	}
	// 2. get all block info.
	blocks := make([]*types.Block, 0)
	bstart := true
	txsinfo := make(map[common.Hash]*types.Block, 1000000)
	find := int32(0)
	extblock := 0

	//realStart := c.startBlock

	log.Infof("collect blocks txfind %d, start block is %d, endblock is %d\n", find, c.startBlock, c.endblock)
	for i := c.startBlock; i <= c.endblock; {
		block, err := client.BlockByNumber(context.Background(), big.NewInt(i))
		if err != nil {
			log.Error("get block by number failed", "err", err, "block number", i)
			time.Sleep(time.Second)
			continue
		}
		i++

		txs := block.Transactions()
		if len(txs) == 0 && bstart {
			continue
		}
		if len(txs) > 0 {
			bstart = false
		}
		blocks = append(blocks, block)
		for _, tx := range txs {
			if c.SetBlockInfo(tx.Hash(), block) {
				find++
				txsinfo[tx.Hash()] = block
			}
		}
		if find >= c.totalsend {
			c.endblock = i
			break
		}
		if i >= c.endblock {
			c.endblock += 1
			extblock += 1
		}
		if extblock > 10 {
			break
		}
	}
	c.startBlock = blocks[0].Number().Int64()
	c.endblock = blocks[len(blocks)-1].Number().Int64()
	log.Infof("collect blocks txfind %d, start block is %d, endblock is %d\n", find, c.startBlock, c.endblock)

	// wait generate next block.
	for {
		height, _ := client.BlockNumber(context.Background())
		if int64(height) > c.endblock {
			break
		}
		log.Infof("wait next block %d generate", c.endblock+1)
		time.Sleep(time.Second)
	}

	for i, b := range blocks {
		costtm := c.getBlockCostTime(blocks, i)
		log.Infof("block %d have tx %d, blocktime %d, cost time %ds\n", b.Number().Int64(), len(b.Transactions()), b.Time(), costtm)
	}
	// 3. calc tps
	totalcost := uint64(0)
	totalcount := 0

	if len(blocks) <= 2 {
		for i, b := range blocks {
			costtm := c.getBlockCostTime(blocks, i)
			totalcost += costtm
			totalcount += len(b.Transactions())
		}
	} else {
		newblocks := make([]*types.Block, len(blocks)-2)
		copy(newblocks, blocks[1:len(blocks)-1])
		for i, b := range blocks {
			costtm := c.getBlockCostTime(blocks, i)
			totalcost += costtm
			totalcount += len(b.Transactions())
		}
	}
	log.Infof("total tx %d and cost %d, tps is %d\n", totalcount, totalcost, uint64(totalcount)/totalcost)
}
