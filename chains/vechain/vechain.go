package vechain

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"github.com/vechain/thor/thor"
	"github.com/xueqianLu/txpress/types"
	"math/big"
)

var _ types.ChainPlugin = &VeChain{}

type VeChain struct {
	chainTag byte
	client   *Client
	accounts []*AccountInfo
	ctx      context.Context
	rpc      string
	index    int
	config   types.ChainConfig
}

func (e VeChain) CreateTxs(count int, checkNonce bool) ([]types.ChainTx, error) {
	txs := make([]types.ChainTx, 0)

	amount, _ := new(big.Int).SetString(e.config.Amount, 10)
	txcfg := txConfig{
		receiver: thor.BytesToAddress(common.Hex2Bytes(e.config.Receiver)),
		amount:   amount,
		chainTag: e.chainTag,
	}
	for i := 0; i < count; i++ {
		acc := e.accounts[i%len(e.accounts)]
		tx := acc.MakeTx(txcfg, acc.Nonce)
		txs = append(txs, VeTx{tx})

		acc.Nonce++
	}
	return txs, nil
}

func (e VeChain) SendTxs(txs []types.ChainTx) ([]string, error) {
	hashes := make([]string, 0)
	for _, tx := range txs {
		vtx := tx.(VeTx)
		err := e.client.SendTransaction(e.ctx, vtx.Transaction)
		if err != nil {
			log.WithFields(log.Fields{
				"chain": e.config.Name,
				"rpc":   e.rpc,
				"index": e.index,
				"txid":  vtx.Transaction.ID().String(),
				"err":   err,
			}).Error("send tx failed")
			continue
		} else {
			//log.WithFields(log.Fields{
			//	"chain": e.config.Name,
			//	"rpc":   e.rpc,
			//	"index": e.index,
			//	"tx":    etx.Transaction.Hash().String(),
			//}).Info("send tx success")
		}
		hashes = append(hashes, vtx.Transaction.ID().String())

	}
	return hashes, nil
}

func (e VeChain) TxReceipt(hash string) error {
	_, err := e.client.TransactionReceipt(e.ctx, hash)
	if err != nil {
		log.WithFields(log.Fields{
			"chain": e.config.Name,
			"rpc":   e.rpc,
			"index": e.index,
			"err":   err,
		}).Error("get tx receipt failed")
		return err
	}
	return err
}

func (e VeChain) TxBlock(hash string) (int, error) {
	receipt, err := e.client.TransactionReceipt(e.ctx, hash)
	if err != nil {
		log.WithFields(log.Fields{
			"chain": e.config.Name,
			"rpc":   e.rpc,
			"index": e.index,
			"err":   err,
			"hash":  hash,
		}).Error("get tx receipt failed")
		return 0, err
	}
	return int(receipt.Meta.BlockNumber), nil
}

func (e VeChain) GetBlockInfo(number int64) (types.BlockInfo, error) {
	blk := new(big.Int).SetInt64(number)
	info, err := e.client.BlockByNumber(e.ctx, blk)
	if err != nil {
		log.WithFields(log.Fields{
			"chain": e.config.Name,
			"rpc":   e.rpc,
			"index": e.index,
			"err":   err,
		}).Error("get block info failed")
		return types.BlockInfo{}, err
	}
	return types.BlockInfo{
		Number:    info.Number,
		Timestamp: info.Timestamp,
		TxCount:   int64(len(info.Txs)),
	}, nil
}

func (e VeChain) Id() string {
	return fmt.Sprintf("%s-%d", e.config.Name, e.index)
}

func (e VeChain) SecondPerBlock() int {
	return 10
}

var (
	totalAccounts []*AccountInfo
)

func NewVeChain(rpc string, index int, config types.ChainConfig) (types.ChainPlugin, error) {
	ctx := context.TODO()
	client := NewVechainClient(ctx, rpc)
	chainTag, err := client.ChainTag()
	if err != nil {
		return nil, err
	}
	total := len(config.Rpcs)

	if totalAccounts == nil {
		totalAccounts = GetAccountJson(config.Accounts)
	}
	chainAccounts := make([]*AccountInfo, 0)
	for i := 0; i < len(totalAccounts); i++ {
		if i%total == index {
			chainAccounts = append(chainAccounts, totalAccounts[i])
		}
	}

	return &VeChain{
		ctx:      ctx,
		rpc:      rpc,
		index:    index,
		accounts: chainAccounts,
		client:   client,
		chainTag: chainTag,
		config:   config,
	}, nil
}
