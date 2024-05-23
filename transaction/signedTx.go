package transaction

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	log "github.com/sirupsen/logrus"
	"github.com/xueqianLu/txpress/config"
	"golang.org/x/crypto/sha3"
	"math/big"
	"sync"
	"time"
)

type TranConfig struct {
	Amount                                 int64
	GoRoutineCount                         int
	ChainId                                int64
	Contract                               bool
	ContractAddr, PrivateKey, ReceivedAddr string
	Nonce                                  *big.Int
}

type Transactor struct {
	config       TranConfig
	signerKey    *ecdsa.PrivateKey
	sender       common.Address
	receivedAddr common.Address
}

func newTransactor(cfg TranConfig) (*Transactor, error) {
	signerKey, err := crypto.HexToECDSA(cfg.PrivateKey)
	if err != nil {
		log.Error("Error crypto HexToECDSA")
		return nil, err
	}
	// through privateKey get account address
	publicKey := signerKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Error("Error casting public key to ECDSA")
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	res := Transactor{
		signerKey:    signerKey,
		receivedAddr: common.HexToAddress(cfg.ReceivedAddr),
		config:       cfg,
		sender:       address,
	}
	return &res, nil
}

// SignedTxArr 获取全部签名数据
func SignedTxArr(sendTxAccountArr [][]string, cfg *config.Config, isContract bool, accountsNonceMap *sync.Map) []*types.Transaction {
	tranArr := make([]*types.Transaction, 0)
	var signedTx *types.Transaction
	for _, rows := range sendTxAccountArr {
		privateKey := rows[1]
		value, ok := accountsNonceMap.Load(privateKey)
		if !ok {
			log.Error("Load nonce map error...........")
			continue
		}
		nonce := big.NewInt(value.(int64))
		tranCfg := TranConfig{
			Amount:         cfg.Amount.Int64(),
			ChainId:        cfg.ChainId,
			Contract:       isContract,
			ContractAddr:   cfg.Erc20Contract.String(),
			PrivateKey:     privateKey,
			ReceivedAddr:   cfg.ReceiveAddr.String(),
			GoRoutineCount: cfg.SendRoutine,
			Nonce:          nonce,
		}
		t, err := newTransactor(tranCfg)
		if t.config.Contract {
			signedTx, err = t.signedContractTx()
		} else {
			signedTx, err = t.signedTx()
		}
		nonce = big.NewInt(1).Add(nonce, big.NewInt(1))
		if err != nil || signedTx == nil {
			log.Errorf("signed tx error %s ", err)
			continue
		}
		tranArr = append(tranArr, signedTx)
	}
	return tranArr
}

// signedContractTx 签名合约代币转账交易
func (t *Transactor) signedContractTx() (*types.Transaction, error) {
	value := big.NewInt(0)
	toAddress := common.HexToAddress(t.config.ReceivedAddr)
	tokenAddress := common.HexToAddress(t.config.ContractAddr)
	// 转账签名方法
	transferFnSignature := []byte("transfer(address,uint256)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	amount, _ := new(big.Int).SetString("10000000000000000000", 10)
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)
	txData := types.LegacyTx{
		Nonce:    t.config.Nonce.Uint64(),
		To:       &tokenAddress,
		Value:    value,
		Gas:      100000,
		GasPrice: big.NewInt(1000000001),
		Data:     data,
	}
	newtx := types.NewTx(&txData)
	signedTx, err := types.SignTx(newtx, types.NewEIP155Signer(big.NewInt(t.config.ChainId)), t.signerKey)
	if err != nil {
		log.Errorf("Signed contract tx error: %s", err)
		return nil, err
	}
	return signedTx, nil
}

// signedTx sign normal transaction
func (t *Transactor) signedTx() (*types.Transaction, error) {
	txData := types.LegacyTx{
		Nonce:    t.config.Nonce.Uint64(),
		To:       &t.receivedAddr,
		Value:    big.NewInt(t.config.Amount),
		Gas:      21000,
		GasPrice: big.NewInt(1000000000),
		Data:     nil,
	}
	newtx := types.NewTx(&txData)
	signedTx, err := types.SignTx(newtx, types.NewEIP155Signer(big.NewInt(t.config.ChainId)), t.signerKey)
	if err != nil {
		log.Errorf("Send tx nonce: %d , From: %s , to: %s , error: %s", t.config.Nonce, crypto.PubkeyToAddress(t.signerKey.PublicKey), t.receivedAddr, err.Error())
		time.Sleep(time.Second)
		return nil, err
	}
	return signedTx, nil
}
