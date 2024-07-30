package chains

import (
	log "github.com/sirupsen/logrus"
	"github.com/xueqianLu/txpress/chains/ethchain"
	"github.com/xueqianLu/txpress/chains/vechain"
	"github.com/xueqianLu/txpress/types"
)

func NewChains(config types.ChainConfig) []types.ChainPlugin {
	chains := make([]types.ChainPlugin, 0)

	var createFunc func(rpc string, index int, config types.ChainConfig) (types.ChainPlugin, error)
	switch config.Name {
	case "eth":
		createFunc = ethchain.NewEthChain
	case "vechain":
		createFunc = vechain.NewVeChain
	default:
		log.Errorf("unsupport chain %s", config.Name)
		return nil
	}

	for i, rpc := range config.Rpcs {
		chain, err := createFunc(rpc, i, config)
		if err != nil {
			log.Errorf("create chain %s for with rpc(%s) failed", config.Name, rpc)
			continue
		}
		chains = append(chains, chain)
	}
	return chains
}
