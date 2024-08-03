package hackcontrol

import (
	"context"
	log "github.com/sirupsen/logrus"
	pb "github.com/xueqianLu/txpress/hackcontrol/hackcenter"
	"github.com/xueqianLu/txpress/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type HackControl struct {
	conf   types.HackControlConfig
	client pb.CenterServiceClient
	quit   chan struct{}
	chain  types.ChainPlugin
}

func NewHackControl(chain types.ChainPlugin, conf types.HackControlConfig) *HackControl {
	hc := &HackControl{
		conf:  conf,
		chain: chain,
		quit:  make(chan struct{}),
	}
	conn, err := grpc.NewClient(conf.Url,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(1024*1024*1024),
			grpc.MaxCallSendMsgSize(1024*1024*1024)),
	)
	if err != nil {
		log.WithError(err).Error("grpc server connect failed")
		return nil
	}
	hc.client = pb.NewCenterServiceClient(conn)
	return hc
}

func (hc *HackControl) Loop() {
	tm := time.NewTicker(time.Second * 3)
	defer tm.Stop()
	last := 0
	for {
		select {
		case <-hc.quit:
			return
		case <-tm.C:
			// 1. get latest block number.
			blk, err := hc.chain.LatestBlockInfo()
			if err != nil {
				continue
			}
			if last == 0 {
				last = int(blk.Number)
				continue
			}
			if int(blk.Number) >= (last + hc.conf.IntervalBlock) {
				begin := blk.Number + 10
				end := begin + int64(hc.conf.KeepBlock)
				_, err := hc.client.UpdateHack(context.Background(), &pb.UpdateHackRequest{
					Begin: begin,
					End:   end,
				})
				if err != nil {
					log.WithError(err).Error("update hack failed")
				} else {
					log.WithFields(log.Fields{
						"begin": begin,
						"end":   end,
					}).Info("set hack success")
					last = int(blk.Number)
				}
			}
		}
	}

}

func (hc *HackControl) Close() {
	close(hc.quit)
}
