package clientpool

import (
	"github.com/xueqianLu/txpress/config"
	"github.com/ethereum/go-ethereum/ethclient"
	"sync"
)

type poolConfig struct {
	RPCNode []string
	Index int32
	clients []*ethclient.Client
	mux sync.Mutex
}

var pconfig = &poolConfig{}

func InitPool(cfg *config.Config) {
	pconfig.Index = 0
	pconfig.RPCNode = cfg.RpcNode
	pconfig.clients = make([]*ethclient.Client,0)
	for i:=0; i < len(pconfig.RPCNode); i++ {
		c,_ := ethclient.Dial(pconfig.RPCNode[i])
		pconfig.clients = append(pconfig.clients, c)
	}
}

func (p *poolConfig) getClient() *ethclient.Client {
	p.mux.Lock()
	defer p.mux.Unlock()
	c := p.clients[p.Index]
	p.Index += 1
	if pconfig.Index >= int32(len(p.RPCNode)) {
		p.Index = 0
	}
	return c
}

func GetClient() *ethclient.Client {
	return pconfig.getClient()
}