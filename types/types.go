package types

import "time"

type ChainTx interface {
	IsChainTx() bool
}

type BlockInfo struct {
	Timestamp int64
	Number    int64
	TxCount   int64
}

type ChainPlugin interface {
	CreateTxs(count int, checkNonce bool) ([]ChainTx, error)
	SendTxs(txs []ChainTx) ([]string, error)
	TxReceipt(hash string) error
	TxBlock(hash string) (int, error)
	GetBlockInfo(number int64) (BlockInfo, error)
	Id() string
}

type RunConfig struct {
	BaseCount int
	Interval  time.Duration
	Batch     int
}

type ChainConfig struct {
	Rpcs      []string `json:"rpc-nodes"`
	Name      string   `json:"chain-name"`
	BaseCount int      `json:"base-count"`
	Interval  int      `json:"interval"`
	Batch     int      `json:"batch"`
	Receiver  string   `json:"receiver"`
	Amount    string   `json:"amount"`
	Accounts  string   `json:"accounts"`
}
