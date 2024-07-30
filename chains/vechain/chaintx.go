package vechain

import (
	"github.com/vechain/thor/tx"
)

type VeTx struct {
	*tx.Transaction
}

func (v VeTx) IsChainTx() bool {
	return true
}
