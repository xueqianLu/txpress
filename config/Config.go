package config

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"math/big"
)

type Config struct {
	AccountFile           string            `json:"account_file"`
	ReceiveAddr           common.Address    `json:"receive_addr"`
	Count                 int               `json:"tx_count"`
	Type                  int               `json:"tx_type"` // 0: normal tx 1: token tx
	Erc20Contract         common.Address    `json:"erc20_contract"`
	RpcNode               []string          `json:"rpc_node"`
	SendRoutine           int               `json:"send_routine_count"`
	Amount                *big.Int          `json:"amount"`
	ChainId               int64             `json:"chain_id"`
	InitAccountPrv        *ecdsa.PrivateKey `json:"main_account_priv"`
	BatchTransferContract common.Address    `json:"batch_transfer_contract"`
	SendSpeed             int               `json:"speed"`

	SleepTime              int `json:"sleep_time"`
	AfterSendTranSleepTime int `json:"afterSendTranSleepTime"`
}

// MarshalJSON marshals as JSON.
func (h Config) MarshalJSON() ([]byte, error) {
	type IConfig struct {
		AccountFile           string   `json:"account_file"`
		ReceiveAddr           string   `json:"receive_addr"`
		SendSpeed             int      `json:"speed"`
		Count                 int      `json:"tx_count"`
		Type                  int      `json:"tx_type"`
		Erc20Contract         string   `json:"erc20_contract"`
		RpcNode               []string `json:"rpc_node"`
		SendRoutine           int      `json:"send_routine_count"`
		Amount                int64    `json:"amount"`
		ChainId               int64    `json:"chain_id"`
		InitAccountPrv        string   `json:"main_account_priv"`
		BatchTransferContract string   `json:"batch_transfer_contract"`

		SleepTime              int `json:"sleep_time"`
		AfterSendTranSleepTime int `json:"afterSendTranSleepTime"`
	}
	var enc IConfig
	enc.AccountFile = h.AccountFile
	enc.ReceiveAddr = h.ReceiveAddr.String()
	enc.Count = h.Count
	enc.SendSpeed = h.SendSpeed
	enc.Type = h.Type
	enc.Erc20Contract = h.Erc20Contract.String()
	enc.RpcNode = h.RpcNode
	enc.SendRoutine = h.SendRoutine
	enc.Amount = h.Amount.Int64()
	enc.ChainId = h.ChainId
	enc.InitAccountPrv = hexutil.Encode(crypto.FromECDSA(h.InitAccountPrv))
	enc.BatchTransferContract = h.BatchTransferContract.String()
	enc.SleepTime = h.SleepTime
	enc.AfterSendTranSleepTime = h.AfterSendTranSleepTime
	return json.Marshal(&enc)
}

// UnmarshalJSON unmarshals from JSON.
func (h *Config) UnmarshalJSON(input []byte) error {
	type IConfig struct {
		AccountFile           string   `json:"account_file"`
		ReceiveAddr           string   `json:"receive_addr"`
		SendSpeed             int      `json:"speed"`
		Count                 int      `json:"tx_count"`
		Type                  int      `json:"tx_type"`
		Erc20Contract         string   `json:"erc20_contract"`
		RpcNode               []string `json:"rpc_node"`
		SendRoutine           int      `json:"send_routine_count"`
		Amount                int64    `json:"amount"`
		ChainId               int64    `json:"chain_id"`
		InitAccountPrv        string   `json:"main_account_priv"`
		BatchTransferContract string   `json:"batch_transfer_contract"`

		SleepTime              int `json:"sleep_time"`
		AfterSendTranSleepTime int `json:"afterSendTranSleepTime"`
	}
	var dec IConfig
	var err error
	if err = json.Unmarshal(input, &dec); err != nil {
		return err
	}
	h.AccountFile = dec.AccountFile
	h.ReceiveAddr = common.HexToAddress(dec.ReceiveAddr)
	h.Count = dec.Count
	h.SendSpeed = dec.SendSpeed
	h.Type = dec.Type
	h.Erc20Contract = common.HexToAddress(dec.Erc20Contract)
	h.RpcNode = dec.RpcNode
	h.SendRoutine = dec.SendRoutine
	h.Amount = big.NewInt(dec.Amount)
	h.ChainId = dec.ChainId
	h.InitAccountPrv, err = crypto.HexToECDSA(dec.InitAccountPrv)
	if err != nil {
		return errors.New("invalid init account private key")
	}
	h.BatchTransferContract = common.HexToAddress(dec.BatchTransferContract)
	h.SleepTime = dec.SleepTime
	h.AfterSendTranSleepTime = dec.AfterSendTranSleepTime
	return nil
}

var _cfg *Config = nil

func ParseConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Error("get config failed", "err", err)
		panic(err)
	}
	err = json.Unmarshal(data, &_cfg)
	if err != nil {
		log.Error("unmarshal config failed", "err", err)
		panic(err)
	}
	return _cfg, nil
}

func GetConfig() *Config {
	return _cfg
}
