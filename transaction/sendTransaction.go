package transaction

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
	"github.com/xueqianLu/txpress/config"
	"math/big"
	"sync/atomic"
	"time"
)

var (
	tranTxCh                                                                       chan *types.Transaction
	clientMap                                                                      map[int]*ethclient.Client
	bathHandleSendCount, onChainTxCount, allSignedTxCount, totalSendTime, forCount int64
	endTran                                                                        common.Hash
)

func init() {
	tranTxCh = make(chan *types.Transaction, 1000000)
	clientMap = make(map[int]*ethclient.Client)
}

// SendTransactions 发送的签名交易
func SendTransactions(client *ethclient.Client, sleepTime int) {
	var beforeSendTxTime time.Time
	var err error
	for {
		select {
		case signedTx := <-tranTxCh:
			if sleepTime != 0 {
				time.Sleep(time.Millisecond * time.Duration(sleepTime))
			}
			sendTranStartTime := time.Now()
			if beforeSendTxTime.UnixMilli() > 0 && time.Since(beforeSendTxTime).Milliseconds() < 1 {
				time.Sleep(time.Millisecond * time.Duration(1))
			}
			err = client.SendTransaction(context.Background(), signedTx)
			beforeSendTxTime = time.Now()
			atomic.AddInt64(&bathHandleSendCount, 1)
			if bathHandleSendCount == allSignedTxCount {
				endTran = signedTx.Hash()
			}
			if err != nil {
				log.Errorf("SendTranErr: %s", err)
				err := client.SendTransaction(context.Background(), signedTx)
				if err != nil {
					log.Errorf("Send tran twice err: %s", err)
					time.Sleep(time.Second * 1)
					continue
				}
			}
			sinceTime := time.Since(sendTranStartTime).Milliseconds()
			atomic.AddInt64(&totalSendTime, sinceTime)
			log.Infof("Send transaction time: %d ms,tranHash: %s", sinceTime, signedTx.Hash())
		}
	}
}

// BatchSendTran 处理批量发送的签名交易
func BatchSendTran(tranArr []*types.Transaction, rpcNodes []string, cfg *config.Config) error {
	rpcCount := len(cfg.RpcNode)
	for i := 0; i < rpcCount; i++ {
		client, err := ethclient.Dial(rpcNodes[i])
		if err != nil {
			log.Errorf("Connect rpc node err: %s,url: %s", err, rpcNodes[i])
			continue
		}
		clientMap[i] = client
	}
	startSendTxBlockNumber, err := clientMap[0].BlockNumber(context.Background())
	if err != nil {
		return err
	}
	for i := 0; i < cfg.SendRoutine; i++ {
		index := i % len(clientMap)
		go SendTransactions(clientMap[index], cfg.SleepTime)
	}
	allSignedTxCount = int64(len(tranArr))
	for i := 0; i < int(allSignedTxCount); i++ {
		tranTxCh <- tranArr[i]
	}
	for {
		if bathHandleSendCount == allSignedTxCount {
			log.Infof("Send tran count: %d", bathHandleSendCount)
			break
		}
	}
	var endBlockNumber *big.Int
	for {
		if forCount == 5 {
			getEndBlockNumber, _ := clientMap[0].BlockNumber(context.Background())
			endBlockNumber = big.NewInt(int64(getEndBlockNumber))
			break
		}
		receipt, err := clientMap[0].TransactionReceipt(context.Background(), endTran)
		if receipt != nil {
			endBlockNumber = receipt.BlockNumber
			break
		}
		if err != nil {
			forCount++
			time.Sleep(time.Second * time.Duration(cfg.AfterSendTranSleepTime))
			log.Infof("Wait end tx receipt.......")
			continue
		}
	}
	StatisticsTxRes(endBlockNumber, big.NewInt(int64(startSendTxBlockNumber)), cfg)
	return nil
}

func StatisticsTxRes(maxBlockNum *big.Int, minBlockNum *big.Int, cfg *config.Config) {
	client := clientMap[0]
	defer closeClientMap(clientMap)
	if maxBlockNum.Int64()-minBlockNum.Int64() > 1 {
		maxBlockNum.Sub(maxBlockNum, big.NewInt(1))
		minBlockNum.Add(minBlockNum, big.NewInt(3))
	}
	log.Infof("Total send time: %d s", totalSendTime/1000)
	log.Infof("maxBlockNum: %d,minBlockNum: %d", maxBlockNum, minBlockNum)
	maxBlockNumInt64 := maxBlockNum.Int64()
	minBlockNumInt64 := minBlockNum.Int64()
	blockNumArr := make([]int, 0)
	blockNumTimeArr := make([]uint64, 0)
	if maxBlockNumInt64 >= minBlockNumInt64 {
		for i := minBlockNumInt64; i <= maxBlockNumInt64; i++ {
			block, err := client.BlockByNumber(context.Background(), big.NewInt(i))
			if err != nil {
				log.Errorf("Statistics block tx count error %s, blockNum:%d", err, i)
				block, err = client.BlockByNumber(context.Background(), big.NewInt(i))
				if err != nil {
					log.Errorf("Statistics block tx count twice error:%s , blockNum:%d", err, i)
					continue
				}
			}
			txCount := len(block.Transactions())
			blockNumTimeArr = append(blockNumTimeArr, block.Time())
			blockNumArr = append(blockNumArr, txCount)
			atomic.AddInt64(&onChainTxCount, int64(txCount))
		}
	} else {
		atomic.AddInt64(&onChainTxCount, int64(cfg.Count))
	}
	maxBlock, err := client.BlockByNumber(context.Background(), maxBlockNum)
	if err != nil {
		log.Errorf("BlockByNumber maxBlock error: %s", err)
		maxBlock, err = client.BlockByNumber(context.Background(), maxBlockNum)
		if err != nil {
			log.Errorf("BlockByNumber maxBlock error: %s", err)
			return
		}
	}
	minBlock, err := client.BlockByNumber(context.Background(), minBlockNum)
	if err != nil {
		log.Errorf("BlockByNumber minBlock error: %s", err)
		minBlock, err = client.BlockByNumber(context.Background(), minBlockNum)
		if err != nil {
			log.Errorf("BlockByNumber minBlock error: %s", err)
			return
		}
	}
	blockNumTimeDifferArr := make([]uint64, 0)
	for i := 0; i < len(blockNumTimeArr); i++ {
		if i == 0 {
			blockNumTimeDifferArr = append(blockNumTimeDifferArr, blockNumTimeArr[i])
		} else {
			blockNumTimeDifferArr = append(blockNumTimeDifferArr, blockNumTimeArr[i]-blockNumTimeArr[i-1])
		}
	}
	totalTxChain := maxBlock.Time() - minBlock.Time()
	log.Infof("Wait last tx time: %d s", forCount*int64(cfg.AfterSendTranSleepTime))
	log.Infof("maxBlockTime: %d ,minBlockTime: %d ,totalTxChain: %d s", maxBlock.Time(), minBlock.Time(), int64(totalTxChain))
	log.Info("All block tx num:", blockNumArr)
	log.Info("All block differ time:", blockNumTimeDifferArr)
	log.Infof("Tx onChain onChainTxCount:%d,  tps %d/%d tx/s -> %d tx/s", onChainTxCount, onChainTxCount, int64(totalTxChain), onChainTxCount/int64(totalTxChain))
}

func closeClientMap(clients map[int]*ethclient.Client) {
	for _, v := range clients {
		v.Close()
	}
}
