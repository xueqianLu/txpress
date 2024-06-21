package main

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
	"github.com/xueqianLu/txpress/clientpool"
	"github.com/xueqianLu/txpress/collection"
	"github.com/xueqianLu/txpress/config"
	"github.com/xueqianLu/txpress/tool"
	"time"
)

func maketxSequences(cfg *config.Config, accounts []*tool.Account) [][]*types.Transaction {
	type param struct {
		account *tool.Account
		nonce   uint64
		count   int
	}
	countPerAccount := (cfg.Count + len(accounts)) / len(accounts)
	taskpool := make(chan interface{}, len(accounts))
	totalCount := 0
	for i := 0; i < len(accounts); i++ {
		p := param{
			account: accounts[i],
			nonce:   accounts[i].Nonce,
		}
		if i == len(accounts)-1 {
			p.count = cfg.Count - totalCount
		} else {
			p.count = countPerAccount
			totalCount += countPerAccount
		}
		taskpool <- p
	}
	output := make(chan []*types.Transaction, len(accounts))
	tasks := tool.NewTasks(10, func(task interface{}) {
		var tx *types.Transaction
		p := task.(param)
		userTxs := make([]*types.Transaction, 0)
		for i := 0; i < p.count; i++ {
			if cfg.Type == 1 {
				tx = p.account.MakeTokenTx(cfg, p.nonce+uint64(i))
			} else {
				tx = p.account.MakeNormalTx(cfg, p.nonce+uint64(i))
			}
			signedtx, err := p.account.SignTx(cfg, tx)
			if err != nil {
				log.Error("sign tx failed", "tx", tx, "err", err)
			} else {
				log.Debugf("account (%s) sign tx (%s)", p.account.Address, signedtx.Hash())
				userTxs = append(userTxs, signedtx)
			}
		}
		output <- userTxs
	}, taskpool)
	tasks.Run()
	close(taskpool)
	tasks.Done()
	total := 0
	txs := make([][]*types.Transaction, 0, len(accounts))
	for total < cfg.Count {
		signedtxs := <-output
		total += len(signedtxs)
		txs = append(txs, signedtxs)
	}
	close(output)
	return txs
}

func maketx(cfg *config.Config, accounts []*tool.Account) []*types.Transaction {
	type param struct {
		account *tool.Account
		nonce   uint64
	}
	taskpool := make(chan interface{}, cfg.Count)
	for i := 0; i < cfg.Count; i++ {
		idx := i % len(accounts)
		loop := uint64(i / len(accounts))
		p := param{
			account: accounts[idx],
			nonce:   accounts[idx].Nonce + loop,
		}
		taskpool <- p
	}
	txs := make([]*types.Transaction, 0, cfg.Count)
	output := make(chan interface{}, cfg.Count)
	tasks := tool.NewTasks(10, func(task interface{}) {
		var tx *types.Transaction
		p := task.(param)
		if cfg.Type == 1 {
			tx = p.account.MakeTokenTx(cfg, p.nonce)
		} else {
			tx = p.account.MakeNormalTx(cfg, p.nonce)
		}
		signedtx, err := p.account.SignTx(cfg, tx)
		if err != nil {
			log.Error("sign tx failed", "tx", tx, "err", err)
		} else {
			log.Debugf("account (%s) sign tx (%s)", p.account.Address, signedtx.Hash())
			output <- signedtx
		}
	}, taskpool)
	tasks.Run()
	close(taskpool)
	tasks.Done()
	total := len(output)
	for len(txs) < total {
		signedtx := <-output
		txs = append(txs, signedtx.(*types.Transaction))
	}
	close(output)
	return txs
}

func sendWithTimeout(client *ethclient.Client, ctx context.Context, timeout time.Duration, tx *types.Transaction) error {
	tm := time.NewTicker(timeout)
	defer tm.Stop()
	res := make(chan struct{}, 1)
	var err error
	go func() {
		err = client.SendTransaction(ctx, tx)
		res <- struct{}{}
	}()
	select {
	case <-res:
		return err
	case <-tm.C:
		return nil
	}
}

func sendTxSequence(cfg *config.Config, txs [][]*types.Transaction, collect *collection.Collect) {
	taskpool := make(chan interface{}, len(txs))
	for _, accountTxs := range txs {
		taskpool <- accountTxs
	}

	task := tool.NewTasksWithSpeed(cfg.SendRoutine, func(task interface{}) {
		ctx := context.Background()
		client := clientpool.GetClient()
		accountTxs := task.([]*types.Transaction)
		for _, tx := range accountTxs {
			s1 := time.Now()
			err := sendWithTimeout(client, ctx, time.Second*3, tx)
			if err != nil {
				log.Errorf("send tx (%s) failed err %s", tx.Hash(), err.Error())
			} else {
				collect.SetSendTime(tx, time.Now())
			}
			s2 := time.Now()
			log.Debugf("send tx cost tm %vms\n", s2.Sub(s1).Milliseconds())
		}
	}, taskpool, cfg.SendSpeed)
	task.Run()
	close(taskpool)
	task.Done()
}

func sendTx(cfg *config.Config, txs []*types.Transaction, collect *collection.Collect) {
	taskpool := make(chan interface{}, len(txs))
	for _, tx := range txs {
		taskpool <- tx
	}

	task := tool.NewTasksWithSpeed(cfg.SendRoutine, func(task interface{}) {
		ctx := context.Background()
		client := clientpool.GetClient()
		tx := task.(*types.Transaction)
		s1 := time.Now()
		err := sendWithTimeout(client, ctx, time.Second*3, tx)
		if err != nil {
			log.Errorf("send tx (%s) failed err %s", tx.Hash(), err.Error())
		} else {
			collect.SetSendTime(tx, time.Now())
		}
		s2 := time.Now()
		log.Debugf("send tx cost tm %vms\n", s2.Sub(s1).Milliseconds())
	}, taskpool, cfg.SendSpeed)
	task.Run()
	close(taskpool)
	task.Done()
}

func initcollect(cfg *config.Config) *collection.Collect {
	collect := collection.NewCollect(cfg)
	client := clientpool.GetClient()
	blocknumber, _ := client.BlockNumber(context.Background())
	collect.SetBeginBlock(int64(blocknumber))

	return collect
}

func randomReceive() common.Address {
	pk, _ := crypto.GenerateKey()
	addr := crypto.PubkeyToAddress(pk.PublicKey)
	return addr
}

func start(cfg *config.Config, accounts []*tool.Account) {
	var empty = common.Address{}
	if cfg.ReceiveAddr == empty {
		cfg.ReceiveAddr = randomReceive()
	}
	log.Infof("start make tx to receive (%s) count %v\n", cfg.ReceiveAddr.String(), cfg.Count)
	txs := maketxSequences(cfg, accounts)
	//txs := maketx(cfg, accounts)
	log.Info("make tx finished")
	collect := initcollect(cfg)
	log.Info("start send tx")
	sendTxSequence(cfg, txs, collect)
	//sendTx(cfg, txs, collect)
	log.Infof("send tx succeed and total %v.\n", collect.TxCount())
	collect.SetLatestTx(txs[len(txs)-1][len(txs[len(txs)-1])-1])
	collect.Run()
	log.Info("test finished")
}
