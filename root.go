package main

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/crypto"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/xueqianLu/txpress/chains"
	"github.com/xueqianLu/txpress/config"
	"github.com/xueqianLu/txpress/types"
	"github.com/xueqianLu/txpress/workflow"
	"io"
	"os"
	"runtime/pprof"
	"strconv"
	"time"
)

var (
	cpuProfile   bool
	configpath   string
	startCommand bool
	logfile      string
)

func init() {
	rootCmd.PersistentFlags().BoolVar(&startCommand, "start", false, "Start after initializing the account")
	rootCmd.PersistentFlags().BoolVar(&cpuProfile, "cpuProfile", false, "Statistics cpu profile")
	rootCmd.PersistentFlags().StringVar(&configpath, "config", "app.json", "config file path")
	rootCmd.PersistentFlags().StringVar(&logfile, "log", "", "log file path")

	accountCmd.Flags().Int("count", 1, "account count")
	accountCmd.Flags().String("balance", "1", "account balance")
	accountCmd.Flags().String("nonce", "0", "account nonce")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(accountCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Errorf("Program execute error: %s", err)
		os.Exit(1)
	}
}

func logInit() {
	if logfile != "" {
		file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			log.SetOutput(io.MultiWriter(file, os.Stdout))
		} else {
			log.Info("Failed to log to file, using default stderr")
		}
	}
}

var rootCmd = &cobra.Command{
	Use:   "txpress",
	Short: "Stress test tools",
	Run: func(cmd *cobra.Command, args []string) {
		logInit()

		log.Info("check start and ", "start is", startCommand)

		cfg, err := config.ParseConfig(configpath)
		if err != nil {
			os.Exit(1)
		}
		if startCommand {
			var allchain []types.ChainPlugin
			for {
				allchain = chains.NewChains(cfg)
				if len(allchain) == 0 {
					log.Error("have no chain to start, wait")
					time.Sleep(3 * time.Second)
					continue
				}
				break
			}

			flow := workflow.NewWorkFlow(allchain, types.RunConfig{
				BaseCount:     cfg.BaseCount,
				Round:         cfg.Round,
				Batch:         cfg.Batch,
				Interval:      time.Duration(cfg.Interval) * time.Second,
				IncRate:       cfg.IncRate,
				BeginToStart:  cfg.BeginToStart,
				ForceIncrease: cfg.ForceIncrease,
			})

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
			flow.Start()

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
	Short: "Create account",
	Run: func(cmd *cobra.Command, args []string) {
		// account count.
		// get account count from a flag.
		accfile := "accounts.json"
		genfile := "gen-alloc.json"
		count, _ := strconv.ParseInt(cmd.Flags().Lookup("count").Value.String(), 10, 64)
		balance := cmd.Flags().Lookup("balance").Value.String()
		nonce, _ := strconv.ParseInt(cmd.Flags().Lookup("nonce").Value.String(), 10, 64)

		type GenesisInfo struct {
			Alloc map[string]map[string]string `json:"alloc"`
		}
		type AccountInfo struct {
			Address string `json:"address"`
			Private string `json:"private"`
			Nonce   int    `json:"nonce"`
		}
		accounts := make([]AccountInfo, 0)
		genesisInfo := GenesisInfo{
			Alloc: make(map[string]map[string]string),
		}
		for i := int64(0); i < count; i++ {
			pk, err := crypto.GenerateKey()
			if err != nil {
				log.Error("Generate key error: ", err)
				return
			}
			address := crypto.PubkeyToAddress(pk.PublicKey).Hex()
			private := pkPadding(pk.D.Text(16))
			accounts = append(accounts, AccountInfo{
				Address: address,
				Private: private,
				Nonce:   int(nonce),
			})
			usergeninfo := make(map[string]string)
			usergeninfo["balance"] = toWeiHex(balance)
			genesisInfo.Alloc[address] = usergeninfo
		}
		accinfo, _ := json.MarshalIndent(accounts, "", "    ")
		if err := os.WriteFile(accfile, accinfo, 0666); err != nil {
			log.Error("Write account info error: ", err)
			return
		}
		geninfo, _ := json.MarshalIndent(genesisInfo, "", "    ")
		if err := os.WriteFile(genfile, geninfo, 0666); err != nil {
			log.Error("Write genesis info error: ", err)
			return
		}
		log.Infof("account generate success, please see %s and %s", accfile, genfile)
		return
	},
}
