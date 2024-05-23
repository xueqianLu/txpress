package main

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/xueqianLu/txpress/clientpool"
	"github.com/xueqianLu/txpress/config"
	"github.com/xueqianLu/txpress/tool"
	"io/fs"
	"math/big"
	"os"
	"runtime/pprof"
)

var (
	initAccCount     int64
	initEthAmount    int64
	cpuProfile       bool
	configpath       string
	startCommand     bool
	noCheckNonce     bool
	tokentransaction bool
)

func init() {
	rootCmd.PersistentFlags().BoolVar(&startCommand, "start", false, "Start after initializing the account")
	rootCmd.PersistentFlags().BoolVar(&tokentransaction, "token", false, "Start test with erc20 token transfer")
	rootCmd.PersistentFlags().BoolVar(&cpuProfile, "cpuProfile", false, "Statistics cpu profile")
	rootCmd.PersistentFlags().StringVar(&configpath, "config", "app.json", "config file path")
	rootCmd.PersistentFlags().BoolVar(&noCheckNonce, "nocheck", false, "no need check account nonce")

	accountCmd.PersistentFlags().Int64Var(&initAccCount, "count", 1, "Init account count")
	accountCmd.PersistentFlags().Int64Var(&initEthAmount, "eth", 1, "Init account balance,default: 1ETH")
	rootCmd.AddCommand(accountCmd)
	rootCmd.AddCommand(versionCmd)

	cfg, err := config.ParseConfig(configpath)
	if err != nil {
		log.Error("parse config failed", "err", err)
		os.Exit(1)
	}

	clientpool.InitPool(cfg)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Errorf("Program execute error: %s", err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "txpress",
	Short: "Stress test tools",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("check start and ", "start is", startCommand)
		cfg := config.GetConfig()
		if startCommand {
			accounts := tool.GetAccountJson(cfg)
			if len(accounts) == 0 {
				log.Error("get count failed")
				os.Exit(1)
			}
			if len(accounts) > cfg.Count {
				allaccounts := accounts
				accounts = make([]*tool.Account, cfg.Count)
				copy(accounts, allaccounts[:cfg.Count])
			}

			if !noCheckNonce {
				taskpool := make(chan interface{}, 1000000)
				for _, account := range accounts {
					taskpool <- account
				}
				tasks := tool.NewTasks(10, func(task interface{}) {
					client := clientpool.GetClient()
					account := task.(*tool.Account)
					tool.CheckAccountNonce(client, account)
				}, taskpool)
				tasks.Run()
				close(taskpool)
				tasks.Done()
			}

			if cpuProfile {
				f, err := os.Create("cpuprofile.log")
				if err != nil {
					log.Fatal(err)
				}
				err = pprof.StartCPUProfile(f)
				if err != nil {
					log.Error("Start cpu profile err:", err)
					return
				}
			}
			if tokentransaction {
				cfg.Type = 1
			}

			start(config.GetConfig(), accounts)
			if cpuProfile {
				pprof.StopCPUProfile()
			}

		}
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of txpress",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Version: ", Version)
		log.Info("Git Commit: ", Commit)
	},
}

var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "Account cmd to create account",
	Run: func(cmd *cobra.Command, args []string) {
		type Info struct {
			Balance string `json:"balance"`
		}
		cfg := config.GetConfig()
		accounts := tool.CreateAccounts(cfg, int(initAccCount))
		if len(accounts) > 0 {
			var infos = make(map[string]Info)
			balance := new(big.Int).Mul(big.NewInt(1e18), big.NewInt(initEthAmount))
			// export to genesis format.
			for _, account := range accounts {
				infos[account.Address.String()] = Info{
					Balance: fmt.Sprintf("0x%s", balance.Text(16)),
				}
			}
			d, _ := json.MarshalIndent(infos, "", "  ")
			err := os.WriteFile("balance.json", d, fs.ModePerm)
			if err != nil {
				log.Error("write account info failed", "err", err)
			}
		} else {
			log.Error("create accounts failed")
		}
	},
}
