package tool

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
	"github.com/xueqianLu/txpress/clientpool"
	"github.com/xueqianLu/txpress/config"
	"github.com/xueqianLu/txpress/contract"
	"io/fs"
	"math/big"
	"os"
	"strings"
)

type Account struct {
	Address common.Address    `json:"address"`
	Private string            `json:"private"`
	Nonce   uint64            `json:"nonce"`
	PK      *ecdsa.PrivateKey `json:"-"`
}

func AccountFromPrivk(pk *ecdsa.PrivateKey) *Account {
	//var err error
	a := &Account{
		PK: pk,
	}
	a.Address = crypto.PubkeyToAddress(a.PK.PublicKey)
	a.Private = hexutil.Encode(crypto.FromECDSA(pk))
	client := clientpool.GetClient()
	a.Nonce, _ = client.NonceAt(context.Background(), a.Address, nil)
	return a
}

func (acc *Account) SendInitTokenTx(client *ethclient.Client, cfg *config.Config, count int64, peramount *big.Int) (*types.Transaction, error) {
	unit := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	totalCount := new(big.Int).Mul(peramount, big.NewInt(count))
	tokenAmount := new(big.Int).Mul(unit, totalCount)

	tokenContract, err := contract.NewToken(cfg.Erc20Contract, client)
	auth, err := bind.NewKeyedTransactorWithChainID(acc.PK, new(big.Int).SetInt64(cfg.ChainId))
	if err != nil {
		log.Error("NewKeyedTransactorWithChainID Error:", err)
	}
	auth.GasLimit = 90000000
	auth.GasPrice, _ = new(big.Int).SetString("1000000000", 10)

	nonce, err := client.NonceAt(context.Background(), acc.Address, nil)
	auth.Nonce = big.NewInt(int64(nonce))

	return tokenContract.Transfer(auth, cfg.BatchTransferContract, tokenAmount)

}

func (acc *Account) MakeInitTx(cfg *config.Config, nonce uint64, count int64, peramount *big.Int) *types.Transaction {
	unit := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	totalCount := new(big.Int).Mul(peramount, big.NewInt(count))
	txData := &types.LegacyTx{
		Nonce:    nonce,
		To:       &cfg.BatchTransferContract,
		Value:    new(big.Int).Mul(unit, totalCount),
		Gas:      300000,
		GasPrice: big.NewInt(1000000000),
		Data:     nil,
	}
	tx, err := acc.SignTx(cfg, types.NewTx(txData))
	if err != nil {
		log.Error("sign init tx failed", "err", err)
		return nil
	}
	return tx
}

func (acc *Account) MakeNormalTx(cfg *config.Config, nonce uint64) *types.Transaction {
	txData := &types.LegacyTx{
		Nonce:    nonce,
		To:       &cfg.ReceiveAddr,
		Value:    cfg.Amount,
		Gas:      300000,
		GasPrice: big.NewInt(1000000000),
		Data:     nil,
	}
	return types.NewTx(txData)
}

func (acc *Account) MakeTokenTx(cfg *config.Config, nonce uint64) *types.Transaction {
	parsed, err := abi.JSON(strings.NewReader(contract.TokenABI))
	if err != nil {
		log.Println("contract json failed.")
		log.Fatal(err)
	}

	data, _ := parsed.Pack("transfer", cfg.ReceiveAddr, cfg.Amount)

	txData := &types.LegacyTx{
		Nonce:    nonce,
		To:       &cfg.Erc20Contract,
		Value:    big.NewInt(0),
		Gas:      90000000,
		GasPrice: big.NewInt(1000000000),
		Data:     data,
	}
	return types.NewTx(txData)
}

func (acc *Account) SignTx(cfg *config.Config, tx *types.Transaction) (*types.Transaction, error) {
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(cfg.ChainId)), acc.PK)
	if err != nil {
		log.Error("sign tx failed", "err", err)
		return nil, err
	}
	return signedTx, nil
}

func createAccounts(count int) []*Account {
	accs := make([]*Account, 0, count)
	for i := 0; i < count; i++ {
		pk, _ := crypto.GenerateKey()
		addr := crypto.PubkeyToAddress(pk.PublicKey)
		private := hexutil.Encode(crypto.FromECDSA(pk))
		accs = append(accs, &Account{Address: addr, Private: private, PK: pk})
	}
	return accs
}

func CreateAccounts(cfg *config.Config, count int) []*Account {
	accs := createAccounts(count)
	d, _ := json.MarshalIndent(accs, "", "    ")
	err := os.WriteFile(cfg.AccountFile, d, fs.ModePerm)
	if err != nil {
		log.Error("write account failed", "err", err)
	}
	return accs
}

func GetAccountJson(cfg *config.Config) []*Account {
	data, err := os.ReadFile(cfg.AccountFile)
	if err != nil || len(data) == 0 {
		return []*Account{}
	} else {
		accs := make([]*Account, 0)
		err = json.Unmarshal(data, &accs)
		if err != nil {
			log.Error("unmarshal account failed", "err", err)
		}
		for _, acc := range accs {
			acc.PK, err = crypto.HexToECDSA(acc.Private[2:])
			if err != nil {
				log.Error("hex to ecdsa failed", "err", err)
			}
		}
		log.Info("get accounts from json", "len", len(accs))
		return accs
	}
}

func CheckAccountNonce(client *ethclient.Client, account *Account) {

	nonce, err := client.NonceAt(context.Background(), account.Address, nil)
	if err != nil {
		// ignore err
	} else {
		log.Debug("check account nonce", "account", account.Address, "nonce", nonce)
		account.Nonce = nonce
	}

	return
}

// BatchTransferAccount :Realize batch transfer local currency transactions by calling contracts
func BatchTransferAccount(client *ethclient.Client, account *Account, cfg *config.Config, amount *big.Int, accounts []*Account) {
	batchTransferAddress := cfg.BatchTransferContract
	price, _ := new(big.Int).SetString("1000000000", 10)

	handOutPool, err := contract.NewHandOutPool(batchTransferAddress, client)
	auth, err := bind.NewKeyedTransactorWithChainID(account.PK, new(big.Int).SetInt64(cfg.ChainId))
	if err != nil {
		log.Errorf("NewKeyedTransactorWithChainID Error:%r", err)
	}
	auth.GasLimit = 10000000
	auth.GasPrice = price

	nonce := account.Nonce

	bathRows := make([]common.Address, 0)
	for i, account := range accounts {
		bathRows = append(bathRows, account.Address)
		if len(bathRows) == 100 || i == len(accounts)-1 {
			auth.Nonce = big.NewInt(int64(nonce))
			tran, err := handOutPool.Handout(auth, bathRows, amount)
			if err != nil {
				log.Error("BatchTransferAccount() handOutPool.Handout Error:", err)
			} else {
				log.Info("Batch transaction eth hash:", tran.Hash())
				nonce += 1
			}
			bathRows = make([]common.Address, 0)
			auth.Nonce = big.NewInt(int64(nonce))
		}
	}
	account.Nonce = nonce
	log.Info("Bath tran eth ending ..............")
}

// BatchTransferAccountToken :Realize batch transfer Token transactions by calling contracts
func BatchTransferAccountToken(client *ethclient.Client, account *Account, cfg *config.Config, amount *big.Int, accounts []*Account) {
	price, _ := new(big.Int).SetString("1000000000", 10)

	tokenAddress := cfg.Erc20Contract
	bathTranContractAddress := cfg.BatchTransferContract

	handOutPool, err := contract.NewHandOutPool(bathTranContractAddress, client)
	auth, err := bind.NewKeyedTransactorWithChainID(account.PK, new(big.Int).SetInt64(cfg.ChainId))
	if err != nil {
		log.Error("NewKeyedTransactorWithChainID Error:", err)
	}
	auth.GasLimit = 90000000
	auth.GasPrice = price

	nonce := account.Nonce

	bathRows := make([]common.Address, 0)
	for i, account := range accounts {
		bathRows = append(bathRows, account.Address)
		if len(bathRows) == 100 || i == len(accounts)-1 {
			auth.Nonce = big.NewInt(int64(nonce))
			tran, err := handOutPool.HandoutToken(auth, tokenAddress, bathRows, amount)
			if err != nil {
				log.Error("handOutPool.HandoutToken error:", err)
			} else {
				nonce += 1
				log.Info("Bath tran token hash:", tran.Hash())
			}
			bathRows = make([]common.Address, 0)
			auth.Nonce = big.NewInt(int64(nonce))
		}
	}
	account.Nonce = nonce
	log.Info("Bath tran token ending ..............")
}
