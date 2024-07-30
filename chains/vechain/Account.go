package vechain

import (
	"crypto/ecdsa"
	"encoding/json"
	"github.com/ethereum/go-ethereum/crypto"
	log "github.com/sirupsen/logrus"
	"github.com/vechain/thor/thor"
	"github.com/vechain/thor/tx"
	"math/big"
	"os"
	"strings"
)

type AccountInfo struct {
	Address string            `json:"address"`
	Private string            `json:"private"`
	Nonce   int               `json:"nonce"`
	PK      *ecdsa.PrivateKey `json:"-"`
}
type txConfig struct {
	receiver thor.Address
	amount   *big.Int
	chainTag byte
}

func (acc *AccountInfo) MakeTx(cfg txConfig, nonce int) *tx.Transaction {
	cla := tx.NewClause(&cfg.receiver).WithValue(cfg.amount)
	tx := new(tx.Builder).
		ChainTag(cfg.chainTag).
		GasPriceCoef(1).
		Expiration(10000).
		Gas(21000).
		Nonce(uint64(nonce)).
		Clause(cla).
		BlockRef(tx.NewBlockRef(0)).
		Build()

	sig, err := crypto.Sign(tx.SigningHash().Bytes(), acc.PK)
	if err != nil {
		log.WithError(err).Error("sign tx failed")
		return nil
	}
	tx = tx.WithSignature(sig)
	return tx
}

func GetAccountJson(accountFile string) []*AccountInfo {
	data, err := os.ReadFile(accountFile)
	if err != nil || len(data) == 0 {
		return []*AccountInfo{}
	} else {
		accs := make([]*AccountInfo, 0)
		err = json.Unmarshal(data, &accs)
		if err != nil {
			log.Error("unmarshal account failed", "err", err)
		}

		for _, acc := range accs {
			private := acc.Private
			if strings.HasPrefix(private, "0x") {
				private = acc.Private[2:]
			}
			acc.PK, err = crypto.HexToECDSA(private)
			if err != nil {
				log.Error("hex to ecdsa failed", "err", err)
			}
		}
		log.Info("get accounts from json", "len", len(accs))
		return accs
	}
}
